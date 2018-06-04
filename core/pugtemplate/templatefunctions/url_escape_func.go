package templatefunctions

import (
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	url2 "net/url"
)

type (
	// URLFunc allows templates to access the routers `URL` helper method
	URLEscapeFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (u URLEscapeFunc) Name() string {
	return "urlescape"
}

// Func as implementation of url method
func (u *URLEscapeFunc) Func(ctx web.Context) interface{} {
	return func(where string) string {
		return url2.PathEscape(where)
	}
}
