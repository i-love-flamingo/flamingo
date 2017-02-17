package cms

import "flamingo/core/flamingo"

func Register(r *flamingo.ServiceContainer) {
	// default handlers
	r.Handle("cms.page.view", new(PageController))

	// default routes
	r.Route("/page/{name}", "cms.page.view")
}
