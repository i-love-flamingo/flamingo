package template_functions

import "encoding/json"

type (
	DebugFunc struct{}
)

// Name alias for use in template
func (_ DebugFunc) Name() string {
	return "debug"
}

// Func as implementation of debug method
func (_ DebugFunc) Func() interface{} {
	return func(o interface{}) string {
		d, _ := json.MarshalIndent(o, "", "    ")
		return string(d)
	}
}
