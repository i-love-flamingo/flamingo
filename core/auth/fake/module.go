package fake

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
)

// Module provides OpenID Connect support
type Module struct{}

// Configure dependency injection
func (*Module) Configure(injector *dingo.Injector) {

}

// CueConfig schema
func (*Module) CueConfig() string {
	return `

`
}

// Depends marks dependency to auth.WebModule
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(auth.WebModule),
	}
}
