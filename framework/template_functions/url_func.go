package template_functions

import (
	"flamingo/framework/router"
	"html/template"
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
func (u *URLFunc) Func() interface{} {
	return func(where string, params ...map[interface{}]interface{}) template.URL {
		if len(params) > 0 {
			p := make([]string, len(params[0])*2)
			for k, v := range params[0] {
				p = append(p, k.(string), v.(string))
			}
			return template.URL(u.Router.URL(where, p...).String())
		}
		return template.URL(u.Router.URL(where).String())
	}
}
