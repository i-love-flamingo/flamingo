package interfaces

import (
	"context"
	"net/url"
)

type (
	// IsExternalURL is exported as a template function
	IsExternalURL struct {
		service ApplicationService
	}
)

// Inject CanonicalURLFunc dependencies
func (c *IsExternalURL) Inject(service ApplicationService) *IsExternalURL {
	c.service = service
	return c
}

// Func returns a boolean if a given URL is external
func (c *IsExternalURL) Func(context.Context) interface{} {
	return func(urlStr string) bool {
		if url, err := url.Parse(urlStr); err == nil {
			return c.service.GetBaseDomain() != url.Host
		}

		return false
	}
}
