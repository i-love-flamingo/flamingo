package templatefunctions

import "encoding/json"

type (
	// JsJSON is exported as a template function
	JsJSON struct{}

	// JSON is our Javascript's JSON equivalent
	JSON struct{}
)

// Name of JS object
func (jl JsJSON) Name() string {
	return "JSON"
}

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
