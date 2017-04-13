package cms

import (
	di "flamingo/framework/dependencyinjection"
	"flamingo/framework/router"
)

// Register adds handles for cms page routes.
func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		// default handlers
		r.Handle("cms.page.view", new(PageController))

		// default routes
		r.Route("/page/{name}", "cms.page.view")
	}, "router.register")
}
