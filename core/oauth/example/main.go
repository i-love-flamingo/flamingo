package main

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/http"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/framework/web"
)

type testapp struct{}

func (*testapp) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	c *controller
}

func (r *routes) Inject(c *controller) {
	r.c = c
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	_, _ = registry.Route("/", "example")
	registry.HandleAny("example", r.c.Handle)
}

type controller struct {
	am        *application.AuthManager
	responder *web.Responder
}

func (c *controller) Inject(am *application.AuthManager, responder *web.Responder) {
	c.am = am
	c.responder = responder
}

func (c *controller) Handle(ctx context.Context, r *web.Request) web.Result {
	a, err := c.am.Auth(ctx, r.Session())
	return c.responder.Data(map[string]interface{}{
		"auth":  a,
		"error": err,
	})
}

func main() {
	flamingo.App([]dingo.Module{
		new(requestlogger.Module),
		new(auth.WebModule),
		new(oauth.Module),
		new(http.Module),
		new(testapp),
	})
}
