package interfaces

import (
	"flamingo.me/flamingo/core/canonicalUrl/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	// CanonicalUrlFunc is exported as a template function
	CanonicalUrlFunc struct {
		service *application.Service
	}
)

// Inject CanonicalUrlFunc dependencies
func (c *CanonicalUrlFunc) Inject(service *application.Service) {
	c.service = service
}

// Name alias for use in template
func (c *CanonicalUrlFunc) Name() string {
	return "canonicalUrl"
}

// Func returns the CanonicalUrlFunc function
func (c *CanonicalUrlFunc) Func(ctx web.Context) interface{} {
	return func() interface{} {
		return c.service.GetCanonicalUrlForCurrentRequest(ctx)
	}
}
