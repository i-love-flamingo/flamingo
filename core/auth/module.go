package auth

import (
	"flamingo.me/dingo"
	interfaces2 "flamingo.me/flamingo/v3/core/auth/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module for core.auth
type Module struct {
}

// Configure core.auth module
func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	auth *interfaces2.Authcontroller
}

// Inject routes dependencies
func (r *routes) Inject(
	auth *interfaces2.Authcontroller,
) {
	r.auth = auth
}

// Routes module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.Route("/auth/authenticate", `auth.authenticate(redirecturl?="")`)
	registry.HandleGet("auth.authenticate", r.auth.AuthAction)
}
