package main

import (
	"context"
	"net/http"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/framework/web"
)

func main() {
	flamingo.App([]dingo.Module{
		new(requestlogger.Module),
		new(helloworldModule),
	})
}

type helloworldModule struct{}

func (*helloworldModule) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct{}

func (*routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/", "index")
	registry.HandleAny("index", func(ctx context.Context, req *web.Request) web.Result {
		return new(web.Responder).HTTP(
			http.StatusOK,
			strings.NewReader(`Hello Flamingo World!`),
		)
	})
}
