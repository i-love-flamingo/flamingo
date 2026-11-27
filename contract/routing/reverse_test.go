package routing_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/web"
)

// reverseProbeModule binds a route that reverses route names in the area that loads the module.
type reverseProbeModule struct{}

func (*reverseProbeModule) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(reverseProbeRoutes))
}

type reverseProbeRoutes struct {
	responder *web.Responder
	router    web.ReverseRouter
}

func (r *reverseProbeRoutes) Inject(responder *web.Responder, router web.ReverseRouter) {
	r.responder = responder
	r.router = router
}

func (r *reverseProbeRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/probe/product", "probe.product")
	registry.HandleGet("probe.product", func(_ context.Context, req *web.Request) web.Result {
		params := map[string]string{"id": "1"}

		relative, err := r.router.Relative("catalog.product", params)
		if err != nil {
			return r.responder.Data(map[string]string{"error": err.Error()})
		}

		absolute, err := r.router.Absolute(req, "catalog.product", params)
		if err != nil {
			return r.responder.Data(map[string]string{"error": err.Error()})
		}

		return r.responder.Data(map[string]string{"relative": relative.String(), "absolute": absolute.String()})
	})

	registry.MustRoute("/probe/auth", "probe.auth")
	registry.HandleGet("probe.auth", func(_ context.Context, _ *web.Request) web.Result {
		broker := map[string]string{"broker": "shop"}
		routes := map[string]map[string]string{
			"core.auth.login": broker, "core.auth.callback": broker, "core.auth.logout": broker,
			"core.auth.logoutall": nil, "core.auth.logoutCallback": nil,
		}
		result := make(map[string]string, len(routes))

		for name, params := range routes {
			reversed, err := r.router.Relative(name, params)
			if err != nil {
				return r.responder.Data(map[string]string{"error": err.Error()})
			}

			result[name] = reversed.String()
		}

		return r.responder.Data(result)
	})
}

func TestReverse_RelativeAndAbsoluteIncludeAreaPrefix(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootAreas(t, nil, sameModules(new(catalog.Module), new(reverseProbeModule)), areaFiles())

	resp := app.Get(t, "/en/probe/product")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.JSONEq(t, `{"relative":"/en/product/1","absolute":"`+app.URL+`/en/product/1"}`, resp.Body)
}

// Applications build login and logout links from these route names.
func TestReverse_AuthRouteNamesStable(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	files := areaFiles()
	files["config.yml"] = "core.auth.web.broker: []\n"

	app := bootAreas(t, []dingo.Module{new(auth.WebModule), new(oauth.Module)}, sameModules(new(reverseProbeModule)), files)

	resp := app.Get(t, "/en/probe/auth")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.JSONEq(t, `{
		"core.auth.login": "/en/core/auth/login/shop",
		"core.auth.callback": "/en/core/auth/callback/shop",
		"core.auth.logout": "/en/core/auth/logout/shop",
		"core.auth.logoutall": "/en/core/auth/logout",
		"core.auth.logoutCallback": "/en/core/auth/logoutCallback"
	}`, resp.Body)
}
