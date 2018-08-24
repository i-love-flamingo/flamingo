package config

import "flamingo.me/flamingo/framework/dingo"

type Module struct {
	Map
}

func (m *Module) Configure(injector *dingo.Injector) {
	for k, v := range m.Flat() {
		if v == nil {
			// log.Printf("Warning: %s has nil value Configured!", k)
			continue
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}
}
