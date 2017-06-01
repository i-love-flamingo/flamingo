package cms

import (
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/core/cms/interfaces"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("cms.page.view", new(interfaces.ViewController))
	m.RouterRegistry.Handle("cms.block", new(interfaces.DataController))
	m.RouterRegistry.Route("/page/:name", "cms.page.view")
}
