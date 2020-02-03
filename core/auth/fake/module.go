package fake

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/fake/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module provides Fake OpenID Connect support
type (
	Module struct{}

	routes struct {
		fakeController *interfaces.IdpController
	}
)

// Configure dependency injection
func (*Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: auth:
    fake::
	  loginTemplate: string | *"" 
      usernameFieldId: string | *""
      passwordFieldId: string | *""
      otpFieldId: string | *""
      userConfig:
        validatePassword: bool
        validateOtp: bool
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
	_, _ = router.Route("/core/auth/fake/auth", "core.auth.fake.auth")
	router.HandleAny("core.auth.fake.auth", r.fakeController.Auth)
}
