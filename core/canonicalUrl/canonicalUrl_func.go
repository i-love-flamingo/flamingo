package canonicalUrl

import (
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"strings"
)

type (
	// CanonicalUrlFunc is exported as a template function
	CanonicalUrlFunc struct {
		Router  *router.Router `inject:""`
		BaseUrl string         `inject:"config:canonicalurl.baseurl"`
	}
)

// Name alias for use in template
func (c *CanonicalUrlFunc) Name() string {
	return "canonicalUrl"
}

// Func returns the canonical URL
func (c *CanonicalUrlFunc) Func(ctx web.Context) interface{} {
	baseUrl := strings.TrimRight(c.BaseUrl, "/")

	return func() interface{} {
		// @todo: Add logic to add allowed parameters via controller
		url := baseUrl + c.Router.Base().Path + ctx.Request().URL.Path
		return url
	}
}
