package template_functions

import "encoding/json"

type (
	// JsonLib is exported as a template function
	JsonLib struct{}

	// Json is our Javascript's JSON equivalent
	Json struct{}
)

// Name alias for use in template
func (jl JsonLib) Name() string {
	return "JSON"
}

// Func as implementation of debug method
func (jl JsonLib) Func() interface{} {
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

func (j Json) NoConvert() {}
