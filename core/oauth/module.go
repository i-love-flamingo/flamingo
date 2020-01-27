package oauth

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
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
	UseFake                     bool   `inject:"config:core.oauth.useFake"`
	PreventSimultaneousSessions bool   `inject:"config:core.oauth.preventSimultaneousSessions"`
	SessionBackend              string `inject:"config:flamingo.session.backend"`
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

	injector.BindMap(new(auth.RequestIdentifierFactory), "flamingo.core.oauth").ToInstance(func(config config.Map) (auth.RequestIdentifier, error) {
		return &interfaces.LegacyIdentifier{}, nil
	})
}

// CueConfig for oauth module
func (*Module) CueConfig() string {
	return `
core oauth: {
	server: string
	secret: string
	clientid: string
	disableOfflineToken: bool | *false
	enabled: bool | *true
	useFake: bool | *false
	fakeUserData: [string]: _
	fakeLoginTemplate: string | *""
	scopes: [...string] | *["profile", "email"]
	claims: {
		idToken: [...string]
		userInfo: [...string]
	}
	tokenExtras: [...string]
	mapping: {
		idToken: { [string]: string } & {
			sub: string | *"sub"
			email: string | *"email"
			name: string | *"name"
		}
		userInfo: { [string]: string } & {
			sub: string | *"sub"
			email: string | *"email"
			name: string | *"name"
		}
	}
	preventSimultaneousSessions: bool | *false

	legacyAuthIdentifier: {
		broker: "flamingo.core.oauth"
		typ: "flamingo.core.oauth"
	}
}

core: auth: web: broker: [core.oauth.legacyAuthIdentifier, ...]
`
}

// FlamingoLegacyConfigAlias mapping for backwards compatibility
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	alias := make(map[string]string)
	for _, v := range []string{
		"oauth.server",
		"oauth.secret",
		"oauth.clientid",
		"oauth.disableOfflineToken",
		"oauth.enabled",
		"oauth.useFake",
		"oauth.fakeUserData",
		"oauth.fakeLoginTemplate",
		"oauth.scopes",
		"oauth.claims.idToken", "oauth.claims.userInfo",
		"oauth.tokenExtras",
		"oauth.mapping.idToken",
		"oauth.mapping.userInfo",
		"oauth.preventSimultaneousSessions",
	} {
		alias[v] = "core." + v
	}
	return alias
}

// Depends on the session module
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(flamingo.SessionModule),
	}
}

type routes struct {
	login    interfaces.LoginControllerInterface
	logout   interfaces.LogoutControllerInterface
	callback interfaces.CallbackControllerInterface
	user     *interfaces.UserController
	UseFake  bool `inject:"config:core.oauth.useFake"`
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
