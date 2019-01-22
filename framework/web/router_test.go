package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	router := new(Router)

	router.eventrouter = new(flamingo.DefaultEventRouter)

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
		registry := NewRegistry()
		router.routerRegistry = registry

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
		router.routerRegistry = registry

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
		router.routerRegistry = registry

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
	router := new(Router)

	registry := NewRegistry()
	_, err := registry.Route("/test", "test")
	assert.NoError(t, err)
	registry.HandleAny("test", func(context.Context, *Request) Result { return nil })
	router.routerRegistry = registry

	router.eventrouter = new(flamingo.DefaultEventRouter)

	server := httptest.NewServer(router)
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

func TestRouterMiniMocks(t *testing.T) {
	routerVar := new(Router)
	routerVar.eventrouter = new(flamingo.DefaultEventRouter)

	registry := NewRegistry()
	_, err := registry.Route("/test", "test")
	assert.NoError(t, err)
	registry.HandleAny("test", func(context.Context, *Request) Result { return nil })

	routerVar.routerRegistry = registry

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
	assert.NoError(t, res.Body.Close())
	require.NoError(t, err)

	fmt.Printf("%s", greeting)
}

func TestRouterTimeout(t *testing.T) {
	tests := []struct {
		name            string
		exceededTimeout float64
		controller      func(context.Context, *Request) Result
	}{
		{
			name:            "Timeout enforced",
			exceededTimeout: float64(10),
			controller: testControllerFactory(t, 18*time.Millisecond, func(t *testing.T, ctx context.Context) {
				t.Helper()
				select {
				case <-ctx.Done():
				default:
					t.Error("Timeout was not caught")
				}
			}),
		},
		{
			name:            "Timeout not enforced",
			exceededTimeout: float64(10),
			controller: testControllerFactory(t, 0, func(t *testing.T, ctx context.Context) {
				t.Helper()
				select {
				case <-ctx.Done():
					t.Error("Timeout was caught but shouldn't have")
				default:
				}
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routerVar := &Router{
				routerTimeout: tt.exceededTimeout,
			}

			routerVar.eventrouter = new(flamingo.DefaultEventRouter)

			registry := NewRegistry()
			_, err := registry.Route("/test", "test")
			assert.NoError(t, err)

			registry.HandleAny("test", tt.controller)

			routerVar.routerRegistry = registry

			server := httptest.NewServer(routerVar)
			defer server.Close()

			request, err := http.NewRequest("GET", "/test", nil)

			if err != nil {
				t.Fatal(err)
			}

			responseRecorder := httptest.NewRecorder()

			routerVar.ServeHTTP(responseRecorder, request)
		})
	}
}

func testControllerFactory(t *testing.T, timeout time.Duration, validator func(*testing.T, context.Context)) func(context.Context, *Request) Result {
	t.Helper()
	return func(ctx context.Context, _ *Request) Result {
		time.Sleep(timeout)

		validator(t, ctx)

		return nil
	}
}
