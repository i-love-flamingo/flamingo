package systemendpoint

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/systemendpoint/application"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
)

type (
	// Module basic struct
	Module struct {
		handlerProvider domain.HandlerProvider
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).To(&application.SystemServer{})
}

// DefaultConfig for the module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"systemendpoint.serviceAddr": ":13210",
	}
}
