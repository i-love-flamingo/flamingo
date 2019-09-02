package requestlogger

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module for core/requestlogger
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).To(logger{})
}

// DefaultConfig configures module's default configuration
func (m *Module) DefaultConfig() config.Map {
	return config.Map{}
}
