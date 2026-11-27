package responder_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/catalog"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// partialEngine is a template engine that renders a full page and fails to render partials.
// partialCalls counts the partial renders.
type partialEngine struct {
	partialCalls *atomic.Int32
}

var _ flamingo.PartialTemplateEngine = partialEngine{}

func (partialEngine) Render(_ context.Context, _ string, _ interface{}) (io.Reader, error) {
	return strings.NewReader("full-page"), nil
}

func (e partialEngine) RenderPartials(_ context.Context, _ string, _ interface{}, _ []string) (map[string]io.Reader, error) {
	e.partialCalls.Add(1)

	return nil, errors.New("catalog-secret-error")
}

// partialEngineModule binds partialEngine as the template engine.
type partialEngineModule struct {
	engine partialEngine
}

func (m *partialEngineModule) Configure(injector *dingo.Injector) {
	injector.Bind(new(flamingo.TemplateEngine)).ToInstance(m.engine)
}

func TestTemplateError_FailedRenderShowsErrorPage(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootResponder(t)

	resp := app.Get(t, "/broken")

	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Contains(t, resp.Body, "error-with-code")
}

func TestTemplateError_PartialRenderHidesEngineError(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	engine := partialEngine{partialCalls: new(atomic.Int32)}

	app := harness.Boot(t, []dingo.Module{&partialEngineModule{engine: engine}, new(catalog.Module)}, harness.Config{
		Files: map[string]string{"config.yml": configYAML()},
	})

	resp := app.Get(t, "/broken", "X-Partial", "content")

	assert.Positive(t, engine.partialCalls.Load(), "the partial render path was not used")
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.NotContains(t, resp.Body, "catalog-secret-error")
}
