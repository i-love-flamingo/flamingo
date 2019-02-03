package baseurl

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/baseurl/application"
	"flamingo.me/flamingo/v3/framework/baseurl/domain"
	"flamingo.me/flamingo/v3/framework/baseurl/interfaces"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/template"
)

type (
	// Module basic struct
	Module struct {
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.Service)(nil)).To(&application.Service{})

	template.BindFunc(injector, "canonicalDomain", new(interfaces.CanonicalDomainFunc))
	template.BindFunc(injector, "isExternalUrl", new(interfaces.IsExternalURL))
}

// DefaultConfig for baseurl module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"baseurl.url": "",
	}
}
