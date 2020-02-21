package fake

import (
	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
)

// Module provides Fake OpenID Connect support
type (
	Module struct{}
)

// Interface compliance compile type checks
var (
	_ dingo.Module           = (*Module)(nil)
	_ dingo.Depender         = (*Module)(nil)
	_ config.CueConfigModule = (*Module)(nil)
)

// Configure dependency injection
func (*Module) Configure(injector *dingo.Injector) {
	injector.BindMap(new(auth.RequestIdentifierFactory), "fake").ToInstance(FakeIdentityProviderFactory)
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: auth: {
	fake :: {
		UserConfig :: {
			password?: string
		}

		typ: "fake"
		broker: string
		loginTemplate?: string
		userConfig: {
			[string]: UserConfig
		}

		validatePassword: bool | *true
		usernameFieldId: string | *"username"
		passwordFieldId: string | *"password"
	}
}
`
}

// Depends marks dependency to auth.WebModule
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(auth.WebModule),
	}
}
