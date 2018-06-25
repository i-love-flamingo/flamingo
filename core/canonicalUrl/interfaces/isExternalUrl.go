package interfaces

import (
	"net/url"

	"flamingo.me/flamingo/core/canonicalUrl/application"
)

type (
	// IsExternalUrl is exported as a template function
	IsExternalUrl struct {
		Service *application.Service `inject:""`
	}
)

// Name alias for use in template
func (c *IsExternalUrl) Name() string {
	return "isExternalUrl"
}

// Func returns a boolean if a given URL is external
func (c *IsExternalUrl) Func() interface{} {
	return func(urlStr string) bool {
		if url, err := url.Parse(urlStr); err == nil {
			baseUrl := c.Service.GetBaseDomain()
			return baseUrl != url.Host
		}

		return false
	}
}
