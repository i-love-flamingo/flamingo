package canonicalUrl

import (
	"flamingo.me/flamingo/core/canonicalUrl/interfaces"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/template"
)

// Module for core/canonicalUrl
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*template.Function)(nil)).To(interfaces.CanonicalDomainFunc{})
	injector.BindMulti((*template.Function)(nil)).To(interfaces.IsExternalUrl{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(interfaces.CanonicalUrlFunc{})
}
