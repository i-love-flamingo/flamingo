/*
the flamingo package provides the most necessary basics, such as
 - service_locator
 - router
 - web (including context and response)
 - web/responder

Additionally it provides a router at /_flamingo/json/{handler} for convenient access to DataControllers
Additionally it registers two template functions, `get(...)` and `url(...)`
*/
package flamingo

import (
	"flamingo/core/flamingo/controller"
	"flamingo/core/flamingo/service_container"
	"flamingo/core/flamingo/template_functions"
)

// Register flamingo json Handler
func Register(r *service_container.ServiceContainer) {
	r.Route("/_flamingo/json/{Handler}", "_flamingo.json")
	r.Handle("_flamingo.json", new(controller.DataController))

	r.Register(new(template_functions.GetFunc), "template.func")
	r.Register(new(template_functions.UrlFunc), "template.func")
}
