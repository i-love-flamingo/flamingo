package flamingo

import (
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module is a generic flamingo module, to reduce boilerplating
type Module struct {
	Routes router.Module
}

// Configure the Module and bind routes
func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, m.Routes)
}
