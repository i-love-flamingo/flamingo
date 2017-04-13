package profiler

import (
	"flamingo/core/dingo"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		Router *router.Router `inject:""`
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
	m.Router.Route("/_profiler/view/{profile}", "_profiler.view")
	m.Router.Handle("_profiler.view", new(ProfileController))

	injector.Bind(new(DefaultProfiler)).In(new(dingo.RequestScope))

	// Use a profiler to inject scope-bound DefaultProfiler
	injector.Bind(new(profiler.Profiler)).ToProvider(ProfilerProvider)
	injector.BindMulti(new(event.Subscriber)).ToProvider(ProfilerProvider)
}

func ProfilerProvider(profiler *DefaultProfiler) *DefaultProfiler {
	return profiler
}
