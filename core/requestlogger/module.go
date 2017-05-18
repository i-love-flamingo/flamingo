package requestlogger

import (
	"flamingo/framework/dingo"
	"flamingo/framework/event"
)

// Module for core/requestlogger
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*event.Subscriber)(nil)).To(Logger{})
}
