package gotemplate

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRouter struct{}

func (m mockRouter) Relative(name string, params map[string]string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://name-%v.com/param-amount-%d/", name, len(params)))
}

func (m mockRouter) Data(ctx context.Context, handler string, params map[interface{}]interface{}) interface{} {
	var stringParams []string
	for key, value := range params {
		stringParams = append(stringParams, fmt.Sprintf("%v:%v", key, value))
	}

	sort.Strings(stringParams)

	return fmt.Sprintf("%v %v", handler, stringParams)
}

var _ urlRouter = mockRouter{}

func Test_dataFunc_getFunc(t *testing.T) {
	df := &dataFunc{}
	df.Inject(mockRouter{})
	gf := &getFunc{}
	gf.Inject(mockRouter{})

	dataTemplateFunction := df.Func(context.Background()).(func(what string, params ...string) interface{})
	getTemplateFunction := gf.Func(context.Background()).(func(what string, params ...string) interface{})

	// no params
	var what = "test"
	var params []string
	var want = "test []"

	gotDF := dataTemplateFunction(what, params...)
	gotGF := getTemplateFunction(what, params...)
	assert.Equal(t, want, gotDF)
	assert.Equal(t, want, gotGF)

	// valid params
	what = "test"
	params = []string{"key-1", "value-1", "key-2", "value-2"}
	want = "test [key-1:value-1 key-2:value-2]"

	gotDF = dataTemplateFunction(what, params...)
	gotGF = getTemplateFunction(what, params...)
	assert.Equal(t, want, gotDF)
	assert.Equal(t, want, gotGF)

	// invalid params
	what = "test"
	params = []string{"key-1"}
	want = "test []"

	gotDF = dataTemplateFunction(what, params...)
	gotGF = getTemplateFunction(what, params...)
	assert.Equal(t, want, gotDF)
	assert.Equal(t, want, gotGF)
}

func Test_plainHTMLFunc(t *testing.T) {
	var in = "string abc"
	var want template.HTML = "string abc"

	phf := &plainHTMLFunc{}
	templateFunction := phf.Func(context.Background()).(func(in string) template.HTML)

	got := templateFunction(in)
	assert.Equal(t, want, got)
}

func Test_plainJSFunc(t *testing.T) {
	var in = "string abc"
	var want template.JS = "string abc"

	pjf := &plainJSFunc{}
	templateFunction := pjf.Func(context.Background()).(func(in string) template.JS)

	got := templateFunction(in)
	assert.Equal(t, want, got)
}

func Test_urlFunc_Func(t *testing.T) {
	uf := &urlFunc{}
	uf.Inject(mockRouter{})

	templateFunction := uf.Func(context.Background()).(func(where string, params ...string) template.URL)

	// no params
	var where = "test"
	var params []string
	var want template.URL = "http://name-test.com/param-amount-0/"

	got := templateFunction(where, params...)
	assert.Equal(t, want, got)

	// valid params
	where = "test"
	params = []string{"key-1", "value-1", "key-2", "value-2"}
	want = "http://name-test.com/param-amount-2/"

	got = templateFunction(where, params...)
	assert.Equal(t, want, got)

	// invalid params
	where = "abcd"
	params = []string{"key-1"}
	want = "http://name-abcd.com/param-amount-0/"

	got = templateFunction(where, params...)
	assert.Equal(t, want, got)
}
