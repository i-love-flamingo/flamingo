package flamingo_test

import (
	"context"
	"net/http"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/framework/web"
)

func Example() {
	flamingo.App(
		[]dingo.Module{
			new(requestlogger.Module),
		},
		flamingo.WithRoutes(web.RoutesFunc(func(registry *web.RouterRegistry) {
			registry.MustRoute("/", "index")
			registry.HandleAny("index", func(ctx context.Context, req *web.Request) web.Result {
				return new(web.Responder).HTTP(http.StatusOK, strings.NewReader("Hello World!"))
			})
		})),
	)
}
