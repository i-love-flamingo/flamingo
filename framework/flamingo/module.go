package flamingo

import (
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module is a generic flamingo module, to reduce boilerplating
type Module struct {
	Routes router.Module
	DependensOn []dingo.Module
}

// Configure the Module and bind routes
func (m *Module) Configure(injector *dingo.Injector) {
	panic("do not use me")
	router.Bind(injector, m.Routes)
}

func (m *Module) Depends() []dingo.Module {
	return m.DependensOn
}
