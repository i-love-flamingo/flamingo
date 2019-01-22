package application

import (
	"context"
	"net/url"
	"strings"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// RouterRouter dependency for base url
	RouterRouter interface {
		Base() *url.URL
	}

	// Service exposes helper methods to handle canonical base urls
	Service struct {
		router  RouterRouter
		baseURL string
	}
)

// Inject Service dependencies
func (s *Service) Inject(router RouterRouter, config *struct {
	BaseURL string `inject:"config:canonicalurl.baseurl"`
}) *Service {
	s.router = router
	s.baseURL = config.BaseURL
	return s
}

// GetBaseDomain returns the canonical base domain
func (s *Service) GetBaseDomain() string {
	url, err := url.Parse(s.baseURL)

	if err != nil {
		panic(err)
	}

	return url.Host
}

// GetBaseURL returns the canonical base url
func (s *Service) GetBaseURL() string {
	return strings.TrimRight(s.baseURL, "/")
}

// GetCanonicalURLForCurrentRequest return the canonical url for the current request
// @todo: Add logic to add allowed parameters via controller
func (s *Service) GetCanonicalURLForCurrentRequest(ctx context.Context) string {
	r := web.RequestFromContext(ctx)
	if r == nil {
		return s.GetBaseURL() + s.router.Base().Path
	}
	return s.GetBaseURL() + s.router.Base().Path + r.Request().URL.Path
}
