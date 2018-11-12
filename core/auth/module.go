package auth

import (
	"flamingo.me/flamingo/core/auth/application"
	fakeService "flamingo.me/flamingo/core/auth/application/fake"
	"flamingo.me/flamingo/core/auth/interfaces"
	fakeController "flamingo.me/flamingo/core/auth/interfaces/fake"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module for core.auth
type Module struct {
	UseFake bool `inject:"config:auth.useFake"`
}

// Configure core.auth module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(application.AuthManager{}).In(dingo.ChildSingleton)
	injector.Bind((*interfaces.LogoutRedirectAware)(nil)).To(interfaces.DefaultLogoutRedirect{})
	if !m.UseFake {
		injector.Bind((*application.UserServiceInterface)(nil)).To(application.UserService{})
		injector.Bind((*interfaces.LoginControllerInterface)(nil)).To(interfaces.LoginController{})
		injector.Bind((*interfaces.CallbackControllerInterface)(nil)).To(interfaces.CallbackController{})
		injector.Bind((*interfaces.LogoutControllerInterface)(nil)).To(interfaces.LogoutController{})
	} else {
		injector.Bind((*application.UserServiceInterface)(nil)).To(fakeService.UserService{})
		injector.Bind((*interfaces.LoginControllerInterface)(nil)).To(fakeController.LoginController{})
		injector.Bind((*interfaces.CallbackControllerInterface)(nil)).To(fakeController.CallbackController{})
		injector.Bind((*interfaces.LogoutControllerInterface)(nil)).To(fakeController.LogoutController{})
	}

	router.Bind(injector, new(routes))
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"auth": config.Map{
			"useFake":           false,
			"fakeUserData":      config.Map{},
			"fakeLoginTemplate": "",
			"scopes":            config.Slice{"profile", "email"},
			"claims": config.Map{
				"idToken":  config.Slice{},
				"userInfo": config.Slice{},
			},
			"tokenExtras": config.Slice{},
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
	login    interfaces.LoginControllerInterface
	logout   interfaces.LogoutControllerInterface
	callback interfaces.CallbackControllerInterface
	user     *interfaces.UserController
}

// Inject routes dependencies
func (r *routes) Inject(
	login interfaces.LoginControllerInterface,
	logout interfaces.LogoutControllerInterface,
	callback interfaces.CallbackControllerInterface,
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
