package brand

import (
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/om3/brand/interfaces/controller"
)

type Module struct {
	RouterRegistry *router.RouterRegistry `inject:""`
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Handle("brand.view", new(controller.ViewController))
	module.RouterRegistry.Route("/brand/{uid}", "brand.view")
}
