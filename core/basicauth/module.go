package basicauth

import (
	"flamingo.me/dingo"
	authInterfaces "flamingo.me/flamingo/v3/core/auth/interfaces"
	"flamingo.me/flamingo/v3/core/basicauth/interfaces"
)

// Module for core.basicauth
type Module struct {
}

// Configure core.basicauth module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(authInterfaces.Authservice)).To(interfaces.Authservice{})
}

