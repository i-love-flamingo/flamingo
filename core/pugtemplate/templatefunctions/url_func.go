package templatefunctions

import (
	"html/template"
	"net/url"

	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// URLFunc allows templates to access the routers `URL` helper method
	URLFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (u URLFunc) Name() string {
	return "url"
}

// Func as implementation of url method
func (u *URLFunc) Func(ctx web.Context) interface{} {
	return func(where string, params ...*pugjs.Map) template.URL {
		if where == "" {
			q := ctx.Request().URL.Query()
			if len(params) == 1 {
				for k, v := range params[0].Items {
					q.Del(k.String())
					if arr, ok := v.(*pugjs.Array); ok {
						for _, i := range arr.Items() {
							q.Add(k.String(), i.String())
						}
					} else if v.String() != "" {
						q.Set(k.String(), v.String())
					}
				}
			}
			return template.URL((&url.URL{RawQuery: q.Encode(), Path: u.Router.Base().Path + ctx.Request().URL.Path}).String())
		}

		var p = make(map[string]string)
		var q = make(map[string][]string)
		if len(params) == 1 {
			for k, v := range params[0].Items {
				if arr, ok := v.(*pugjs.Array); ok {
					for _, i := range arr.Items() {
						q[k.String()] = append(q[k.String()], i.String())
					}
				} else {
					p[k.String()] = v.String()
				}
			}
		}
		url := u.Router.URL(where, p)
		query := url.Query()
		for k, v := range q {
			for _, i := range v {
				query.Add(k, i)
			}
		}
		url.RawQuery = query.Encode()
		return template.URL(url.String())
	}
}
