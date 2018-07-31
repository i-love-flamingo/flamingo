package profiler

import (
	"flamingo.me/flamingo/framework/dingo"
)

// Module registers our profiler
// deprecated: use opencensus tracing
type Module struct{}

// Configure DI
func (*Module) Configure(*dingo.Injector) {}
