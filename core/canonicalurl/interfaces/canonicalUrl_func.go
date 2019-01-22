package interfaces

import (
	"context"
)

type (
	// CanonicalURLFunc is exported as a template function
	CanonicalURLFunc struct {
		service ApplicationService
	}
)

// Inject CanonicalURLFunc dependencies
func (c *CanonicalURLFunc) Inject(service ApplicationService) *CanonicalURLFunc {
	c.service = service
	return c
}

// Func returns the CanonicalURLFunc function
func (c *CanonicalURLFunc) Func(ctx context.Context) interface{} {
	return func() string {
		return c.service.GetCanonicalURLForCurrentRequest(ctx)
	}
}
