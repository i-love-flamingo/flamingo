package services_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/ops"
	"flamingo.me/flamingo/v3/core/gotemplate"
	"flamingo.me/flamingo/v3/framework/web"
)

// funcsModule binds the route "/funcs" that renders the template "funcs".
type funcsModule struct{}

func (*funcsModule) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(funcsRoutes))
}

type funcsRoutes struct {
	responder *web.Responder
}

func (r *funcsRoutes) Inject(responder *web.Responder) {
	r.responder = responder
}

func (r *funcsRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/funcs", "funcs")
	registry.HandleGet("funcs", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Render("funcs", nil)
	})
}

func TestTemplateFunc_Callable(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "funcs.html"), []byte(`{{ upper "shop" }}`), 0o600)
	if err != nil {
		t.Fatalf("write funcs.html: %v", err)
	}

	app := harness.Boot(t, []dingo.Module{new(gotemplate.Module), new(ops.Module), new(funcsModule)}, harness.Config{
		Files: map[string]string{
			"config.yml": "core.gotemplate.engine.templates.basepath: " + dir + "\nflamingo.debug.mode: false\n",
		},
	})

	resp := app.Get(t, "/funcs")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, "SHOP", resp.Body)
}
