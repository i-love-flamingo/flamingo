package template_functions

import (
	"flamingo/core/flamingo/router"
	"html/template"
)

type (
	UrlFunc struct {
		Router *router.Router `inject:""`
	}
)

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
			return template.URL(u.Router.Url(where, p...).String())
		}
		return template.URL(u.Router.Url(where).String())
	}
}
