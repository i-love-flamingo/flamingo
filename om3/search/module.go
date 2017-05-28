package search

import (
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/om3/search/interfaces"
)

type (
	// Module registers our search package
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure the search URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("search.view", new(interfaces.ViewController))
	m.RouterRegistry.Route("/search", `search.view(type="product")`)
	m.RouterRegistry.Route("/search/:type", `search.view(type)`)
}
