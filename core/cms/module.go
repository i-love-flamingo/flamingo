package cms

import (
	"flamingo/core/cms/interfaces"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		Debug bool `inject:"config:debug.mode"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("cms.page.view", new(interfaces.ViewController))
	m.RouterRegistry.Handle("cms.block", new(interfaces.DataController))
	m.RouterRegistry.Route("/page/:name", "cms.page.view(name)")
	if m.Debug {
		m.RouterRegistry.Route("/cmstest", "cms.page.view(name='test',template='pages/cms/test')")
	}
}
