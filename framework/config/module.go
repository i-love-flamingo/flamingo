package config

import (
	"flamingo.me/dingo"
)

// Module defines a dingo module which automatically binds provided config.
// Normaly this module is not included in your flamingo projects bootstrap.
//
// Its can be useful for testing dingo.Module that require certain configuration to be set before. E.g.:
//
// cfgModule := &config.Module{
//		Map: config.Map{
//			"redirects.useInRouter":       true,
//			"redirects.useInPrefixRouter": true,
//		},
//	}
//
//	if err := dingo.TryModule(cfgModule, module); err != nil {
//		t.Error(err)
//	}
type Module struct {
	Map
}

// Configure the Module
func (m *Module) Configure(injector *dingo.Injector) {
	for k, v := range m.Flat() {
		if v == nil {
			continue
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}
}
