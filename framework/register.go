/*
Package flamingo provides the most necessary basics, such as
 - service_locator
 - router
 - web (including context and response)
 - web/responder

Additionally it provides a router at /_flamingo/json/{handler} for convenient access to DataControllers
Additionally it registers two template functions, `get(...)` and `url(...)`
*/
package framework

import (
	"flamingo/framework/controller"
	di "flamingo/framework/dependencyinjection"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
	"flamingo/framework/template_functions"
	"flamingo/framework/web"
)

// Register flamingo json Handler
func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		r.Route("/_flamingo/json/{Handler}", "_flamingo.json")
		r.Handle("_flamingo.json", new(controller.DataController))
	}, router.RouterRegister)

	c.Register(web.ContextFactory(web.ContextFromRequest))

	c.RegisterFactory(func() event.Router { return new(event.DefaultRouter) })
	c.RegisterFactory(func() profiler.Profiler { return new(profiler.NullProfiler) })

	c.Register(new(template_functions.GetFunc), "template.func")
	c.Register(new(template_functions.URLFunc), "template.func")
}
