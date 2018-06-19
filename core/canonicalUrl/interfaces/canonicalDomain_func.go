package interfaces

import (
	"flamingo.me/flamingo/core/canonicalUrl/application"
)

type (
	// CanonicalDomainFunc is exported as a template function
	CanonicalDomainFunc struct {
		Service *application.Service `inject:""`
	}
)

// Name alias for use in template
func (c *CanonicalDomainFunc) Name() string {
	return "canonicalDomain"
}

// Func returns the CSRF NONCE
func (c *CanonicalDomainFunc) Func() interface{} {
	return func() interface{} {
		return c.Service.GetBaseDomain()
	}
}
