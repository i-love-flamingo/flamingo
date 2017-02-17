package app

import "flamingo/core/core/app/web"

type GetFunc struct {
	App *App `inject:""`
}

func (_ GetFunc) Name() string {
	return "get"
}

func (g *GetFunc) Func() interface{} {
	return func(ctx web.Context) interface{} {
		return func(what string) interface{} {
			return g.App.Get(what, ctx)
		}
	}
}

type GlobalFunc struct {
	GetFunc `inject:"inline"`
}

func (_ GlobalFunc) Name() string {
	return "global"
}
