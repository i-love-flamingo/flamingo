//go:build !race

package perf_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/web"
)

const (
	memorySessionLoadCeiling = 250
	redisSessionLoadCeiling  = 230
	sessionLoadRuns          = 100
)

// sessionLoadAllocs counts the allocations of SessionStore.LoadByRequest for a request with the cookie.
func sessionLoadAllocs(t *testing.T, app *harness.App, cookie *http.Cookie) float64 {
	t.Helper()

	instance, err := app.Injector(t).GetInstance(new(web.SessionStore))
	require.NoError(t, err)

	store, ok := instance.(*web.SessionStore)
	require.True(t, ok, "the store is %T", instance)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/perf/data", nil)
	req.AddCookie(cookie)

	// a load that misses gives an empty session and would allocate less
	session, err := store.LoadByRequest(req.Context(), req)
	require.NoError(t, err)
	require.Equal(t, "a", session.Try("basket"), "the store must find the session of the cookie")

	return minAllocs(t, sessionLoadRuns, func() {
		_, err := store.LoadByRequest(req.Context(), req)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPerf_AllocsMemorySessionLoad(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app, cookie := bootPerf(t)

	assertAllocs(t, sessionLoadAllocs(t, app, cookie), memorySessionLoadCeiling)
}

// miniredis runs in the test process, so its allocations are part of the count.
func TestPerf_AllocsRedisSessionLoad(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app, cookie := bootPerfRedis(t)

	assertAllocs(t, sessionLoadAllocs(t, app, cookie), redisSessionLoadCeiling)
}
