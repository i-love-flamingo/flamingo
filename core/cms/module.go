package cms

import (
	"flamingo/core/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("cms.page.view", new(PageController))
	m.RouterRegistry.Route("/page/{name}", "cms.page.view")
}
