package csp

import (
	"go.aoe.com/flamingo/core/csrfPreventionFilter"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
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
