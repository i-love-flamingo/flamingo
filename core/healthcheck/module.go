package healthcheck

import (
	"flamingo.me/flamingo/core/healthcheck/interfaces/controllers"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("healthcheck", new(controllers.Healthcheck))
	m.RouterRegistry.Route("/status/healthcheck", "healthcheck")
}
