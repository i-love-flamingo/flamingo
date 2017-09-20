package templatefunctions

import (
	"flamingo/core/pugtemplate/pugjs"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"html/template"
	"strings"
)

type (
	// AssetFunc returns the proper URL for the asset, either local or via CDN
	AssetFunc struct {
		Router *router.Router `inject:""`
		Engine *pugjs.Engine  `inject:""`
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

		if af.Engine.Assetrewrites[assetName] != "" {
			result = url.String() + "/" + af.Engine.Assetrewrites[assetName]
		} else {
			result = url.String() + "/" + string(asset)
		}

		ctx.Push(result, nil) // h2 server push
		return template.URL(result)
	}
}
