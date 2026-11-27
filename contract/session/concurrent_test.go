package session_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/basket"
)

const holdTimeout = 10 * time.Second

// holdResult is the answer to a request that the holder keeps back.
type holdResult struct {
	item   string
	status int
	err    error
}

// Both requests load the session before either one writes. The two saves do not overlap.
func TestConcurrent_OnWriteKeepsWritesToDifferentKeys(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	t.Run("memory", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		holder, extra := newHolder()

		app := bootSession(t, extra, "flamingo.session.saveMode: OnWrite")

		assertBothWritesKept(t, app, holder)
	})

	t.Run("redis", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		holder, extra := newHolder()

		app, _ := bootRedis(t, extra, "flamingo.session.saveMode: OnWrite")

		assertBothWritesKept(t, app, holder)
	})
}

// newHolder gives a holder for the items "a" and "b", and the module that binds it in the root area.
func newHolder() (*basket.Holder, []dingo.Module) {
	holder := &basket.Holder{
		Loaded:  make(chan string, 2),
		Release: map[string]chan struct{}{"a": make(chan struct{}), "b": make(chan struct{})},
	}

	module := dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.Bind(new(basket.Holder)).ToInstance(holder)
	})

	return holder, []dingo.Module{module}
}

// assertBothWritesKept sends two requests with the same session. Both load the session before one of them writes.
// The request of "a" ends first, then the request of "b".
func assertBothWritesKept(t *testing.T, app *harness.App, holder *basket.Holder) {
	t.Helper()

	// the first request gives the session cookie to the cookie jar
	app.Get(t, "/en/basket/count")

	results := make(chan holdResult, 2)

	for _, item := range []string{"a", "b"} {
		go func() {
			results <- hold(t, app, item)
		}()
	}

	loaded := map[string]bool{}

	for range 2 {
		select {
		case item := <-holder.Loaded:
			loaded[item] = true
		case <-time.After(holdTimeout):
			t.Fatalf("not both requests loaded the session, loaded: %v", loaded)
		}
	}

	require.Equal(t, map[string]bool{"a": true, "b": true}, loaded)

	for _, item := range []string{"a", "b"} {
		close(holder.Release[item])

		select {
		case result := <-results:
			require.NoError(t, result.err)
			assert.Equal(t, item, result.item)
			assert.Equal(t, http.StatusOK, result.status)
		case <-time.After(holdTimeout):
			t.Fatalf("no answer after the release of %q", item)
		}
	}

	resp := app.Get(t, "/en/basket")

	assert.JSONEq(t, `["basket.a","basket.b"]`, resp.Body)
}

// hold sends the request "basket.hold" for the item. It uses the client and not App.Get, because it runs in a goroutine.
func hold(t *testing.T, app *harness.App, item string) holdResult {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, app.URL+"/en/basket/hold/"+item, nil)
	if err != nil {
		return holdResult{item: item, err: err}
	}

	resp, err := app.Client.Do(req)
	if err != nil {
		return holdResult{item: item, err: err}
	}

	defer func() { _ = resp.Body.Close() }()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return holdResult{item: item, err: err}
	}

	return holdResult{item: item, status: resp.StatusCode}
}
