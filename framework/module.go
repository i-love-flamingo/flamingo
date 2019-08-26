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
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/controller"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
)

const (
	// VERSION of flamingo core
	VERSION = "3"
)

type (
	// InitModule initial module for basic setup
	InitModule struct{}

	routes struct {
		flashController *controller.SessionFlashController
		render          *controller.Render
		redirect        *controller.Redirect
		errorController *controller.Error
		static          *controller.Static
	}
)

// Configure the InitModule
func (initmodule *InitModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(cobra.Command)).ToProvider(web.RoutesCmd)
	injector.BindMulti(new(cobra.Command)).ToProvider(web.HandlerCmd)
	injector.BindMulti(new(cobra.Command)).ToProvider(config.Cmd)

	web.BindRoutes(injector, new(routes))

	injector.Bind(new(flamingo.EventRouter)).To(flamingo.DefaultEventRouter{})

	injector.Bind(web.Router{}).In(dingo.ChildSingleton)
	injector.Bind(new(web.ReverseRouter)).To(web.Router{})
	injector.Bind(web.RouterRegistry{}).In(dingo.Singleton).ToProvider(web.NewRegistry)

	flamingo.BindTemplateFunc(injector, "config", new(config.TemplateFunc))
	flamingo.BindTemplateFunc(injector, "setPartialData", new(web.SetPartialDataFunc))
	flamingo.BindTemplateFunc(injector, "getPartialData", new(web.GetPartialDataFunc))
	flamingo.BindTemplateFunc(injector, "canonicalDomain", new(web.CanonicalDomainFunc))
	flamingo.BindTemplateFunc(injector, "isExternalUrl", new(web.IsExternalURL))
}

// Inject controller for flamingo default handler
func (r *routes) Inject(
	flashController *controller.SessionFlashController,
	render *controller.Render,
	redirect *controller.Redirect,
	errorController *controller.Error,
	staticController *controller.Static,
) {
	r.flashController = flashController
	r.render = render
	r.redirect = redirect
	r.errorController = errorController
	r.static = staticController
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleData("session.flash", r.flashController.Data)

	registry.HandleAny("flamingo.render", r.render.Render)

	registry.HandleAny("flamingo.redirect", r.redirect.Redirect)
	registry.HandleAny("flamingo.redirectUrl", r.redirect.RedirectURL)
	registry.HandleAny("flamingo.redirectPermanent", r.redirect.RedirectPermanent)
	registry.HandleAny("flamingo.redirectPermanentUrl", r.redirect.RedirectPermanentURL)
	registry.HandleAny("flamingo.static.file", r.static.File)

	registry.HandleAny(web.FlamingoError, r.errorController.Error)
	registry.HandleAny(web.FlamingoNotfound, r.errorController.NotFound)
}

// DefaultConfig for this module
func (initmodule *InitModule) DefaultConfig() config.Map {
	return config.Map{
		"debug.mode":                    true,
		"flamingo.router.notfound":      web.FlamingoNotfound,
		"flamingo.router.error":         web.FlamingoError,
		"flamingo.router.timeout":       float64(60000),
		"flamingo.template.err403":      "error/403",
		"flamingo.template.err404":      "error/404",
		"flamingo.template.errWithCode": "error/withCode",
		"flamingo.template.err503":      "error/503",
		"session.name":                  "flamingo",
	}
}
