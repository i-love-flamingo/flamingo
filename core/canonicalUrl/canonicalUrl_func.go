package canonicalUrl

import (
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// CanonicalUrlFunc is exported as a template function
	CanonicalUrlFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (c *CanonicalUrlFunc) Name() string {
	return "canonicalUrl"
}

// Func returns the CSRF NONCE
func (c *CanonicalUrlFunc) Func(ctx web.Context) interface{} {
	return func() interface{} {
		// @todo: Add host
		// @todo: Add logic to add allowed parameters in controller
		url := c.Router.Base().Path + ctx.Request().URL.Path
		return url
	}
}
