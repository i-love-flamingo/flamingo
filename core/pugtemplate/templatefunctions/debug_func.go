package templatefunctions

import (
	"encoding/json"
)

type (
	// DebugFunc renders data as JSON, which allows debugging in templates
	// TODO move into profiler ?
	DebugFunc struct{}
)

// Name alias for use in template
func (df DebugFunc) Name() string {
	return "debug"
}

// Func as implementation of debug method
func (df DebugFunc) Func() interface{} {
	return func(o interface{}) string {
		d, _ := json.MarshalIndent(o, "", "    ")
		return string(d)
	}
}
