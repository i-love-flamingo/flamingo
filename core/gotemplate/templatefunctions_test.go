package gotemplate

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"testing"
)

type MockRouter struct{}

func (m MockRouter) Relative(name string, params map[string]string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://name-%v.com/param-amount-%d/", name, len(params)))
}

func (m MockRouter) Data(ctx context.Context, handler string, params map[interface{}]interface{}) interface{} {
	return fmt.Sprintf("%v %v", handler, params)
}

var _ urlRouter = MockRouter{}

func Test_dataFunc_Func(t *testing.T) {
	type args struct {
		what   string
		params []string
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "no params",
			args: args{
				what: "test",
			},
			want: "test map[]",
		},
		{
			name: "valid params set",
			args: args{
				what:   "test",
				params: []string{"key-1", "value-1", "key-2", "value-2"},
			},
			want: "test map[key-1:value-1 key-2:value-2]",
		},
		{
			name: "invalid params set",
			args: args{
				what:   "test",
				params: []string{"key-1"},
			},
			want: "test map[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := &dataFunc{}
			df.Inject(MockRouter{})

			templateFunction := df.Func(context.Background()).(func(what string, params ...string) interface{})
			if got := templateFunction(tt.args.what, tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_getFunc_Func(t *testing.T) {
	type args struct {
		what   string
		params []string
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "no params",
			args: args{
				what: "test",
			},
			want: "test map[]",
		},
		{
			name: "valid params set",
			args: args{
				what:   "test",
				params: []string{"key-1", "value-1", "key-2", "value-2"},
			},
			want: "test map[key-1:value-1 key-2:value-2]",
		},
		{
			name: "invalid params set",
			args: args{
				what:   "test",
				params: []string{"key-1"},
			},
			want: "test map[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gf := &getFunc{}
			gf.Inject(MockRouter{})

			templateFunction := gf.Func(context.Background()).(func(what string, params ...string) interface{})
			if got := templateFunction(tt.args.what, tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_plainHTMLFunc_Func(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want template.HTML
	}{
		{
			name: "string gets converted to template.HTML",
			in:   "string abc",
			want: "string abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phf := &plainHTMLFunc{}
			templateFunction := phf.Func(context.Background()).(func(in string) template.HTML)
			if got := templateFunction(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("plainHTMLFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_plainJSFunc_Func(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want template.JS
	}{
		{
			name: "string gets converted to template.JS",
			in:   "string abc",
			want: "string abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pjf := &plainJSFunc{}
			templateFunction := pjf.Func(context.Background()).(func(in string) template.JS)
			if got := templateFunction(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("plainJSFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_urlFunc_Func(t *testing.T) {
	type args struct {
		where  string
		params []string
	}

	tests := []struct {
		name string
		args args
		want template.URL
	}{
		{
			name: "no params",
			args: args{
				where: "test",
			},
			want: "http://name-test.com/param-amount-0/",
		},
		{
			name: "valid params set",
			args: args{
				where:  "abc",
				params: []string{"key-1", "value-1", "key-2", "value-2"},
			},
			want: "http://name-abc.com/param-amount-2/",
		},
		{
			name: "invalid params set",
			args: args{
				where:  "abcd",
				params: []string{"key-1"},
			},
			want: "http://name-abcd.com/param-amount-0/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf := &urlFunc{}
			uf.Inject(MockRouter{})

			templateFunction := uf.Func(context.Background()).(func(where string, params ...string) template.URL)
			if got := templateFunction(tt.args.where, tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("urlFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}
