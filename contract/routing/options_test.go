package routing_test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/web"
)

// optionsModule binds an OPTIONS action that wraps a plain http.Handler.
type optionsModule struct {
	sawRequest, sawSession *atomic.Bool
}

func (m *optionsModule) Configure(injector *dingo.Injector) {
	// web.BindRoutes binds the type, and dingo then creates the value itself, without the fields.
	injector.BindMulti(new(web.RoutesModule)).ToInstance(&optionsRoutes{sawRequest: m.sawRequest, sawSession: m.sawSession})
}

type optionsRoutes struct {
	sawRequest, sawSession *atomic.Bool
}

func (r *optionsRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/x", "x")
	registry.HandleOptions("x", web.WrapHTTPHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		r.sawRequest.Store(web.RequestFromContext(req.Context()) != nil)
		r.sawSession.Store(web.SessionFromContext(req.Context()) != nil)

		rw.WriteHeader(http.StatusNoContent)
	})))
}

func TestOptions_WrapHTTPHandlerSeesRequestAndSession(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	var sawRequest, sawSession atomic.Bool

	app := bootAreas(t, nil, sameModules(&optionsModule{sawRequest: &sawRequest, sawSession: &sawSession}), areaFiles())

	req, err := http.NewRequestWithContext(t.Context(), http.MethodOptions, app.URL+"/en/x", nil)
	require.NoError(t, err)

	resp, err := app.Client.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.True(t, sawRequest.Load(), "the context holds the flamingo request")
	assert.True(t, sawSession.Load(), "the context holds the session")
}
