package healthcheck

import (
	"go.aoe.com/flamingo/core/healthcheck/interfaces/controllers"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("healthcheck", new(controllers.Healthcheck))
	m.RouterRegistry.Route("/status/healthcheck", "healthcheck")
}
