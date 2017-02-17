package cms

import "flamingo/core/app"

func Register(r *app.ServiceContainer) {
	// default handlers
	r.Handle("cms.page.view", new(PageController))

	// default routes
	r.Route("/page/{name}", "cms.page.view")
}
