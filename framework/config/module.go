package config

import (
	"flamingo.me/flamingo/v3/framework/dingo"
)

// Module defines a dingo module which automatically binds available config
type Module struct {
	Map
}

// Configure the Module
func (m *Module) Configure(injector *dingo.Injector) {
	for k, v := range m.Flat() {
		if v == nil {
			// log.Printf("Warning: %s has nil value Configured!", k)
			continue
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}
}
