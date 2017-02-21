package template

import (
	"encoding/json"
	"flamingo/core/flamingo"
	"flamingo/core/flamingo/web"
	"html/template"
	"strings"
)

type (
	AssetFunc struct {
		Router *flamingo.Router   `inject:""`
		Engine *PugTemplateEngine `inject:""`
	}
	DebugFunc struct{}
)

func (_ AssetFunc) Name() string {
	return "asset"
}

func (af *AssetFunc) Func(ctx web.Context) interface{} {
	return func(asset string) template.URL {
		// let webpack devserver handle URL's
		if af.Engine.webpackserver {
			return template.URL("/assets/" + asset)
		}

		// get the _static URL
		url := af.Router.Url("_static", "n", "")
		var result string

		assetSplitted := strings.Split(asset, "/")
		assetName := assetSplitted[len(assetSplitted)-1]

		if af.Engine.assetrewrites[assetName] != "" {
			result = url.String() + "/" + af.Engine.assetrewrites[assetName]
		} else {
			result = url.String() + "/" + asset
		}

		ctx.Push(result, nil) // h2 server push
		return template.URL(result)
	}
}

func (_ DebugFunc) Name() string {
	return "debug"
}

func (_ DebugFunc) Func() interface{} {
	return func(o interface{}) string {
		d, _ := json.MarshalIndent(o, "", "    ")
		return string(d)
	}
}
