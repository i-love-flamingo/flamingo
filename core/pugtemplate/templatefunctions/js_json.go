package templatefunctions

import (
	"encoding/json"

	"flamingo.me/flamingo/core/pugtemplate/pugjs"
)

type (
	// JsJSON is exported as a template function
	JsJSON struct{}

	// JSON is our Javascript's JSON equivalent
	JSON struct{}
)

// Func returns the JSON object
func (jl JsJSON) Func() interface{} {
	return func() JSON {
		return JSON{}
	}
}

// Stringify rounds a value up to the next biggest integer
func (j JSON) Stringify(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Stringify rounds a value up to the next biggest integer
func (j JSON) Parse(x string) pugjs.Object {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(x), &m)
	if err != nil {
		panic(err)
	}
	return pugjs.Convert(m)
}
