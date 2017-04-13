package brand

import (
	"flamingo/core/dingo"
	"flamingo/framework/router"
	"flamingo/om3/brand/controller"
)

type Module struct {
	Router *router.Router `inject:""`
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.Router.Handle("brand.view", new(controller.ViewController))
	module.Router.Route("/brand/{uid}", "brand.view")
}
