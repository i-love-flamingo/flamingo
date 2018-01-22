package canonicalUrl

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/template"
)

// Module for core/canonicalUrl
type (
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*event.Subscriber)(nil)).To(canonicalTagCreator{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(CanonicalUrlFunc{})
}
