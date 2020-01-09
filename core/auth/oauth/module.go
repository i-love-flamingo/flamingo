package oauth

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
)

// Module provides OpenID Connect support
type Module struct{}

// Configure dependency injection
func (*Module) Configure(injector *dingo.Injector) {
	// injector.BindMap(new(auth.WebIdentifierFactory), "oauth2").ToInstance(oauth2Factory)
	injector.BindMap(new(auth.IdentifierFactory), "oidc").ToInstance(oidcFactory)
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: auth: {
	oauth2Config :: {
		broker: string
		clientID: string
		clientSecret: string
		endpoint: string
		scopes: [...string] | *["profile", "email"]
		claims: {
			accessToken: { [string]: string }
		}
	}

	oidc :: {
		oauth2Config
		typ: "oidc"
		enableOfflineToken: bool | *true
		claims: {
			idToken: { [string]: string } & {
				sub: string | *"sub"
				email: string | *"email"
				name: string | *"name"
			}
		}
		requestClaims: {
			idToken: [...string]
			userInfo: [...string]
		}
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
