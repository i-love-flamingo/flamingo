package web

import (
	"context"
	"fmt"
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
