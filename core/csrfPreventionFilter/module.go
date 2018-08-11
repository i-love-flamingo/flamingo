package csrfPreventionFilter

import (
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/template"
)

// Module for core/csrfPreventionFilter
type (
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*router.Filter)(nil)).To(csrfFilter{})
	injector.BindMulti((*event.Subscriber)(nil)).To(hiddenCsrfTagCreator{})
	template.BindCtxFunc(injector, "csrftoken", new(CsrfFunc))
	injector.Bind((*NonceGenerator)(nil)).To(UuidGenerator{})
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"csrfPreventionFilter.tokenLimit": 10,
	}
}
