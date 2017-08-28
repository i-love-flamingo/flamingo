package template_functions

import (
	"flamingo/framework/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsJson(t *testing.T) {
	var jsJson template.Function = new(JsJson)

	assert.Equal(t, jsJson.Name(), "JSON")

	json := jsJson.Func().(func() Json)()
	assert.Equal(t, `{"foo":123}`, json.Stringify(map[string]int{"foo": 123}))
}
