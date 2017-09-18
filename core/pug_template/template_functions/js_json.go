package template_functions

import "encoding/json"

type (
	// JsJson is exported as a template function
	JsJson struct{}

	// Json is our Javascript's JSON equivalent
	Json struct{}
)

// Name of JS object
func (jl JsJson) Name() string {
	return "JSON"
}

// Func returns the Json object
func (jl JsJson) Func() interface{} {
	return func() Json {
		return Json{}
	}
}

// Stringify rounds a value up to the next biggest integer
func (j Json) Stringify(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}
