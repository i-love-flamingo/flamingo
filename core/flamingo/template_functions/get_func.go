package template_functions

import (
	"flamingo/core/flamingo/router"
	"flamingo/core/flamingo/web"
)

type (
	GetFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (_ GetFunc) Name() string {
	return "get"
}

// Func as implementation of get method
func (g *GetFunc) Func(ctx web.Context) interface{} {
	return func(what string) interface{} {
		return g.Router.Get(what, ctx)
	}
}
