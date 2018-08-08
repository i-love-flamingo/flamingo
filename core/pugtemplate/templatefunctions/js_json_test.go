package templatefunctions

import (
	"testing"

	"flamingo.me/flamingo/framework/template"

	"github.com/stretchr/testify/assert"
)

func TestJsJSON(t *testing.T) {
	var jsJSON template.Func = new(JsJSON)

	json := jsJSON.Func().(func() JSON)()
	assert.Equal(t, `{"foo":123}`, json.Stringify(map[string]int{"foo": 123}))
}
