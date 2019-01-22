package templatefunctions

import (
	"context"
	"html/template"
	"strings"

	"flamingo.me/flamingo/v3/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/v3/framework/router"
)

type (
	// AssetFunc returns the proper URL for the asset, either local or via CDN
	AssetFunc struct {
		Router  *router.Router `inject:""`
		Engine  *pugjs.Engine  `inject:""`
		BaseUrl string         `inject:"config:cdn.base_url,optional"`
	}
)

// Func as implementation of asset method
func (af *AssetFunc) Func(ctx context.Context) interface{} {
	return func(asset pugjs.String) template.URL {
		// let webpack dev server handle URL's
		if af.Engine.Webpackserver {
			return template.URL("/assets/" + asset)
		}

		var result string

		assetSplitted := strings.Split(string(asset), "/")
		assetName := assetSplitted[len(assetSplitted)-1]

		af.Engine.RLock()
		if af.Engine.Assetrewrites[assetName] != "" {
			result = af.Engine.Assetrewrites[assetName]
		} else if af.Engine.Assetrewrites[strings.TrimSpace(string(asset))] != "" {
			result = af.Engine.Assetrewrites[strings.TrimSpace(string(asset))]
		} else {
			result = string(asset)
		}
		af.Engine.RUnlock()

		result = strings.TrimLeft(result, "/")

		result = af.Router.URL("_static", router.P{"n": result}).String()

		if af.BaseUrl != "" {
			baseUrl := strings.TrimRight(af.BaseUrl, "/")
			result = baseUrl + result
		}

		return template.URL(result)
	}
}
