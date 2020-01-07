package auth

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// WebModule registers identification for web requests
type WebModule struct{}

// Configure dependency injection
func (m *WebModule) Configure(injector *dingo.Injector) {
	injector.Bind(new([]RequestIdentifier)).ToProvider(buildAuthentifier)
	injector.Bind(new(WebIdentityService)).In(dingo.ChildSingleton)

	web.BindRoutes(injector, new(routes))
}

func buildAuthentifier(
	provider map[string]IdentifierFactory,
	cfg *struct {
		Config config.Slice `inject:"config:core.auth.web.broker"`
	},
) []RequestIdentifier {
	var broker []config.Map
	cfg.Config.MapInto(&broker)

	res := make([]RequestIdentifier, len(broker))

	for i, broker := range broker {
		if res[i] = provider[broker["typ"].(string)](broker); res[i] == nil {
			panic("can not build broker " + broker["typ"].(string))
		}
	}

	return res
}

// Depends marks the WebModule to depend on the flamingo session module
func (*WebModule) Depends() []dingo.Module {
	return []dingo.Module{
		new(flamingo.SessionModule),
	}
}

type routes struct {
	debugController *debugController
	controller      *Controller
	debug           bool
}

// Inject controller
func (r *routes) Inject(debugController *debugController, controller *Controller, cfg *struct {
	Debug bool `inject:"config:flamingo.debug.mode"`
}) {
	r.debugController = debugController
	r.controller = controller
	r.debug = cfg.Debug
}

// Routes configuration
func (r *routes) Routes(router *web.RouterRegistry) {
	if r.debug {
		router.Route("/core/auth/debug", "core.auth.debug")
		router.HandleAny("core.auth.debug", r.debugController.Action)
	}
	router.Route("/core/auth/callback/:broker", "core.auth.callback(broker)")
	router.HandleAny("core.auth.callback", r.controller.Callback)
}

// CueConfig schema
func (*WebModule) CueConfig() string {
	return `
core: auth: {
	web: {
		broker: [...authBroker]
	}

	authBroker :: {
		broker: string
		typ: string
		[string]: _
	}
}
`
}
