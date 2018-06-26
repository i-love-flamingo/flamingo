package interfaces

import (
	"flamingo.me/flamingo/core/canonicalUrl/application"
)

type (
	// CanonicalDomainFunc is exported as a template function
	CanonicalDomainFunc struct {
		service *application.Service
	}
)

// Inject CanonicalDomainFunc dependencies
func (c *CanonicalDomainFunc) Inject(service *application.Service) {
	c.service = service
}

// Name alias for use in template
func (c *CanonicalDomainFunc) Name() string {
	return "canonicalDomain"
}

// Func returns the canonicalDomain func
func (c *CanonicalDomainFunc) Func() interface{} {
	return func() interface{} {
		return c.service.GetBaseDomain()
	}
}
