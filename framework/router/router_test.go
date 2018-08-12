package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	eventMocks "flamingo.me/flamingo/framework/event/mocks"
	"flamingo.me/flamingo/framework/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	router := new(Router)

	registry := NewRegistry()
	router.RouterRegistry = registry

	eventRouter := new(eventMocks.Router)
	eventRouter.On("Dispatch", mock.Anything, mock.Anything)
	router.eventrouter = eventRouter

	server := httptest.NewServer(router)
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
		registry.Route("/test", "test")
		registry.HandleAny("test", func(context context.Context, req *web.Request) web.Response { method = "Handle"; return nil })

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
		registry.Route("/test", "test")
		registry.HandleAny("test", func(context context.Context, req *web.Request) web.Response { method = "HandleAny"; return nil })

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
		registry.Route("/test", "test")
		registry.HandleGet("test", func(context context.Context, req *web.Request) web.Response { method = "HandleGet"; return nil })
		registry.HandlePost("test", func(context context.Context, req *web.Request) web.Response { method = "HandlePost"; return nil })
		registry.HandlePut("test", func(context context.Context, req *web.Request) web.Response { method = "HandlePut"; return nil })
		registry.HandleOptions("test", func(context context.Context, req *web.Request) web.Response { method = "HandleOptions"; return nil })
		registry.HandleHead("test", func(context context.Context, req *web.Request) web.Response { method = "HandleHead"; return nil })
		registry.HandleDelete("test", func(context context.Context, req *web.Request) web.Response { method = "HandleDelete"; return nil })
		registry.HandleMethod("CUSTOM", "test", func(context context.Context, req *web.Request) web.Response { method = "HandleCustom"; return nil })

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

// Example
func TestTest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}

func TestRouterTestify(t *testing.T) {
	router := new(Router)

	registry := NewRegistry()
	registry.Route("/test", "test")
	registry.HandleAny("test", func(context context.Context, req *web.Request) web.Response { return nil })
	router.RouterRegistry = registry

	eventRouter := new(eventMocks.Router)
	eventRouter.On("Dispatch", mock.Anything, mock.Anything)
	router.eventrouter = eventRouter

	server := httptest.NewServer(router)
	defer server.Close()

	request, err := http.NewRequest("GET", "/test", nil)
	assert.True(t, err == nil)
	request.URL, err = url.Parse(server.URL + "/test")
	assert.True(t, err == nil)

	defaultClient := &http.Client{}
	res, err := defaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}

func TestRouterMiniMocks(t *testing.T) {
	routerVar := new(Router)
	routerVar.eventrouter = new(eventMocks.Router)
	routerVar.eventrouter.(*eventMocks.Router).On("Dispatch", mock.Anything, mock.Anything).Return()

	registry := NewRegistry()
	registry.Route("/test", "test")
	registry.HandleAny("test", func(context context.Context, req *web.Request) web.Response { return nil })

	routerVar.RouterRegistry = registry

	server := httptest.NewServer(routerVar)
	defer server.Close()

	request, err := http.NewRequest("GET", "/test", nil)
	assert.True(t, err == nil)
	request.URL, err = url.Parse(server.URL + "/test")
	assert.True(t, err == nil)

	defaultClient := &http.Client{}
	res, err := defaultClient.Do(request)
	require.NoError(t, err)

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)

	fmt.Printf("%s", greeting)
}
