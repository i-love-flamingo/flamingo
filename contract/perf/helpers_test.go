//go:build !race

package perf_test

import (
	"context"
	"math"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// The race detector adds allocations, so the files of this package are not built with -race.
// Each ceiling is the count of go1.27.4 on darwin/arm64 plus 10%, rounded up.

// Background goroutines, such as the view worker of opencensus, also allocate, so the count is the lowest
// of allocRounds rounds.
const allocRounds = 5

// minAllocs gives the minimum of allocRounds testing.AllocsPerRun(runs, f) rounds.
func minAllocs(t *testing.T, runs int, f func()) float64 {
	t.Helper()

	lowest := math.MaxFloat64

	for range allocRounds {
		lowest = min(lowest, testing.AllocsPerRun(runs, f))
	}

	return lowest
}

// assertAllocs logs the count and fails the test when it is above the ceiling.
func assertAllocs(t *testing.T, got, ceiling float64) {
	t.Helper()

	t.Logf("%v allocations per run, ceiling %v", got, ceiling)

	if got > ceiling {
		t.Errorf("%v allocations per run, ceiling %v: find the new allocation, or raise the ceiling", got, ceiling)
	}
}

// filtersWithSession counts the calls of passFilter that found the stored value in the session.
var filtersWithSession atomic.Int32

// passFilter passes each request on.
type passFilter struct{}

func (passFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	if req.Session().Try("basket") == "a" {
		filtersWithSession.Add(1)
	}

	return chain.Next(ctx, req, w)
}

// perfModule binds the routes "/perf/data" and "/perf/store" and 5 pass-through filters.
type perfModule struct{}

func (*perfModule) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(perfRoutes))

	for range 5 {
		injector.BindMulti(new(web.Filter)).ToInstance(passFilter{})
	}
}

type perfRoutes struct {
	responder *web.Responder
}

func (r *perfRoutes) Inject(responder *web.Responder) {
	r.responder = responder
}

func (r *perfRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/perf/data", "perf.data")
	registry.HandleGet("perf.data", func(context.Context, *web.Request) web.Result {
		return r.responder.Data(map[string]string{"a": "b"})
	})

	registry.MustRoute("/perf/store", "perf.store")
	registry.HandleGet("perf.store", func(_ context.Context, req *web.Request) web.Result {
		req.Session().Store("basket", "a")

		return r.responder.Data(map[string]string{})
	})
}

// bootPerf starts an app with the session module and the perf module. The config lines are values of config.yml.
// The returned cookie belongs to a session that holds one value.
func bootPerf(t *testing.T, lines ...string) (*harness.App, *http.Cookie) {
	t.Helper()

	lines = append(lines, "flamingo.session.cookie.secure: false")

	app := harness.Boot(t, []dingo.Module{new(flamingo.SessionModule), new(perfModule)}, harness.Config{
		Files: map[string]string{"config.yml": strings.Join(lines, "\n") + "\n"},
	})

	resp := app.Get(t, "/perf/store")

	for _, cookie := range (&http.Response{Header: resp.Header}).Cookies() {
		if cookie.Name == "flamingo" {
			return app, cookie
		}
	}

	t.Fatalf("response has no session cookie, Set-Cookie: %q", resp.Header.Values("Set-Cookie"))

	return nil, nil
}

// bootPerfRedis starts miniredis and an app that keeps its sessions in it.
func bootPerfRedis(t *testing.T) (*harness.App, *http.Cookie) {
	t.Helper()

	server := miniredis.RunT(t)

	app, cookie := bootPerf(t,
		"flamingo.session.backend: redis",
		"flamingo.session.redis.host: \""+server.Addr()+"\"",
	)

	return app, cookie
}
