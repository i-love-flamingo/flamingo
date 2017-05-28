package brand

import (
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/om3/brand/interfaces/controller"
)

// Module for om3/brand package
type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

// Configure DI
func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Mount("/brand/:uid", new(controller.ViewController))
}
