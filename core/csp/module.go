package csp

import (
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
	injector.Bind((*NonceGenerator)(nil)).To(UuidGenerator{})
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cspFilter.reportMode": true,
	}
}

type routes struct {
	controller *cspReportController
}

func (r *routes) Inject(controller *cspReportController) {
	r.controller = controller
}

func (r *routes) Routes(registry *router.Registry) {
	registry.Route("/_cspreport", "_cspreport.view")
	registry.HandlePost("_cspreport.view", r.controller.Post)
}
