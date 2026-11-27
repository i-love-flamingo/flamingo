package config_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

type contextValues struct {
	K1 string `inject:"config:brand.k1"`
	K2 string `inject:"config:brand.k2"`
	K3 string `inject:"config:brand.k3"`
	K4 string `inject:"config:brand.k4"`
}

// contextFile gives a config file that sets the keys from the first one up to k4 to the value.
// The prefix holds other config of the file.
func contextFile(prefix string, first int, value string) string {
	content := prefix + "brand:\n"

	for key := first; key <= 4; key++ {
		content += "  k" + strconv.Itoa(key) + ": " + value + "\n"
	}

	return content
}

// A later file wins.
func TestContext_LoadsContextFilesInOrderThenLocal(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	const areaPath = "flamingo:\n  router:\n    path: en\n"

	files := make(map[string]string)

	// Each area has its own values, to show that an area reads its own files and not the files of root.
	for dir, prefix := range map[string]string{"": "", "en/": areaPath} {
		files[dir+"config.yml"] = contextFile(prefix, 1, dir+"base")
		files[dir+"config_a.yml"] = contextFile("", 2, dir+"a")
		files[dir+"config_b.yml"] = contextFile("", 3, dir+"b")
		files[dir+"config_local.yml"] = contextFile("", 4, dir+"local")
	}

	app := harness.Boot(t,
		[]dingo.Module{probeModule[contextValues](), new(prefixrouter.Module), harness.PassThroughWrapper()},
		harness.Config{
			Env:   map[string]string{"CONTEXT": "a:b"},
			Files: files,
			Areas: []*config.Area{config.NewArea("en", nil)},
		},
	)

	assert.Equal(t, contextValues{K1: "base", K2: "a", K3: "b", K4: "local"}, rootValues[contextValues](t, app), "root")
	assert.Equal(t, contextValues{K1: "en/base", K2: "en/a", K3: "en/b", K4: "en/local"},
		areaValues[contextValues](t, app, "/en"), "en")
}
