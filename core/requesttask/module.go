package requesttask

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module registers the requestTask request filter
type Module struct{}

// Configure dependency injection
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).To(new(filter))
}
