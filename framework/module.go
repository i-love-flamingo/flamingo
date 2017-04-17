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
	"flamingo/core/template"
	"flamingo/framework/controller"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
	"flamingo/framework/template_functions"
	"flamingo/framework/web"
)

type (
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}

	InitModule struct{}
)

func (initmodule *InitModule) Configure(injector *dingo.Injector) {
	injector.Bind((*event.Router)(nil)).To(event.DefaultRouter{})
	injector.Bind((*profiler.Profiler)(nil)).To(profiler.NullProfiler{})

	injector.Bind((*web.ContextFactory)(nil)).ToInstance(web.ContextFromRequest)

	injector.Bind(router.Router{}).In(dingo.Singleton).ToProvider(router.NewRouter)
	injector.Bind(router.RouterRegistry{}).In(dingo.Singleton).ToProvider(router.NewRouterRegistry)

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.GetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.URLFunc{})
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Route("/_flamingo/json/{Handler}", "_flamingo.json")
	module.RouterRegistry.Handle("_flamingo.json", new(controller.DataController))
}
