package interfaces

import (
	"context"

	"flamingo.me/flamingo/core/canonicalUrl/application"
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

// Func returns the CanonicalUrlFunc function
func (c *CanonicalUrlFunc) Func(ctx context.Context) interface{} {
	return func() interface{} {
		return c.service.GetCanonicalUrlForCurrentRequest(ctx)
	}
}
