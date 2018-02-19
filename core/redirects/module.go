package redirects

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

// Module for core/redirects
type (
	Module struct{}
	LogFilter struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*router.Filter)(nil)).To(redirect0r{})
}
