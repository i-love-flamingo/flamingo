package interfaces

import (
	"net/url"
)

type (
	// IsExternalUrl is exported as a template function
	IsExternalUrl struct {
		service ApplicationService
	}
)

// Inject CanonicalUrlFunc dependencies
func (c *IsExternalUrl) Inject(service ApplicationService) *IsExternalUrl {
	c.service = service
	return c
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
