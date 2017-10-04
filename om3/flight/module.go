package flight

import (
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/om3/flight/interfaces/controller"
)

type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {

	m.RouterRegistry.Handle("flight.api.autosuggest", (*controller.FlightApiController).AutosuggestAction)
	m.RouterRegistry.Handle("flight.api.search.flightsByAirport", (*controller.FlightApiController).SearchFlightsByAirportAction)
	m.RouterRegistry.Handle("flight.api.saveFlight", (*controller.FlightApiController).SaveFlightAction)
	m.RouterRegistry.Handle("flight.api.getSessionFlight", (*controller.FlightApiController).GetSessionFlightAction)

	m.RouterRegistry.Route("/api/flight/autosuggest", "flight.api.autosuggest")
	m.RouterRegistry.Route("/api/flight/searchByAirport", "flight.api.search.flightsByAirport")
	m.RouterRegistry.Route("/api/flight/saveFlight", "flight.api.saveFlight")
	m.RouterRegistry.Route("/api/flight/getSessionFlight", "flight.api.getSessionFlight")
}
