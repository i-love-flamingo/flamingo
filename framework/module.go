/*
Package framework provides the most necessary basics, such as
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
	"flamingo/framework/dingo"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
	"flamingo/framework/template"
	"flamingo/framework/template_functions"
	"flamingo/framework/web"
)

const (
	VERSION = "1.0"
)

type (
	// InitModule: initial module for basic setup
	InitModule struct{}

	// Module for framework functionality
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

// Configure the InitModule
func (initmodule *InitModule) Configure(injector *dingo.Injector) {
	injector.Bind((*event.Router)(nil)).To(event.DefaultRouter{})
	injector.Bind((*profiler.Profiler)(nil)).To(profiler.NullProfiler{})

	injector.Bind((*web.ContextFactory)(nil)).ToInstance(web.ContextFromRequest)

	injector.Bind(router.Router{}).In(dingo.ChildSingleton).ToProvider(router.NewRouter)
	injector.Bind(router.RouterRegistry{}).In(dingo.Singleton).ToProvider(router.NewRouterRegistry)

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.GetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.URLFunc{})
}

// Configure the Module
func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Route("/_flamingo/json/{Handler}", "_flamingo.json")
	module.RouterRegistry.Handle("_flamingo.json", new(controller.DataController))
	module.RouterRegistry.Handle("session.flash", new(controller.SessionFlashController))
}
