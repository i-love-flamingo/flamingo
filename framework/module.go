// Package framework provides the most necessary basics, such as
//  - service_locator
//  - router
//  - web (including context and response)
//  - web/responder
//
// Additionally it provides a router at /_flamingo/json/{handler} for convenient access to DataControllers
// Additionally it registers two template functions, `get(...)` and `url(...)`
package framework

import (
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/controller"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/profiler"
	"go.aoe.com/flamingo/framework/profiler/collector"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/template"
	"go.aoe.com/flamingo/framework/templatefunctions"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

const (
	// VERSION of flamingo core
	VERSION = "1.0"
)

type (
	// InitModule initial module for basic setup
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

	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.ConfigFunc{})
}

// Configure the Module
func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Route("/_flamingo/json/:handler", "flamingo.data.handler")
	module.RouterRegistry.Handle("flamingo.data.handler", new(controller.DataController))
	module.RouterRegistry.Handle("session.flash", new(controller.SessionFlashController))

	module.RouterRegistry.Handle("flamingo.render", (*controller.Render).Render)

	module.RouterRegistry.Handle("flamingo.redirect", (*controller.Redirect).Redirect)
	module.RouterRegistry.Handle("flamingo.redirectUrl", (*controller.Redirect).RedirectURL)
	module.RouterRegistry.Handle("flamingo.redirectPermanent", (*controller.Redirect).RedirectPermanent)
	module.RouterRegistry.Handle("flamingo.redirectPermanentUrl", (*controller.Redirect).RedirectPermanentURL)

	module.RouterRegistry.Handle(router.FlamingoError, (*controller.Error).Error)
	module.RouterRegistry.Handle(router.FlamingoNotfound, (*controller.Error).NotFound)

	injector.BindMulti((*collector.DataCollector)(nil)).To(router.DataCollector{})

	injector.Bind((*responder.RedirectAware)(nil)).To(responder.FlamingoRedirectAware{})
	injector.Bind((*responder.RenderAware)(nil)).To(responder.FlamingoRenderAware{})
	injector.Bind((*responder.ErrorAware)(nil)).To(responder.FlamingoErrorAware{})
	injector.Bind((*responder.JSONAware)(nil)).To(responder.FlamingoJSONAware{})
}

// DefaultConfig for this module
func (module *Module) DefaultConfig() config.Map {
	return config.Map{
		"debug.mode":               true,
		"flamingo.router.notfound": router.FlamingoNotfound,
		"flamingo.router.error":    router.FlamingoError,
		"flamingo.template.err404": "error/404",
		"flamingo.template.err503": "error/503",
		"session.name":             "flamingo",
	}
}
