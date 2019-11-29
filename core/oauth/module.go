package oauth

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/oauth/application"
	fakeService "flamingo.me/flamingo/v3/core/oauth/application/fake"
	"flamingo.me/flamingo/v3/core/oauth/interfaces"
	fakeController "flamingo.me/flamingo/v3/core/oauth/interfaces/fake"
	"flamingo.me/flamingo/v3/core/security/application/role"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module for core.auth
type Module struct {
	UseFake                     bool   `inject:"config:oauth.useFake"`
	PreventSimultaneousSessions bool   `inject:"config:oauth.preventSimultaneousSessions"`
	SessionBackend              string `inject:"config:session.backend"`
}

// Configure core.auth module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(application.AuthManager{}).In(dingo.ChildSingleton)
	injector.Bind(new(interfaces.LogoutRedirectAware)).To(interfaces.DefaultLogoutRedirect{})
	flamingo.BindEventSubscriber(injector).To(&application.EventHandler{})
	if !m.UseFake {
		injector.Bind(new(application.UserServiceInterface)).To(application.UserService{})
		injector.Bind(new(interfaces.LoginControllerInterface)).To(interfaces.LoginController{})
		injector.Bind(new(interfaces.CallbackControllerInterface)).To(interfaces.CallbackController{})
		injector.Bind(new(interfaces.LogoutControllerInterface)).To(interfaces.LogoutController{})
	} else {
		injector.Bind(new(application.UserServiceInterface)).To(fakeService.UserService{})
		injector.Bind(new(interfaces.LoginControllerInterface)).To(fakeController.LoginController{})
		injector.Bind(new(interfaces.CallbackControllerInterface)).To(fakeController.CallbackController{})
		injector.Bind(new(interfaces.LogoutControllerInterface)).To(fakeController.LogoutController{})
	}

	injector.BindMulti(new(role.Provider)).To(application.AuthRoleProvider{})

	web.BindRoutes(injector, new(routes))
}

// DefaultConfig for auth module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"oauth": config.Map{
			"enabled":           true,
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
			"preventSimultaneousSessions": false,
		},
	}
}

type routes struct {
	login    interfaces.LoginControllerInterface
	logout   interfaces.LogoutControllerInterface
	callback interfaces.CallbackControllerInterface
	user     *interfaces.UserController
	UseFake  bool `inject:"config:oauth.useFake"`
}

// Inject routes dependencies
func (r *routes) Inject(
	login interfaces.LoginControllerInterface,
	logout interfaces.LogoutControllerInterface,
	callback interfaces.CallbackControllerInterface,
	user *interfaces.UserController,
	fake *bool,
) {
	r.login = login
	r.logout = logout
	r.callback = callback
	r.user = user
	r.UseFake = *fake
}

// Routes module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.Route("/auth/login", `auth.login(redirecturl?="")`)
	registry.HandleGet("auth.login", r.login.Get)
	if r.UseFake {
		registry.Route("/auth/callback", `auth.callback(group?="")`)
	} else {
		registry.Route("/auth/callback", `auth.callback`)
	}
	registry.HandleGet("auth.callback", r.callback.Get)
	registry.Route("/auth/logout", "auth.logout")
	registry.HandleGet("auth.logout", r.logout.Get)

	registry.HandleData("user", r.user.Data)
}
