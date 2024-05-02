// Package framework provides the most necessary basics, such as
//   - service_locator
//   - router
//   - web (including context and response)
//   - web/responder
//
// Additionally it provides a router at /_flamingo/json/{handler} for convenient access to DataControllers
// Additionally it registers two template functions, `get(...)` and `url(...)`
package framework

import (
	"flamingo.me/dingo"
	"github.com/spf13/cobra"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/controller"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/web/filter"
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
func (*InitModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(cobra.Command)).ToProvider(web.RoutesCmd)
	injector.BindMulti(new(cobra.Command)).ToProvider(web.HandlerCmd)
	injector.BindMulti(new(cobra.Command)).ToProvider(config.ModulesCmd)
	injector.BindMulti(new(cobra.Command)).ToProvider(config.Cmd)

	web.BindRoutes(injector, new(routes))

	injector.Bind(new(flamingo.EventRouter)).To(flamingo.DefaultEventRouter{})

	injector.Bind(web.Router{}).In(dingo.ChildSingleton)
	injector.Bind(new(web.ReverseRouter)).To(web.Router{})
	injector.Bind(web.RouterRegistry{}).In(dingo.Singleton).ToProvider(web.NewRegistry)
	injector.BindMulti(new(web.Filter)).To(new(filter.MetricsFilter))

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

// CueConfig definition for flamingo framework
func (*InitModule) CueConfig() string {
	return `
flamingo: {
	debug: mode: bool | *true
	router: {
		notfound: string | *"flamingo.notfound"
		error: string | *"flamingo.error"
		timeout: int | *60000
		scheme?: string
		host?: string
		path?: string
	}
	template: {
		err400: string | *"error/400"
		err403: string | *"error/403"
		err404: string | *"error/404"
		errWithCode: string | *"error/withCode"
		err503: string | *"error/503"
	}
	session: {
		name: string | *"flamingo"
		saveMode: *"Always" | "OnRead" | "OnWrite" 
	}
}
`
}

// FlamingoLegacyConfigAlias maps legacy configuration to new
func (*InitModule) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"debug.mode":   "flamingo.debug.mode",
		"session.name": "flamingo.session.name",
	}
}
