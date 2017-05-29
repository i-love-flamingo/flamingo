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
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure the InitModule
func (initmodule *InitModule) Configure(injector *dingo.Injector) {
	injector.Bind((*event.Router)(nil)).To(event.DefaultRouter{})
	injector.Bind((*profiler.Profiler)(nil)).To(profiler.NullProfiler{})

	injector.Bind((*web.ContextFactory)(nil)).ToInstance(web.ContextFromRequest)

	injector.Bind(router.Router{}).In(dingo.ChildSingleton).ToProvider(router.NewRouter)
	injector.Bind(router.Registry{}).In(dingo.Singleton).ToProvider(router.NewRegistry)

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.GetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.URLFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.ConfigFunc{})
}

// Configure the Module
func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Mount("/_flamingo/json/:handler", new(controller.DataController))
	module.RouterRegistry.Handle("session.flash", new(controller.SessionFlashController))

	module.RouterRegistry.Handle("flamingo.redirect", (*controller.Redirect).Redirect)
	module.RouterRegistry.Handle("flamingo.redirectUrl", (*controller.Redirect).RedirectUrl)
	module.RouterRegistry.Handle("flamingo.redirectPermanent", (*controller.Redirect).RedirectPermanent)
	module.RouterRegistry.Handle("flamingo.redirectPermanentUrl", (*controller.Redirect).RedirectPermanentUrl)
}
