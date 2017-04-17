package brand

import (
	"flamingo/core/dingo"
	"flamingo/framework/router"
	"flamingo/om3/brand/controller"
)

type Module struct {
	RouterRegistry *router.RouterRegistry `inject:""`
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Handle("brand.view", new(controller.ViewController))
	module.RouterRegistry.Route("/brand/{uid}", "brand.view")
}
