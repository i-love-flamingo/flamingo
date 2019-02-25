package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/framework/baseurl/application"
)

type (
	// IsExternalURL is exported as a template function
	IsExternalURL struct {
		service *application.Service
	}
)

// Inject dependencies
func (c *IsExternalURL) Inject(service *application.Service) *IsExternalURL {
	c.service = service
	return c
}

// Func returns a boolean if a given URL is external
func (c *IsExternalURL) Func(_ context.Context) interface{} {
	return func(urlStr string) bool {
		if u, err := url.Parse(urlStr); err == nil {
			return c.service.BaseDomain() != u.Host
		}

		return false
	}
}
