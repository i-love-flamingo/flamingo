package bootstrap_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"flamingo.me/flamingo/v3/framework/web"
)

// areaFiles routes the areas "en" and "de" at the prefixes "en" and "de".
func areaFiles() map[string]string {
	return map[string]string{
		"en/config.yml": "flamingo:\n  router:\n    path: en\n",
		"de/config.yml": "flamingo:\n  router:\n    path: de\n",
	}
}

// counterModule counts how often dingo calls Configure.
type counterModule struct {
	configured *atomic.Int32
}

func (m *counterModule) Configure(*dingo.Injector) {
	m.configured.Add(1)
}

type brandName string

// brandModule binds a value in the area where it is loaded and a route that returns it with the area name.
type brandModule struct {
	configured *atomic.Int32
}

func (m *brandModule) Configure(injector *dingo.Injector) {
	m.configured.Add(1)

	injector.Bind(new(brandName)).ToInstance(brandName("brand"))
	web.BindRoutes(injector, new(brandRoutes))
}

type brandRoutes struct {
	brand brandName
	area  string
}

func (r *brandRoutes) Inject(brand brandName, cfg *struct {
	Area string `inject:"config:area"`
}) {
	r.brand = brand
	r.area = cfg.Area
}

func (r *brandRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/brand", "brand")
	registry.HandleGet("brand", func(_ context.Context, _ *web.Request) web.Result {
		return new(web.Responder).HTTP(http.StatusOK, strings.NewReader(string(r.brand)+" in "+r.area))
	})
}

func TestArea_ConfigureRunsOncePerOwnModuleList(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	var configured atomic.Int32

	counter := &counterModule{configured: &configured}

	harness.Boot(t,
		[]dingo.Module{counter, new(prefixrouter.Module), harness.PassThroughWrapper()},
		harness.Config{
			Files: areaFiles(),
			Areas: []*config.Area{
				config.NewArea("en", []dingo.Module{counter}),
				config.NewArea("de", nil),
			},
		},
	)

	assert.EqualValues(t, 2, configured.Load(), "root and en configure the module, de does not list it")
}

func TestArea_ChildInheritsParentBindings(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	var configured atomic.Int32

	app := harness.Boot(t,
		[]dingo.Module{&brandModule{configured: &configured}, new(prefixrouter.Module), harness.PassThroughWrapper()},
		harness.Config{
			Files: areaFiles(),
			Areas: []*config.Area{
				config.NewArea("en", nil),
				config.NewArea("de", nil),
			},
		},
	)

	body := app.Get(t, "/en/brand").Body

	assert.Equal(t, "brand in root/en", body)
	assert.EqualValues(t, 1, configured.Load())
}
