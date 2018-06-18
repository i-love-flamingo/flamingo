package auth

import (
	"encoding/gob"

	"github.com/coreos/go-oidc"
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/interfaces"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/profiler/collector"
	"flamingo.me/flamingo/framework/router"
	"golang.org/x/oauth2"
)

// Module for core.auth
type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

// Configure core.auth module
func (m *Module) Configure(injector *dingo.Injector) {
	gob.Register(&oauth2.Token{})
	gob.Register(&oidc.IDToken{})

	injector.Bind(application.AuthManager{}).In(dingo.ChildSingleton)
	injector.Bind((*interfaces.LogoutRedirectAware)(nil)).To(interfaces.DefaultLogoutRedirect{})
	injector.Bind((*application.UserServiceInterface)(nil)).To(application.UserService{})

	m.RouterRegistry.Route("/auth/login", `auth.login(redirecturl?="")`)
	m.RouterRegistry.Handle("auth.login", new(interfaces.LoginController))
	m.RouterRegistry.Route("/auth/callback", "auth.callback")
	m.RouterRegistry.Handle("auth.callback", new(interfaces.CallbackController))
	m.RouterRegistry.Route("/auth/logout", "auth.logout")
	m.RouterRegistry.Handle("auth.logout", new(interfaces.LogoutController))

	m.RouterRegistry.Handle("user", new(interfaces.UserController))

	injector.BindMulti((*collector.DataCollector)(nil)).To(application.DataCollector{})
}
