package routing_test

import (
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

// optionalHandlerFunc is a prefixrouter.OptionalHandler. It returns false when it has answered the request.
type optionalHandlerFunc func(rw http.ResponseWriter, req *http.Request) bool

func (f optionalHandlerFunc) TryServeHTTP(rw http.ResponseWriter, req *http.Request) (bool, error) {
	return f(rw, req), nil
}

// answerPath gives a handler that answers requests for the paths with the status and the body.
// It lets other requests pass. The body shows that nothing wrote after the handler.
func answerPath(status int, body string, paths ...string) optionalHandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) bool {
		if !slices.Contains(paths, req.URL.Path) {
			return true
		}

		rw.WriteHeader(status)
		_, _ = rw.Write([]byte(body))

		return false
	}
}

// handlerModule binds an optional handler of the prefix router under the annotation.
type handlerModule struct {
	annotation string
	handler    optionalHandlerFunc
}

func (m *handlerModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(prefixrouter.OptionalHandler)).AnnotatedWith(m.annotation).ToInstance(m.handler)
}

// muxModule registers "/mux" on the default mux that prefixrouter binds.
type muxModule struct {
	mux *http.ServeMux
}

func (m *muxModule) Inject(mux *http.ServeMux) {
	m.mux = mux
}

func (m *muxModule) Configure(*dingo.Injector) {
	m.mux.HandleFunc("/mux", func(rw http.ResponseWriter, _ *http.Request) {
		_, _ = rw.Write([]byte("mux"))
	})
}

// bootTagged starts an app whose areas mark their answers with the header "X-Area".
func bootTagged(tb testing.TB, root ...dingo.Module) *harness.App {
	tb.Helper()

	areas := []*config.Area{
		config.NewArea("en", []dingo.Module{new(catalog.Module), &areaTagModule{name: "en"}}),
		config.NewArea("de", []dingo.Module{new(catalog.Module), &areaTagModule{name: "de"}}),
	}

	return bootAreas(tb, root, areas, areaFiles())
}

func TestPrefix_RoutesByPathAndHost(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootTagged(t)

	resp := app.Get(t, "/en/product/1")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "en", resp.Header.Get("X-Area"))

	resp = app.Get(t, "/product/1", "Host", "shop.example")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "de", resp.Header.Get("X-Area"))

	resp = app.Get(t, "/product/1")

	assert.Equal(t, http.StatusNotFound, resp.Status, "another host does not reach the area")
}

func TestPrefix_PrimaryHandlerAnswersFirst(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootTagged(t, &handlerModule{annotation: "primaryHandlers", handler: answerPath(http.StatusTeapot, "primary", "/en/product/1")})

	resp := app.Get(t, "/en/product/1")

	assert.Equal(t, http.StatusTeapot, resp.Status)
	assert.Equal(t, "primary", resp.Body, "the area router did not run")

	resp = app.Get(t, "/en/product/2")

	assert.Equal(t, `{"id":"2"}`+"\n", resp.Body, "a handler that returns true lets the request pass")
}

// The fallback handler also matches /en/product/1, but the area has to answer it.
func TestPrefix_FallbackThenDefaultMux(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootTagged(t,
		&handlerModule{annotation: "fallback", handler: answerPath(http.StatusGone, "gone", "/legacy", "/en/product/1")},
		new(muxModule),
	)

	resp := app.Get(t, "/legacy")

	assert.Equal(t, http.StatusGone, resp.Status)
	assert.Equal(t, "gone", resp.Body, "the default mux did not run")

	resp = app.Get(t, "/en/product/1")

	assert.Equal(t, http.StatusOK, resp.Status, "the fallback handler does not see a request that matches an area")
	assert.Equal(t, "en", resp.Header.Get("X-Area"))

	resp = app.Get(t, "/mux")

	assert.Equal(t, "mux", resp.Body, "a request without area or fallback reaches the default mux")

	resp = app.Get(t, "/fr/x")

	assert.Equal(t, http.StatusNotFound, resp.Status)
	assert.Equal(t, "404 page not found\n", resp.Body)
}
