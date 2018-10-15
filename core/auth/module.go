package auth

import (
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/interfaces"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module for core.auth
type Module struct{}

// Configure core.auth module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(application.AuthManager{}).In(dingo.ChildSingleton)
	injector.Bind((*interfaces.LogoutRedirectAware)(nil)).To(interfaces.DefaultLogoutRedirect{})
	injector.Bind((*application.UserServiceInterface)(nil)).To(application.UserService{})

	router.Bind(injector, new(routes))
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"auth": config.Map{
			"scopes": config.Slice{"profile", "email"},
			"claims": config.Map{
				"idToken":  config.Slice{},
				"userInfo": config.Slice{},
			},
			"mapping": config.Map{
				"idToken": config.Map{
					"sub":   "sub",
					"email": "email",
					"name":  "name",
				},
				"userInfo": config.Map{
					"sub":   "sub",
					"email": "email",
					"name":  "name",
				},
			},
		},
	}
}

type routes struct {
	login    *interfaces.LoginController
	logout   *interfaces.LogoutController
	callback *interfaces.CallbackController
	user     *interfaces.UserController
}

// Inject routes dependencies
func (r *routes) Inject(
	login *interfaces.LoginController,
	logout *interfaces.LogoutController,
	callback *interfaces.CallbackController,
	user *interfaces.UserController,
) {
	r.login = login
	r.logout = logout
	r.callback = callback
	r.user = user
}

// Routes module
func (r *routes) Routes(registry *router.Registry) {
	registry.Route("/auth/login", `auth.login(redirecturl?="")`)
	registry.HandleGet("auth.login", r.login.Get)
	registry.Route("/auth/callback", "auth.callback")
	registry.HandleGet("auth.callback", r.callback.Get)
	registry.Route("/auth/logout", "auth.logout")
	registry.HandleGet("auth.logout", r.logout.Get)

	registry.HandleData("user", r.user.Data)
}
