package brand

import (
di "flamingo/framework/dependencyinjection"
"flamingo/framework/router"
"flamingo/om3/brand/controller"
)

// Register Services for brand package
func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		r.Handle("brand.view", new(controller.ViewController))
		r.Route("/brand/{uid}", "brand.view")
	}, router.RouterRegister)
}

