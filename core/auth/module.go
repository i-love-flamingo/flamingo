package auth

import (
	"encoding/gob"
	"flamingo/core/auth/application"
	"flamingo/core/auth/interfaces"
	"flamingo/framework/dingo"
	"flamingo/framework/router"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Module for core.auth
type Module struct {
	RouterRegistry *router.RouterRegistry `inject:""`
}

// Configure core.auth module
func (m *Module) Configure(injector *dingo.Injector) {
	gob.Register(&oauth2.Token{})
	gob.Register(&oidc.IDToken{})

	injector.Bind(application.AuthManager{}).AsEagerSingleton()

	m.RouterRegistry.Route("/auth/login", "auth.login")
	m.RouterRegistry.Route("/auth/callback", "auth.callback")
	m.RouterRegistry.Route("/auth/logout", "auth.logout")

	m.RouterRegistry.Handle("auth.login", new(interfaces.LoginController))
	m.RouterRegistry.Handle("auth.callback", new(interfaces.CallbackController))
	m.RouterRegistry.Handle("auth.logout", new(interfaces.LogoutController))

	m.RouterRegistry.Handle("user", new(interfaces.UserController))
}
