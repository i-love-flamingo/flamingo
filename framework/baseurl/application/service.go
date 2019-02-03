package application

import (
	"net/http"
	"net/url"
	"strings"
)

type (
	// Service to retrieve the base URL
	Service struct {
		baseURL string
	}
)

// Inject dependencies
func (s *Service) Inject(
	cfg *struct {
		BasURL string `inject:"config:baseurl.url"`
	},
) {
	if cfg != nil {
		s.baseURL = cfg.BasURL
	}
}

// BaseURL returns the configured base URL
func (s *Service) BaseURL() string {
	return strings.TrimRight(s.baseURL, "/")
}

// BaseDomain returns the canonical base domain
func (s *Service) BaseDomain() string {
	u, err := url.Parse(s.baseURL)

	if err != nil {
		return ""
	}

	return u.Host
}

// DetermineBase returns the base URL as stated in the request object
func (s *Service) DetermineBase(r *http.Request) string {
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	return strings.TrimRight(scheme+r.Host, "/")
}
