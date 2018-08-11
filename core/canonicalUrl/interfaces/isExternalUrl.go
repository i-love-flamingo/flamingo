package interfaces

import (
	"net/url"

	"flamingo.me/flamingo/core/canonicalUrl/application"
)

type (
	// IsExternalUrl is exported as a template function
	IsExternalUrl struct {
		service *application.Service
	}
)

// Inject CanonicalUrlFunc dependencies
func (c *IsExternalUrl) Inject(service *application.Service) {
	c.service = service
}

// Func returns a boolean if a given URL is external
func (c *IsExternalUrl) Func() interface{} {
	return func(urlStr string) bool {
		if url, err := url.Parse(urlStr); err == nil {
			return c.service.GetBaseDomain() != url.Host
		}

		return false
	}
}
