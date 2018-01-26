package canonicalUrl

import (
	"go.aoe.com/flamingo/core/canonicalUrl/interfaces"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/template"
)

// Module for core/canonicalUrl
type (
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*template.ContextFunction)(nil)).To(interfaces.CanonicalUrlFunc{})
}
