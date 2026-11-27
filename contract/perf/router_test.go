//go:build !race

package perf_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/web"
)

const (
	routerHitCeiling      = 616
	routerNotFoundCeiling = 637
	routerRuns            = 100
)

// routerAllocs counts the allocations of one request to the path on the handler of the router,
// with a memory session, 5 pass-through filters and the filters of the framework. The count includes the test
// request. A first request checks the status, and that the filters run and see the session.
func routerAllocs(t *testing.T, path string, status int) float64 {
	t.Helper()

	app, cookie := bootPerf(t)

	instance, err := app.Injector(t).GetInstance(new(web.Router))
	require.NoError(t, err)

	router, ok := instance.(*web.Router)
	require.True(t, ok, "the router is %T", instance)

	handler := router.Handler()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()

	filtersWithSession.Store(0)
	handler.ServeHTTP(rec, req)

	require.Equal(t, status, rec.Code)
	require.EqualValues(t, 5, filtersWithSession.Load(), "filters that ran and saw the session")

	return minAllocs(t, routerRuns, func() {
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
		req.AddCookie(cookie)

		handler.ServeHTTP(httptest.NewRecorder(), req)
	})
}

func TestPerf_AllocsRouterHit(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	assertAllocs(t, routerAllocs(t, "/perf/data", http.StatusOK), routerHitCeiling)
}

func TestPerf_AllocsRouterNotFound(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	assertAllocs(t, routerAllocs(t, "/perf/missing", http.StatusNotFound), routerNotFoundCeiling)
}
