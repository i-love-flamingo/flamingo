package redirects

import (
	"flamingo.me/flamingo/core/redirects/infrastructure"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/prefixrouter"
	"flamingo.me/flamingo/framework/router"
)

type (
	// Module for core/redirects
	Module struct {
		UseInRouter       bool `inject:"config:redirects.useInRouter,optional"`
		UseInPrefixRouter bool `inject:"config:redirects.useInPrefixRouter,optional"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(&infrastructure.RedirectData{}).ToProvider(infrastructure.NewRedirectData).AsEagerSingleton()
	injector.Bind(&redirector{}).ToProvider(newRedirector).AsEagerSingleton()

	if m.UseInRouter {
		injector.BindMulti((*router.Filter)(nil)).To(redirector{})
	}

	if m.UseInPrefixRouter {
		injector.BindMulti((*prefixrouter.OptionalHandler)(nil)).AnnotatedWith("primaryHandlers").To(redirector{})
	}
}

// DefaultConfig provider
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"redirects.csv": "resources/redirects.csv",
	}
}
