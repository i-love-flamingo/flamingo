package profiler

import (
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/flamingo/profiler"
	"flamingo/core/flamingo/router"
)

func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		r.Route("/_profiler/view/{Profile}", "_profiler.view")
		r.Handle("_profiler.view", new(ProfileController))
	}, router.RouterRegister)

	c.RegisterFactory(func() profiler.Profiler { return new(DefaultProfiler) }, "event.subscriber")
}
