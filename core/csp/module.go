package csp

import (
	"flamingo.me/flamingo/core/csrfPreventionFilter"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module for core/csp
type (
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*router.Filter)(nil)).To(cspFilter{})
	m.RouterRegistry.Route("/_cspreport", "_cspreport.view")
	m.RouterRegistry.Handle("_cspreport.view", new(cspReportController))
	injector.Bind((*csrfPreventionFilter.NonceGenerator)(nil)).To(csrfPreventionFilter.UuidGenerator{})
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cspFilter.reportMode": true,
	}
}
