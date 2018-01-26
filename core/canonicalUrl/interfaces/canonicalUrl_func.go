package interfaces

import (
	"go.aoe.com/flamingo/core/canonicalUrl/application"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// CanonicalUrlFunc is exported as a template function
	CanonicalUrlFunc struct {
		Service application.Service `inject:""`
	}
)

// Name alias for use in template
func (c *CanonicalUrlFunc) Name() string {
	return "canonicalUrl"
}

// Func returns the CSRF NONCE
func (c *CanonicalUrlFunc) Func(ctx web.Context) interface{} {
	return func() interface{} {
		return c.Service.GetCanonicalUrlForCurrentRequest(ctx)
	}
}
