package interfaces

import (
	"net/url"

	"flamingo.me/flamingo/v3/framework/baseurl/domain"
)

type (
	// IsExternalURL is exported as a template function
	IsExternalURL struct {
		service domain.Service
	}
)

// Inject dependencies
func (c *IsExternalURL) Inject(service domain.Service) *IsExternalURL {
	c.service = service
	return c
}

// Func returns a boolean if a given URL is external
func (c *IsExternalURL) Func() interface{} {
	return func(urlStr string) bool {
		if u, err := url.Parse(urlStr); err == nil {
			return c.service.BaseDomain() != u.Host
		}

		return false
	}
}
