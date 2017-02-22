package cms

import "flamingo/core/flamingo/service_container"

func Register(r *service_container.ServiceContainer) {
	// default handlers
	r.Handle("cms.page.view", new(PageController))

	// default routes
	r.Route("/page/{name}", "cms.page.view")
}
