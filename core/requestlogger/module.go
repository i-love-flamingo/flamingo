package requestlogger

import (
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
)

type (
	// Module for core/requestlogger
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*router.Filter)(nil)).To(logger{})
}
