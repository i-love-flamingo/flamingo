package profiler

import (
	"flamingo/framework/dingo"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Route("/_profiler/view/{profile}", "_profiler.view")
	m.RouterRegistry.Handle("_profiler.view", new(ProfileController))

	injector.Override((*profiler.Profiler)(nil), "").To(DefaultProfiler{})

	injector.BindMulti((*event.Subscriber)(nil)).To(EventSubscriber{})
}
