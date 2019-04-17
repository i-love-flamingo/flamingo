package web

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := &Router{
		eventRouter:    new(flamingo.DefaultEventRouter),
		routesProvider: func() []RoutesModule { return nil },
		filterProvider: func() []Filter { return nil },
	}

	h := router.Handler()

	server := httptest.NewServer(h)
	defer server.Close()

	testReq := func(method, path string) error {
		request, err := http.NewRequest(method, path, nil)
		assert.NoError(t, err)
		request.URL, err = url.Parse(server.URL + path)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		return res.Body.Close()
	}

	var method string

	t.Run("Test Legacy Fallback", func(t *testing.T) {
		registry := NewRegistry()
		h.(*handler).routerRegistry = registry

		_, err := registry.Route("/test", "test")
		assert.NoError(t, err)
		registry.HandleAny("test", func(context.Context, *Request) Result { method = "Handle"; return nil })

		method = ""
		assert.NoError(t, testReq(http.MethodGet, "/test"))
		assert.Equal(t, "Handle", method)
		method = ""
		assert.NoError(t, testReq(http.MethodPost, "/test"))
		assert.Equal(t, "Handle", method)
		method = ""
		assert.NoError(t, testReq(http.MethodHead, "/test"))
		assert.Equal(t, "Handle", method)
		method = ""
		assert.NoError(t, testReq(http.MethodPut, "/test"))
		assert.Equal(t, "Handle", method)
		method = ""
		assert.NoError(t, testReq(http.MethodOptions, "/test"))
		assert.Equal(t, "Handle", method)
		method = ""
		assert.NoError(t, testReq(http.MethodDelete, "/test"))
		assert.Equal(t, "Handle", method)
		method = ""
		assert.NoError(t, testReq("CUSTOM", "/test"))
		assert.Equal(t, "Handle", method)
	})

	t.Run("Test Any", func(t *testing.T) {
		registry := NewRegistry()
		h.(*handler).routerRegistry = registry

		_, err := registry.Route("/test", "test")
		assert.NoError(t, err)
		registry.HandleAny("test", func(context.Context, *Request) Result { method = "HandleAny"; return nil })

		method = ""
		assert.NoError(t, testReq(http.MethodGet, "/test"))
		assert.Equal(t, "HandleAny", method)
		method = ""
		assert.NoError(t, testReq(http.MethodPost, "/test"))
		assert.Equal(t, "HandleAny", method)
		method = ""
		assert.NoError(t, testReq(http.MethodHead, "/test"))
		assert.Equal(t, "HandleAny", method)
		method = ""
		assert.NoError(t, testReq(http.MethodPut, "/test"))
		assert.Equal(t, "HandleAny", method)
		method = ""
		assert.NoError(t, testReq(http.MethodOptions, "/test"))
		assert.Equal(t, "HandleAny", method)
		method = ""
		assert.NoError(t, testReq(http.MethodDelete, "/test"))
		assert.Equal(t, "HandleAny", method)
		method = ""
		assert.NoError(t, testReq("CUSTOM", "/test"))
		assert.Equal(t, "HandleAny", method)
	})

	t.Run("Test HTTP Methods", func(t *testing.T) {
		registry := NewRegistry()
		h.(*handler).routerRegistry = registry

		_, err := registry.Route("/test", "test")
		assert.NoError(t, err)
		registry.HandleGet("test", func(context.Context, *Request) Result { method = "HandleGet"; return nil })
		registry.HandlePost("test", func(context.Context, *Request) Result { method = "HandlePost"; return nil })
		registry.HandlePut("test", func(context.Context, *Request) Result { method = "HandlePut"; return nil })
		registry.HandleOptions("test", func(context.Context, *Request) Result { method = "HandleOptions"; return nil })
		registry.HandleHead("test", func(context.Context, *Request) Result { method = "HandleHead"; return nil })
		registry.HandleDelete("test", func(context.Context, *Request) Result { method = "HandleDelete"; return nil })
		registry.HandleMethod("CUSTOM", "test", func(context.Context, *Request) Result { method = "HandleCustom"; return nil })
		registry.HandleAny("test", func(context.Context, *Request) Result { method = "HandleAny"; return nil })

		method = ""
		assert.NoError(t, testReq(http.MethodGet, "/test"))
		assert.Equal(t, "HandleGet", method)
		method = ""
		assert.NoError(t, testReq(http.MethodPost, "/test"))
		assert.Equal(t, "HandlePost", method)
		method = ""
		assert.NoError(t, testReq(http.MethodHead, "/test"))
		assert.Equal(t, "HandleHead", method)
		method = ""
		assert.NoError(t, testReq(http.MethodPut, "/test"))
		assert.Equal(t, "HandlePut", method)
		method = ""
		assert.NoError(t, testReq(http.MethodOptions, "/test"))
		assert.Equal(t, "HandleOptions", method)
		method = ""
		assert.NoError(t, testReq(http.MethodDelete, "/test"))
		assert.Equal(t, "HandleDelete", method)
		method = ""
		assert.NoError(t, testReq("CUSTOM", "/test"))
		assert.Equal(t, "HandleCustom", method)
		method = ""
		assert.NoError(t, testReq("UNASSIGNED", "/test"))
		assert.Equal(t, "HandleAny", method)
	})
}

func TestRouterTestify(t *testing.T) {
	registry := NewRegistry()
	_, err := registry.Route("/test", "test")
	assert.NoError(t, err)
	registry.HandleAny("test", func(context.Context, *Request) Result { return nil })

	router := &Router{
		eventRouter:    new(flamingo.DefaultEventRouter),
		routesProvider: func() []RoutesModule { return nil },
		filterProvider: func() []Filter { return nil },
	}

	h := router.Handler()
	h.(*handler).routerRegistry = registry

	server := httptest.NewServer(h)
	defer server.Close()

	request, err := http.NewRequest("GET", "/test", nil)
	assert.True(t, err == nil)
	request.URL, err = url.Parse(server.URL + "/test")
	assert.True(t, err == nil)

	defaultClient := &http.Client{}
	res, err := defaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, res.Body.Close())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}

func TestRouterRelativeAndAbsolute(t *testing.T) {
	setupRouter := func(scheme string, host string, path string, external string) *Router {
		registry := NewRegistry()
		router := &Router{}

		router.Inject(&struct {
			Scheme       string         `inject:"config:flamingo.router.scheme,optional"`
			Host         string         `inject:"config:flamingo.router.host,optional"`
			Path         string         `inject:"config:flamingo.router.path,optional"`
			External     string         `inject:"config:flamingo.router.external,optional"`
			SessionStore sessions.Store `inject:",optional"`
		}{
			Scheme:   scheme,
			Host:     host,
			Path:     path,
			External: external,
		}, new(flamingo.DefaultEventRouter), func() []Filter { return nil }, func() []RoutesModule { return nil }, flamingo.NullLogger{}, nil)

		registry.HandleGet("test", func(context.Context, *Request) Result {
			return &Response{}
		})

		_, err := registry.Route("/test", "test")
		assert.NoError(t, err)

		router.routerRegistry = registry

		return router
	}

	t.Run("Test without scheme, without host, without path, without external", func(t *testing.T) {
		router := setupRouter("", "", "", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http:///test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http:///", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http:///test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://flamingo.me/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, without host, without path, without external", func(t *testing.T) {
		router := setupRouter("https", "", "", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https:///test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https:///", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https:///test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://flamingo.me/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, without path, without external", func(t *testing.T) {
		router := setupRouter("https", "other.host", "", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, with path no slashes, without external", func(t *testing.T) {
		router := setupRouter("https", "other.host", "sub-path", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, with path starting slashes, without external", func(t *testing.T) {
		router := setupRouter("https", "other.host", "/sub-path", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, with path ending slashes, without external", func(t *testing.T) {
		router := setupRouter("https", "other.host", "sub-path/", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, with path all slashes, without external", func(t *testing.T) {
		router := setupRouter("https", "other.host", "/sub-path/", "")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/sub-path/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "https://other.host/sub-path/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, with path all slashes, with external no ending slashes", func(t *testing.T) {
		router := setupRouter("https", "other.host", "/sub-path/", "http://external.domain/external-path")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/external-path", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/external-path/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path/test", absoluteUrl.String())
	})

	t.Run("Test with scheme, with host, with path all slashes, with external and ending slashes", func(t *testing.T) {
		router := setupRouter("https", "other.host", "/sub-path/", "http://external.domain/external-path/")

		relativeUrl, err := router.Relative("", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/external-path", relativeUrl.String())

		relativeUrl, err = router.Relative("test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "/external-path/test", relativeUrl.String())

		absoluteUrl, err := router.Absolute(nil, "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path/test", absoluteUrl.String())

		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path/test", absoluteUrl.String())

		req, err = http.NewRequest(http.MethodGet, "https://flamingo.me/flamingo", nil)
		assert.NoError(t, err)
		absoluteUrl, err = router.Absolute(CreateRequest(req, nil), "test", nil)
		assert.NoError(t, err)
		assert.Equal(t, "http://external.domain/external-path/test", absoluteUrl.String())
	})
}
