package brand

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/om3/brand/interfaces/controller"
)

// Module for om3/brand package
type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

// Configure DI
func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Route("/brand/:uid", "brand.view")
	module.RouterRegistry.Handle("brand.view", new(controller.ViewController))
}
