package csrfPreventionFilter

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/template"
)

// Module for core/csrfPreventionFilter
type (
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*router.Filter)(nil)).To(csrfFilter{})
	injector.BindMulti((*event.Subscriber)(nil)).To(hiddenCsrfTagCreator{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(CsrfFunc{})
	injector.Bind((*NonceGenerator)(nil)).To(uuidGenerator{})
}
