package harness_test

import (
	"context"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/web"
)

// helloModule binds the routes of a minimal app.
type helloModule struct {
	greeting string
}

func (m *helloModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.RoutesModule)).ToInstance(&helloRoutes{greeting: m.greeting})
}

type helloRoutes struct {
	greeting string
}

func (r *helloRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/hello", "hello")
	registry.HandleGet("hello", func(_ context.Context, _ *web.Request) web.Result {
		return new(web.Responder).HTTP(http.StatusOK, strings.NewReader(r.greeting))
	})

	registry.MustRoute("/host", "host")
	registry.HandleGet("host", func(_ context.Context, req *web.Request) web.Result {
		return new(web.Responder).HTTP(http.StatusOK, strings.NewReader(req.Request().Host))
	})
}

func TestHarness_BootTwiceInOneProcess(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	var addresses []string

	for _, name := range []string{"first", "second"} {
		t.Run(name, func(t *testing.T) {
			app := harness.Boot(t, []dingo.Module{&helloModule{greeting: name}}, harness.Config{})
			addresses = append(addresses, strings.TrimPrefix(app.URL, "http://"))

			resp := app.Get(t, "/hello")

			assert.Equal(t, http.StatusOK, resp.Status)
			assert.Equal(t, name, resp.Body)
		})
	}

	for _, address := range addresses { //nolint:paralleltest // checks the apps after their subtests ended
		conn, err := new(net.Dialer).DialContext(t.Context(), "tcp", address)
		if !assert.Error(t, err, "the app still listens after its subtest") {
			_ = conn.Close()
		}
	}
}

// net/http ignores a "Host" entry in Request.Header.
func TestHarness_HostPairSetsRequestHost(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, []dingo.Module{new(helloModule)}, harness.Config{})

	body := app.Get(t, "/host", "Host", "shop.example").Body

	assert.Equal(t, "shop.example", body)
}
