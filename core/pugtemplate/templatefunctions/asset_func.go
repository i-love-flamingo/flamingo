package templatefunctions

import (
	"html/template"
	"strings"

	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	// AssetFunc returns the proper URL for the asset, either local or via CDN
	AssetFunc struct {
		Router  *router.Router `inject:""`
		Engine  *pugjs.Engine  `inject:""`
		BaseUrl string         `inject:"config:cdn.base_url,optional"`
	}
)

// Name alias for use in template
func (af AssetFunc) Name() string {
	return "asset"
}

// Func as implementation of asset method
func (af *AssetFunc) Func(ctx web.Context) interface{} {
	return func(asset pugjs.String) template.URL {
		// let webpack dev server handle URL's
		if af.Engine.Webpackserver {
			return template.URL("/assets/" + asset)
		}

		// get the _static URL
		url := af.Router.URL("_static", router.P{"n": ""})
		var result string

		assetSplitted := strings.Split(string(asset), "/")
		assetName := assetSplitted[len(assetSplitted)-1]

		af.Engine.Lock()
		if af.Engine.Assetrewrites[assetName] != "" {
			result = url.String() + "/" + af.Engine.Assetrewrites[assetName]
		} else if af.Engine.Assetrewrites[strings.TrimSpace(string(asset))] != "" {
			result = url.String() + "/" + af.Engine.Assetrewrites[strings.TrimSpace(string(asset))]
		} else {
			result = url.String() + "/" + string(asset)
		}
		af.Engine.Unlock()

		result = strings.Replace(result, "//", "/", -1)

		baseUrl := strings.TrimRight(af.BaseUrl, "/")
		if baseUrl != "" {
			result = baseUrl + result
		}

		ctx.Push(result, nil) // h2 server push
		return template.URL(result)
	}
}
