package template_functions

import (
	"flamingo/core/flamingo/router"
	"flamingo/core/flamingo/web"
	"flamingo/core/packages/pug_template/pugast"
	"html/template"
	"strings"
)

type (
	AssetFunc struct {
		Router *router.Router            `inject:""`
		Engine *pugast.PugTemplateEngine `inject:""`
	}
)

// Name alias for use in template
func (_ AssetFunc) Name() string {
	return "asset"
}

// Func as implementation of asset method
func (af *AssetFunc) Func(ctx web.Context) interface{} {
	return func(asset string) template.URL {
		// let webpack dev server handle URL's
		if af.Engine.Webpackserver {
			return template.URL("/assets/" + asset)
		}

		// get the _static URL
		url := af.Router.Url("_static", "n", "")
		var result string

		assetSplitted := strings.Split(asset, "/")
		assetName := assetSplitted[len(assetSplitted)-1]

		if af.Engine.Assetrewrites[assetName] != "" {
			result = url.String() + "/" + af.Engine.Assetrewrites[assetName]
		} else {
			result = url.String() + "/" + asset
		}

		ctx.Push(result, nil) // h2 server push
		return template.URL(result)
	}
}
