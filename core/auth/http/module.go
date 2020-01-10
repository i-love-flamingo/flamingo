package http

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
)

// Module for HTTP authentication
type Module struct{}

// Configure dependency injection
func (*Module) Configure(injector *dingo.Injector) {
	injector.BindMap(new(auth.RequestIdentifierFactory), "http").ToInstance(identifierFactory)
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: auth: {
	http :: {
		typ: "http"
		broker: string
		realm: string
		users: [string]: string
	}
}
`
}

// Depends on auth.WebModule
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(auth.WebModule),
	}
}
