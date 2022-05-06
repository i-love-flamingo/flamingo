package systemendpoint

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/systemendpoint/application"
)

type (
	// Module basic struct
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).To(&application.SystemServer{}).In(dingo.Singleton)
}

// CueConfig for the module
func (*Module) CueConfig() string {
	return `flamingo: systemendpoint: serviceAddr: string | *":13210"`
}

// FlamingoLegacyConfigAlias maps legacy config to new
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"systemendpoint.serviceAddr": "flamingo.systemendpoint.serviceAddr",
	}
}
