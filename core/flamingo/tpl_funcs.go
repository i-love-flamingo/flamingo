package flamingo

import (
	"flamingo/core/flamingo/web"
	"html/template"
)

type (
	GetFunc struct {
		App *Router `inject:""`
	}

	UrlFunc struct {
		App *Router `inject:""`
	}
)

func (_ GetFunc) Name() string {
	return "get"
}

func (g *GetFunc) Func(ctx web.Context) interface{} {
	return func(what string) interface{} {
		return g.App.Get(what, ctx)
	}
}

func (_ UrlFunc) Name() string {
	return "url"
}

func (u *UrlFunc) Func() interface{} {
	return func(where string, params ...map[interface{}]interface{}) template.URL {
		if len(params) > 0 {
			p := make([]string, len(params[0])*2)
			for k, v := range params[0] {
				p = append(p, k.(string), v.(string))
			}
			return template.URL(u.App.Url(where, p...).String())
		}
		return template.URL(u.App.Url(where).String())
	}
}
