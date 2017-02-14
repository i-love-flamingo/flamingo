package cms

import "flamingo/core/core/app"

func Register(r *app.Registrator) {
	var pc PageController

	// default handlers
	r.Handle("cms.page.view", &pc)

	// default routes
	r.Route("/page/{name}", "cms.page.view")
}
