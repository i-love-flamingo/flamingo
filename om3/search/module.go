package search

import (
	"flamingo/om3/search/interfaces"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

// Configure the search URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("search.view", new(interfaces.ViewController))
	m.RouterRegistry.Route("/search", "search.view")
}
