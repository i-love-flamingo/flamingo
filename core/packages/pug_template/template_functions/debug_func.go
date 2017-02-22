package template_functions

import "encoding/json"

type (
	DebugFunc struct{}
)

func (_ DebugFunc) Name() string {
	return "debug"
}

func (_ DebugFunc) Func() interface{} {
	return func(o interface{}) string {
		d, _ := json.MarshalIndent(o, "", "    ")
		return string(d)
	}
}
