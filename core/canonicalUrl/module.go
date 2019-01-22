package canonicalUrl

import (
	"flamingo.me/flamingo/v3/core/canonicalUrl/application"
	"flamingo.me/flamingo/v3/core/canonicalUrl/interfaces"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
	"flamingo.me/flamingo/v3/framework/template"
)

// Module for core/canonicalUrl
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	template.BindFunc(injector, "canonicalDomain", new(interfaces.CanonicalDomainFunc))
	template.BindFunc(injector, "isExternalUrl", new(interfaces.IsExternalUrl))
	template.BindCtxFunc(injector, "canonicalUrl", new(interfaces.CanonicalUrlFunc))

	injector.Bind(new(interfaces.ApplicationService)).To(new(application.Service))
	injector.Bind(new(application.RouterRouter)).To(new(router.Router))
}
