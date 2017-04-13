package profiler

import (
	di "flamingo/framework/dependencyinjection"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
)

func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		r.Route("/_profiler/view/{profile}", "_profiler.view")
		r.Handle("_profiler.view", new(ProfileController))
	}, router.RouterRegister)

	c.RegisterFactory(func() profiler.Profiler { return new(DefaultProfiler) }, "event.subscriber")
}
