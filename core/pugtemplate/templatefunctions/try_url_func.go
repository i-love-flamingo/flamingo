package templatefunctions

import (
	"html/template"
	"net/url"

	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	// URLFunc allows templates to access the routers `URL` helper method
	TryURLFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (u TryURLFunc) Name() string {
	return "tryUrl"
}

// Func as implementation of url method
func (u *TryURLFunc) Func(ctx web.Context) interface{} {
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

		tryUrlResponse, err := u.Router.TryURL(where, p)
		if err != nil {
			return ""
		}

		query := tryUrlResponse.Query()
		for k, v := range q {
			for _, i := range v {
				query.Add(k, i)
			}
		}
		tryUrlResponse.RawQuery = query.Encode()
		return template.URL(tryUrlResponse.String())
	}
}
