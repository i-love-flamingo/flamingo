package profiler

import (
	"flamingo.me/flamingo/v3/framework/dingo"
)

// Module registers our profiler
// deprecated: use opencensus tracing
type Module struct{}

// Configure DI
func (*Module) Configure(*dingo.Injector) {}
