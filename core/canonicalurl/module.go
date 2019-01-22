package canonicalurl

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/canonicalurl/application"
	"flamingo.me/flamingo/v3/core/canonicalurl/interfaces"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module for core/canonicalUrl
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "canonicalDomain", new(interfaces.CanonicalDomainFunc))
	flamingo.BindTemplateFunc(injector, "isExternalUrl", new(interfaces.IsExternalURL))
	flamingo.BindTemplateFunc(injector, "canonicalUrl", new(interfaces.CanonicalURLFunc))

	injector.Bind(new(interfaces.ApplicationService)).To(new(application.Service))
	injector.Bind(new(application.RouterRouter)).To(new(web.Router))
}
