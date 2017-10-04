package requestlogger

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
)

// Module for core/requestlogger
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*event.Subscriber)(nil)).To(logger{})
}
