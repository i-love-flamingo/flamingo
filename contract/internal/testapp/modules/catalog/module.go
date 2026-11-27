// Package catalog is a test module with the routes of a product catalog.
package catalog

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/web"
)

// cacheLifetime is the max-age in seconds of the route "/cached".
const cacheLifetime = 60

var (
	errMissing     = errors.New("catalog-missing")
	errForbidden   = errors.New("catalog-forbidden")
	errUnavailable = errors.New("catalog-unavailable")
	errSecretPanic = errors.New("catalog-secret-panic")
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

	r.responderRoutes(registry)
}

// responderRoutes registers the routes for the results of the responder.
func (r *routes) responderRoutes(registry *web.RouterRegistry) {
	registry.MustRoute("/page", "catalog.page")
	registry.HandleGet("catalog.page", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Render("page", nil)
	})

	registry.MustRoute("/notfound", "catalog.notfound")
	registry.HandleGet("catalog.notfound", func(ctx context.Context, _ *web.Request) web.Result {
		return r.responder.NotFoundWithContext(ctx, errMissing)
	})

	registry.MustRoute("/forbidden", "catalog.forbidden")
	registry.HandleGet("catalog.forbidden", func(ctx context.Context, _ *web.Request) web.Result {
		return r.responder.ForbiddenWithContext(ctx, errForbidden)
	})

	registry.MustRoute("/unavailable", "catalog.unavailable")
	registry.HandleGet("catalog.unavailable", func(ctx context.Context, _ *web.Request) web.Result {
		return r.responder.UnavailableWithContext(ctx, errUnavailable)
	})

	registry.MustRoute("/redirect/route", "catalog.redirect.route")
	registry.HandleGet("catalog.redirect.route", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.RouteRedirect("catalog.product", map[string]string{"id": "1"})
	})

	registry.MustRoute("/redirect/url", "catalog.redirect.url")
	registry.HandleGet("catalog.redirect.url", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.URLRedirect(&url.URL{Scheme: "https", Host: "example.com", Path: "/target"})
	})

	// The template "broken" calls a template that does not exist, so it fails to render.
	registry.MustRoute("/broken", "catalog.broken")
	registry.HandleGet("catalog.broken", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Render("broken", nil)
	})

	registry.MustRoute("/panic", "catalog.panic")
	registry.HandleGet("catalog.panic", func(_ context.Context, _ *web.Request) web.Result {
		panic(errSecretPanic)
	})

	registry.MustRoute("/cached", "catalog.cached")
	registry.HandleGet("catalog.cached", func(_ context.Context, _ *web.Request) web.Result {
		result := r.responder.Data(map[string]string{"cache": "yes"})
		result.CacheDirective = web.CacheDirectiveBuilder{
			IsReusable:              true,
			AllowIntermediateCaches: true,
			MaxCacheLifetime:        cacheLifetime,
		}.Build()

		return result
	})

	registry.MustRoute("/nocache", "catalog.nocache")
	registry.HandleGet("catalog.nocache", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Data(map[string]string{"cache": "no"}).SetNoCache()
	})

	registry.MustRoute("/download", "catalog.download")
	registry.HandleGet("catalog.download", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Download(strings.NewReader("a,b\n"), "text/csv", "list.csv", true)
	})
}
