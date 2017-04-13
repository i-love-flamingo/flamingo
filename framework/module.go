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
	"flamingo/core/dingo"
	"flamingo/framework/controller"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
	"flamingo/framework/web"
)

type (
	Module struct {
		Router *router.Router `inject:""`
	}

	InitModule struct{}
)

func (initmodule *InitModule) Configure(injector *dingo.Injector) {
	injector.Bind((*event.Router)(nil)).ToProvider(func() event.Router { return new(event.DefaultRouter) })
	injector.Bind((*profiler.Profiler)(nil)).ToProvider(func() profiler.Profiler { return new(profiler.NullProfiler) })

	injector.Bind((*web.ContextFactory)(nil)).ToInstance(web.ContextFromRequest)

	injector.Bind(new(router.Router)).In(dingo.Singleton).ToProvider(router.CreateRouter)
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.Router.Route("/_flamingo/json/{Handler}", "_flamingo.json")
	module.Router.Handle("_flamingo.json", new(controller.DataController))

	//c.Register(web.ContextFactory(web.ContextFromRequest))

	/*
		c.Register(new(template_functions.GetFunc), "template.func")
		c.Register(new(template_functions.URLFunc), "template.func")
	*/
}
