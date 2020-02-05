package fake

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/fake/interfaces"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module provides Fake OpenID Connect support
type (
	Module struct{}

	routes struct {
		fakeController *interfaces.IdpController
	}
)

// Interface compliance compile type checks
var (
	_ dingo.Module           = (*Module)(nil)
	_ dingo.Depender         = (*Module)(nil)
	_ config.CueConfigModule = (*Module)(nil)

	_ web.RoutesModule = (*routes)(nil)
)

// Configure dependency injection
func (*Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))

	injector.BindMap(new(auth.RequestIdentifierFactory), "fake").ToInstance(interfaces.FakeIdentityProviderFactory)
}

// Inject injects routed dependencies
func (r *routes) Inject(fakeController *interfaces.IdpController) *routes {
	r.fakeController = fakeController

	return r
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: auth: fake: {
	UserConfig :: {
		password: string
		otp: string | *""
	}

	loginTemplate: string | *"" 
	userConfig: {
		[string]: UserConfig
	}
	validatePassword: bool | *true
	validateOtp: bool | *false
	usernameFieldId: string | *"username"
	passwordFieldId: string | *"password"
	otpFieldId: string | *"m2fa-otp"
}
`
}

// Depends marks dependency to auth.WebModule
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(auth.WebModule),
	}
}

// Routes configuration
func (r *routes) Routes(router *web.RouterRegistry) {
	_, _ = router.Route(interfaces.FakeAuthURL, "core.auth.fake.auth")
	router.HandleAny("core.auth.fake.auth", r.fakeController.Auth)
}
