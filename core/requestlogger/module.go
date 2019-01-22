package requestlogger

import (
	"flamingo.me/dingo"
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
