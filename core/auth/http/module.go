package http

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
)

type Module struct{}

func (*Module) Configure(injector *dingo.Injector) {
	injector.BindMap(new(auth.IdentifierFactory), "http").ToInstance(identifierFactory)
}

func (*Module) CueConfig() string {
	return `
core: auth: {
	http :: core.auth.authBroker & {
		typ: "http"
		realm: string
		users: [string]: string
	}
}
`
}

func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(auth.WebModule),
	}
}
