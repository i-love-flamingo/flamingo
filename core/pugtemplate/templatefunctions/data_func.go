package templatefunctions

import (
	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	// DataFunc allows templates to access the router's `get` method
	DataFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (g DataFunc) Name() string {
	return "data"
}

// Func as implementation of get method
func (g *DataFunc) Func(ctx web.Context) interface{} {
	return func(what string, params ...*pugjs.Map) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0].Items {
				p[k.String()] = v.String()
			}
		}
		return g.Router.Get(what, ctx, p)
	}
}
