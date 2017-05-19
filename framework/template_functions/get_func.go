package template_functions

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
)

type (
	// GetFunc allows templates to access the router's `get` method
	GetFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (g GetFunc) Name() string {
	return "get"
}

// Func as implementation of get method
func (g *GetFunc) Func(ctx web.Context) interface{} {
	return func(what string, params ...map[interface{}]interface{}) interface{} {
		return g.Router.Get(what, ctx, params...)
	}
}
