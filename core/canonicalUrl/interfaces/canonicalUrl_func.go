package interfaces

import (
	"context"
)

type (
	// CanonicalUrlFunc is exported as a template function
	CanonicalUrlFunc struct {
		service ApplicationService
	}
)

// Inject CanonicalUrlFunc dependencies
func (c *CanonicalUrlFunc) Inject(service ApplicationService) *CanonicalUrlFunc {
	c.service = service
	return c
}

// Func returns the CanonicalUrlFunc function
func (c *CanonicalUrlFunc) Func(ctx context.Context) interface{} {
	return func() string {
		return c.service.GetCanonicalUrlForCurrentRequest(ctx)
	}
}
