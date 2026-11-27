package bootstrap_test

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/web"
)

// fixedReverseRouter gives the same URL for every route.
type fixedReverseRouter struct{}

func (fixedReverseRouter) Relative(string, map[string]string) (*url.URL, error) {
	return &url.URL{Path: "/fixed"}, nil
}

func (fixedReverseRouter) Absolute(*web.Request, string, map[string]string) (*url.URL, error) {
	return &url.URL{Scheme: "http", Host: "shop.example", Path: "/fixed"}, nil
}

// overrideModule replaces the framework binding of web.ReverseRouter and has a route that uses it.
type overrideModule struct{}

func (*overrideModule) Configure(injector *dingo.Injector) {
	injector.Override(new(web.ReverseRouter), "").To(fixedReverseRouter{})
	web.BindRoutes(injector, new(overrideRoutes))
}

type overrideRoutes struct {
	reverseRouter web.ReverseRouter
}

func (r *overrideRoutes) Inject(reverseRouter web.ReverseRouter) {
	r.reverseRouter = reverseRouter
}

func (r *overrideRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/relative", "relative")
	registry.HandleGet("relative", func(ctx context.Context, _ *web.Request) web.Result {
		relative, err := r.reverseRouter.Relative("x", nil)
		if err != nil {
			return new(web.Responder).ServerErrorWithContext(ctx, err)
		}

		return new(web.Responder).HTTP(http.StatusOK, strings.NewReader(relative.String()))
	})
}

func TestOverride_AppModuleReplacesFrameworkBinding(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, []dingo.Module{new(overrideModule)}, harness.Config{})

	body := app.Get(t, "/relative").Body

	assert.Equal(t, "/fixed", body)
}
