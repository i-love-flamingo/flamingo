package application

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"
)

type mockRouter struct{ url string }

func (m *mockRouter) Base() *url.URL {
	u, _ := url.Parse(m.url)
	return u
}

func testService(url string) *Service {
	return new(Service).Inject(&mockRouter{url: url}, &struct {
		BaseURL string `inject:"config:canonicalurl.baseurl"`
	}{BaseURL: url})
}

func TestService_GetBaseDomain(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"http://example.com/path/", "example.com"},
		{"http://example.com", "example.com"},
	}
	for _, tt := range tests {
		s := testService(tt.url)
		if got := s.GetBaseDomain(); got != tt.want {
			t.Errorf("Service.GetBaseDomain() = %v, want %v", got, tt.want)
		}
	}
}

func TestService_GetBaseUrl(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"http://example.com/path/", "http://example.com/path"},
		{"http://example.com", "http://example.com"},
	}
	for _, tt := range tests {
		s := testService(tt.url)
		if got := s.GetBaseURL(); got != tt.want {
			t.Errorf("Service.GetBaseDomain() = %v, want %v", got, tt.want)
		}
	}
}

func TestService_GetCanonicalUrlForCurrentRequest(t *testing.T) {
	tests := []struct {
		ctx  bool
		path string
		base string
		want string
	}{
		{true, "/", "http://example.com/path/", "http://example.com/path/path//"},
		{true, "/", "http://example.com", "http://example.com/"},
		{true, "/path/", "http://example.com/path/", "http://example.com/path/path//path/"},
		{true, "/foo", "http://example.com", "http://example.com/foo"},

		{false, "/", "http://example.com/path/", "http://example.com/path/path/"},
		{false, "/", "http://example.com", "http://example.com"},
		{false, "/path/", "http://example.com/path/", "http://example.com/path/path/"},
		{false, "/foo", "http://example.com", "http://example.com"},
	}
	for _, tt := range tests {
		s := testService(tt.base)
		ctx := context.Background()
		if tt.ctx {
			req, _ := http.NewRequest(http.MethodGet, tt.path, nil)
			ctx = web.ContextWithRequest(ctx, web.CreateRequest(req, nil))
		}

		if got := s.GetCanonicalURLForCurrentRequest(ctx); got != tt.want {
			t.Errorf("Service.GetBaseDomain() = %v, want %v (context: %v)", got, tt.want, tt.ctx)
		}
	}
}
