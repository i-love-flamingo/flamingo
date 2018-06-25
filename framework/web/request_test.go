package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	hr, err := http.NewRequest("GET", "http://example.com/?var1=val1&var2=val2&var2=val21", nil)
	assert.NoError(t, err)
	req := RequestFromRequest(hr, nil).WithVars(map[string]string{"var1": "val1"})

	assert.Equal(t, "val1", req.MustQuery1("var1"))
	assert.Equal(t, []string{"val1"}, req.MustQuery("var1"))

	assert.Equal(t, "val2", req.MustQuery1("var2"))
	assert.Equal(t, []string{"val2", "val21"}, req.MustQuery("var2"))

	assert.Equal(t, "val1", req.MustParam1("var1"))

	assert.Equal(t, map[string]string{"var1": "val1"}, req.ParamAll())

	assert.Panics(t, func() {
		req.MustParam1("unknown")
	})

	assert.Panics(t, func() {
		req.MustQuery1("unknown")
	})

	assert.Panics(t, func() {
		req.MustForm1("unknown")
	})

	assert.Nil(t, new(Request).QueryAll())
}
