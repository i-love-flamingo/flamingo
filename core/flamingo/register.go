/*
Package flamingo provides the most necessary basics, such as
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
	"flamingo/core/flamingo/event"
	"flamingo/core/flamingo/service_container"
	"flamingo/core/flamingo/template_functions"
)

// Register flamingo json Handler
func Register(sc *service_container.ServiceContainer) {
	sc.Route("/_flamingo/json/{Handler}", "_flamingo.json")
	sc.Handle("_flamingo.json", new(controller.DataController))

	sc.Register(func() event.Router { return new(event.DefaultRouter) })

	sc.Register(new(template_functions.GetFunc), "template.func")
	sc.Register(new(template_functions.URLFunc), "template.func")
}
