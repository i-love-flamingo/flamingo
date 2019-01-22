package csrfPreventionFilter

import (
	"flamingo.me/flamingo/v3/core/csrfPreventionFilter/application"
	"flamingo.me/flamingo/v3/core/csrfPreventionFilter/interfaces"
	"flamingo.me/flamingo/v3/core/csrfPreventionFilter/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/core/form2/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
	"flamingo.me/flamingo/v3/framework/template"
)

// Module for core/csrfPreventionFilter
type (
	Module struct {
		All bool `inject:"config:csrf.all"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*application.Service)(nil)).To(application.ServiceImpl{})
	template.BindCtxFunc(injector, "csrfToken", new(templatefunctions.CsrfTokenFunc))
	template.BindCtxFunc(injector, "csrfInput", new(templatefunctions.CsrfInputFunc))

	injector.BindMap((*domain.FormExtension)(nil), "formExtension.csrfToken").To(interfaces.CrsfTokenFormExtension{})

	if m.All {
		injector.BindMulti((*router.Filter)(nil)).To(interfaces.CsrfFilter{})
	}
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"csrf.all":    false,
		"csrf.secret": "somethingSuperSecret",
		"csrf.ttl":    900.0,
	}
}
