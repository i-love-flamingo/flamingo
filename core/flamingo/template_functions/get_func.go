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

func (_ GetFunc) Name() string {
	return "get"
}

func (g *GetFunc) Func(ctx web.Context) interface{} {
	return func(what string) interface{} {
		return g.Router.Get(what, ctx)
	}
}
