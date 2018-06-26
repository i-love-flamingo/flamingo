package application

import (
	"net/url"
	"strings"

	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	// Service exposes helper methods to handle canonical base urls
	Service struct {
		router  *router.Router
		baseURL string
	}
)

// Inject Service dependencies
func (s *Service) Inject(router *router.Router, config *struct {
	BaseURL string `inject:"config:canonicalurl.baseurl"`
}) {
	s.router = router
	s.baseURL = config.BaseURL
}

// GetBaseDomain returns the canonical base domain
func (s *Service) GetBaseDomain() string {
	url, err := url.Parse(s.baseURL)

	if err != nil {
		panic(err)
	}

	return url.Host
}

// GetBaseUrl returns the canonical base url
func (s *Service) GetBaseUrl() string {
	return strings.TrimRight(s.baseURL, "/")
}

// GetCanonicalUrlForCurrentRequest return the canonical url for the current request
// @todo: Add logic to add allowed parameters via controller
func (s *Service) GetCanonicalUrlForCurrentRequest(ctx web.Context) string {
	return s.GetBaseUrl() + s.router.Base().Path + ctx.Request().URL.Path
}
