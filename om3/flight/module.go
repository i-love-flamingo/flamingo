package flight

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/om3/flight/interfaces/controller"
)

type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {

	m.RouterRegistry.Handle("flight.api.search.autosuggest", (*controller.FlightApiController).AutosuggestAction)
	m.RouterRegistry.Handle("flight.api.search.flights", (*controller.FlightApiController).SearchFlightsAction)
	m.RouterRegistry.Handle("flight.api.search.airports", (*controller.FlightApiController).SearchAirportsAction)
	m.RouterRegistry.Handle("flight.api.search.flightsPerAirline", (*controller.FlightApiController).SearchFlightsPerAirlineAction)
	m.RouterRegistry.Handle("flight.api.saveFlight", (*controller.FlightApiController).SaveFlightAction)
	m.RouterRegistry.Handle("flight.api.getSessionFlight", (*controller.FlightApiController).GetSessionFlightAction)

	m.RouterRegistry.Route("/api/flight/search/autosuggest", "flight.api.search.autosuggest")
	m.RouterRegistry.Route("/api/flight/search/flights", "flight.api.search.flights")
	m.RouterRegistry.Route("/api/flight/search/airports", "flight.api.search.airports")
	m.RouterRegistry.Route("/api/flight/search/flightsPerAirline", "flight.api.search.flightsPerAirline")
	m.RouterRegistry.Route("/api/flight/saveFlight", "flight.api.saveFlight")
	m.RouterRegistry.Route("/api/flight/getSessionFlight", "flight.api.getSessionFlight")
}


