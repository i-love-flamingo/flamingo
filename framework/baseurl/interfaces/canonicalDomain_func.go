package interfaces

import (
	"context"

	"flamingo.me/flamingo/v3/framework/baseurl/domain"
)

type (
	// CanonicalDomainFunc is exported as a template function
	CanonicalDomainFunc struct {
		service domain.Service
	}
)

// Inject dependencies
func (c *CanonicalDomainFunc) Inject(service domain.Service) *CanonicalDomainFunc {
	c.service = service
	return c
}

// Func returns the canonicalDomain func
func (c *CanonicalDomainFunc) Func(_ context.Context) interface{} {
	return func() string {
		return c.service.BaseDomain()
	}
}
