package interfaces

import "context"

type (
	// CanonicalDomainFunc is exported as a template function
	CanonicalDomainFunc struct {
		service ApplicationService
	}
)

// Inject CanonicalDomainFunc dependencies
func (c *CanonicalDomainFunc) Inject(service ApplicationService) *CanonicalDomainFunc {
	c.service = service
	return c
}

// Func returns the canonicalDomain func
func (c *CanonicalDomainFunc) Func(context.Context) interface{} {
	return func() string {
		return c.service.GetBaseDomain()
	}
}
