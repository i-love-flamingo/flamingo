package application

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type (
	// Service to retrieve the base URL
	Service struct {
		baseURL string
		scheme  string
	}
)

var scheme = regexp.MustCompile("^https?://")

// Inject dependencies
func (s *Service) Inject(
	cfg *struct {
		BaseURL string `inject:"config:baseurl.url"`
		Scheme  string `inject:"config:baseurl.scheme"`
	},
) {
	if cfg != nil {
		s.baseURL = cfg.BaseURL
		s.scheme = cfg.Scheme
	}
}

// BaseURL returns the configured base URL
func (s *Service) BaseURL() string {
	baseURL := s.baseURL
	// prepend configured scheme if url is not relative and has no scheme itself
	if !strings.HasPrefix(baseURL, "/") && !scheme.MatchString(baseURL) {
		baseURL = s.scheme + baseURL
	}

	return strings.TrimRight(baseURL, "/")
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
	scheme := s.scheme
	// try to fall back if no scheme is configured
	if scheme == "" {
		scheme = "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
	}

	return strings.TrimRight(scheme+r.Host, "/")
}
