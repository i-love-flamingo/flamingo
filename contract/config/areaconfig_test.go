package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

type brandAreaValues struct {
	Area  string `inject:"config:area"`
	Name  string `inject:"config:brand.name"`
	Color string `inject:"config:brand.color"`
}

// bootEnglishArea starts an app with the probe, the area "en" and the given config files.
func bootEnglishArea(tb testing.TB, files map[string]string) *harness.App {
	tb.Helper()

	return harness.Boot(tb,
		[]dingo.Module{probeModule[brandAreaValues](), new(prefixrouter.Module), harness.PassThroughWrapper()},
		harness.Config{
			Files: files,
			Areas: []*config.Area{config.NewArea("en", nil)},
		},
	)
}

func TestAreaConfig_ChildInheritsAndOverridesRootValues(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootEnglishArea(t, map[string]string{
		"config.yml":    "brand:\n  name: shop\n  color: red\n",
		"en/config.yml": "flamingo:\n  router:\n    path: en\nbrand:\n  color: blue\n",
	})

	got := areaValues[brandAreaValues](t, app, "/en")

	assert.Equal(t, "shop", got.Name)
	assert.Equal(t, "blue", got.Color)
}

// Logs and metrics use the "area" value as a label.
func TestAreaConfig_AreaKeyIsSlashPath(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootEnglishArea(t, map[string]string{
		"config.yml":    "brand:\n  name: shop\n  color: red\n",
		"en/config.yml": "flamingo:\n  router:\n    path: en\n",
	})

	assert.Equal(t, "root", rootValues[brandAreaValues](t, app).Area)
	assert.Equal(t, "root/en", areaValues[brandAreaValues](t, app, "/en").Area)
}
