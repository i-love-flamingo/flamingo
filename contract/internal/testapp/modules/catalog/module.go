// Package catalog is a test module with the routes of a product catalog.
package catalog

import (
	"context"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module of the catalog routes.
	Module struct{}

	routes struct {
		responder *web.Responder
	}
)

// Configure DI.
func (*Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

// Inject dependencies.
func (r *routes) Inject(responder *web.Responder) {
	r.responder = responder
}

// Routes of the module.
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/product/:id", "catalog.product")
	registry.HandleGet("catalog.product", func(_ context.Context, req *web.Request) web.Result {
		return r.responder.Data(map[string]string{"id": req.Params["id"]})
	})

	registry.MustRoute("/list", `catalog.list(sort?="name")`)
	registry.HandleGet("catalog.list", func(_ context.Context, req *web.Request) web.Result {
		return r.responder.Data(map[string]string{"sort": req.Params["sort"]})
	})
}
