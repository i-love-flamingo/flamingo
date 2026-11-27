package routing_test

import (
	"context"
	"net/http"
	"testing"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"flamingo.me/flamingo/v3/framework/web"
)

// areaFiles routes the area "en" at the path "en" and the area "de" at the host "shop.example".
func areaFiles() map[string]string {
	return map[string]string{
		"en/config.yml": "flamingo:\n  router:\n    path: en\n",
		"de/config.yml": "flamingo:\n  router:\n    host: shop.example\n",
	}
}

// bootAreas starts an app with the prefix router. The root area loads the given modules.
func bootAreas(tb testing.TB, root []dingo.Module, areas []*config.Area, files map[string]string) *harness.App {
	tb.Helper()

	root = append([]dingo.Module{new(prefixrouter.Module), harness.PassThroughWrapper()}, root...)

	return harness.Boot(tb, root, harness.Config{Files: files, Areas: areas})
}

// sameModules gives the areas "en" and "de", each with the same area modules.
func sameModules(modules ...dingo.Module) []*config.Area {
	return []*config.Area{
		config.NewArea("en", modules),
		config.NewArea("de", modules),
	}
}

// areaTagModule gives each response of its area the header "X-Area".
type areaTagModule struct {
	name string
}

func (m *areaTagModule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).ToInstance(areaTagFilter{name: m.name})
}

type areaTagFilter struct {
	name string
}

func (f areaTagFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	w.Header().Set("X-Area", f.name)

	return chain.Next(ctx, req, w)
}
