package metrics_test

import (
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/tag"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
	"flamingo.me/flamingo/v3/framework/web"
)

// controllerRecorder keeps the value of web.ControllerKey that the filter sees for each path.
type controllerRecorder struct {
	mu    sync.Mutex
	seen  map[string]string
	found map[string]bool
}

func (r *controllerRecorder) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).ToInstance(&recordFilter{recorder: r})
}

func (r *controllerRecorder) get(path string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.seen[path], r.found[path]
}

// recordFilter reads the controller key from the context it receives. The router puts the key
// in the tag map before it calls the filter chain.
type recordFilter struct {
	recorder *controllerRecorder
}

func (f *recordFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	result := chain.Next(ctx, req, w)

	value, ok := tag.FromContext(ctx).Value(web.ControllerKey)

	f.recorder.mu.Lock()
	f.recorder.seen[req.Request().URL.Path] = value
	f.recorder.found[req.Request().URL.Path] = ok
	f.recorder.mu.Unlock()

	return result
}

// A filter that records its own metrics reads web.ControllerKey to group them by controller.
func TestControllerKey_InTagMap(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	recorder := &controllerRecorder{seen: map[string]string{}, found: map[string]bool{}}

	app := harness.Boot(t, []dingo.Module{new(catalog.Module), recorder}, harness.Config{})

	tests := []struct {
		path string
		want string
	}{
		{path: "/product/1", want: "catalog.product"},
		{path: "/list", want: "catalog.list"},
	}

	for _, tt := range tests {
		resp := app.Get(t, tt.path)

		assert.Equal(t, http.StatusOK, resp.Status, tt.path)

		got, ok := recorder.get(tt.path)

		assert.True(t, ok, "%s: the tag map has no controller key", tt.path)
		assert.Equal(t, tt.want, got, tt.path)
	}
}
