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
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/controller"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/template"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

const (
	// VERSION of flamingo core
	VERSION = "2"
)

type (
	// InitModule initial module for basic setup
	InitModule struct{}

	// Module for framework functionality
	Module struct{}

	routes struct {
		dataController  *controller.DataController
		flashController *controller.SessionFlashController
		render          *controller.Render
		redirect        *controller.Redirect
		errorController *controller.Error
	}
)

// Configure the InitModule
func (initmodule *InitModule) Configure(injector *dingo.Injector) {
	router.Bind(injector, new(routes))

	injector.Bind((*event.Router)(nil)).To(event.DefaultRouter{})

	injector.Bind(router.Router{}).In(dingo.ChildSingleton).ToProvider(router.NewRouter)
	injector.Bind(router.Registry{}).In(dingo.Singleton).ToProvider(router.NewRegistry)
	injector.Bind(new(web.ReverseRouter)).To(router.Router{})
}

// Configure the Module
func (module *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*responder.RedirectAware)(nil)).To(responder.FlamingoRedirectAware{})
	injector.Bind((*responder.RenderAware)(nil)).To(responder.FlamingoRenderAware{})
	injector.Bind((*responder.ErrorAware)(nil)).To(responder.FlamingoErrorAware{})
	injector.Bind((*responder.JSONAware)(nil)).To(responder.FlamingoJSONAware{})

	template.BindFunc(injector, "config", new(config.TemplateFunc))
	template.BindCtxFunc(injector, "setPartialData", new(web.SetPartialDataFunc))
	template.BindCtxFunc(injector, "getPartialData", new(web.GetPartialDataFunc))

	router.Bind(injector, new(routes))
}

func (r *routes) Inject(
	dataController *controller.DataController,
	flashController *controller.SessionFlashController,
	render *controller.Render,
	redirect *controller.Redirect,
	errorController *controller.Error,
) {
	r.dataController = dataController
	r.flashController = flashController
	r.render = render
	r.redirect = redirect
	r.errorController = errorController
}

func (r *routes) Routes(registry *router.Registry) {
	registry.Route("/_flamingo/json/:handler", "flamingo.data.handler")
	registry.HandleGet("flamingo.data.handler", r.dataController.Get)
	registry.HandleData("session.flash", r.flashController.Data)

	registry.HandleAny("flamingo.render", r.render.Render)

	registry.HandleAny("flamingo.redirect", r.redirect.Redirect)
	registry.HandleAny("flamingo.redirectUrl", r.redirect.RedirectURL)
	registry.HandleAny("flamingo.redirectPermanent", r.redirect.RedirectPermanent)
	registry.HandleAny("flamingo.redirectPermanentUrl", r.redirect.RedirectPermanentURL)

	registry.HandleAny(router.FlamingoError, r.errorController.Error)
	registry.HandleAny(router.FlamingoNotfound, r.errorController.NotFound)
}

// DefaultConfig for this module
func (module *Module) DefaultConfig() config.Map {
	return config.Map{
		"debug.mode":                    true,
		"flamingo.router.notfound":      router.FlamingoNotfound,
		"flamingo.router.error":         router.FlamingoError,
		"flamingo.router.timeout":       float64(60000),
		"flamingo.template.err403":      "error/403",
		"flamingo.template.err404":      "error/404",
		"flamingo.template.errWithCode": "error/withCode",
		"flamingo.template.err503":      "error/503",
		"session.name":                  "flamingo",
	}
}
