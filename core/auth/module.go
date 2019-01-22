package auth

import (
	"flamingo.me/flamingo/v3/core/auth/application"
	fakeService "flamingo.me/flamingo/v3/core/auth/application/fake"
	"flamingo.me/flamingo/v3/core/auth/application/store"
	"flamingo.me/flamingo/v3/core/auth/interfaces"
	fakeController "flamingo.me/flamingo/v3/core/auth/interfaces/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
)

// Module for core.auth
type Module struct {
	UseFake                     bool   `inject:"config:auth.useFake"`
	PreventSimultaneousSessions bool   `inject:"config:auth.preventSimultaneousSessions"`
	SessionBackend              string `inject:"config:session.backend"`
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

	if !m.PreventSimultaneousSessions {
		injector.Bind((*store.Store)(nil)).To(store.Nil{})
	} else {
		switch m.SessionBackend {
		case "redis":
			injector.Bind((*store.Store)(nil)).To(store.Redis{}).AsEagerSingleton()
		case "memory":
			injector.Bind((*store.Store)(nil)).To(store.Memory{}).AsEagerSingleton()
		case "file":
			injector.Bind((*store.Store)(nil)).To(store.File{}).AsEagerSingleton()
		default:
			injector.Bind((*store.Store)(nil)).To(store.Nil{}).AsEagerSingleton()
		}
	}
	injector.Bind((*application.Synchronizer)(nil)).To(application.SynchronizerImpl{})

	router.Bind(injector, new(routes))
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"auth": config.Map{
			"useFake":                     false,
			"preventSimultaneousSessions": false,
			"fakeUserData":                config.Map{},
			"fakeLoginTemplate":           "",
			"scopes":                      config.Slice{"profile", "email"},
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
