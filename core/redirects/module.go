package redirects

import (
	"go.aoe.com/flamingo/core/prefixrouter"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

// Module for core/redirects
type (
	Module struct {
		UseInRouter       bool `inject:"config:redirects.useInRouter,optional"`
		UseInPrefixRouter bool `inject:"config:redirects.useInPrefixRouter,optional"`
	}
	LogFilter struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	if m.UseInRouter {
		injector.BindMulti((*router.Filter)(nil)).To(redirector{})
	}

	if m.UseInPrefixRouter {
		injector.BindMulti((*prefixrouter.OptionalHandler)(nil)).AnnotatedWith("primaryHandlers").To(redirector{})
	}

}
