package web

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testController(_ context.Context, _ *Request) Result {
	return nil
}

func TestRegistry(t *testing.T) {
	t.Run("Utils", func(t *testing.T) {
		t.Run("parseHandler", func(t *testing.T) {
			t.Run("should treat empty params properly", func(t *testing.T) {
				var handler = parseHandler("foo.bar")
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Empty(t, handler.params)
			})

			t.Run("should treat params properly", func(t *testing.T) {
				var handler = parseHandler("foo.bar(foo, bar)")
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Len(t, handler.params, 2)
				assert.Equal(t, &param{optional: false, value: ""}, handler.params["foo"])
				assert.Equal(t, &param{optional: false, value: ""}, handler.params["bar"])
			})

			t.Run("should treat optional params properly", func(t *testing.T) {
				var handler = parseHandler("foo.bar(foo?, bar?)")
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Len(t, handler.params, 2)
				assert.Equal(t, &param{optional: true, value: ""}, handler.params["foo"])
				assert.Equal(t, &param{optional: true, value: ""}, handler.params["bar"])
			})

			t.Run("should treat hardcoded params properly", func(t *testing.T) {
				var handler = parseHandler(`foo.bar(foo="bar", x="y")`)
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Len(t, handler.params, 2)
				assert.Equal(t, &param{optional: false, value: "bar"}, handler.params["foo"])
				assert.Equal(t, &param{optional: false, value: "y"}, handler.params["x"])
			})

			t.Run("should treat default value params properly", func(t *testing.T) {
				var handler = parseHandler(`foo.bar(foo?="bar")`)
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Len(t, handler.params, 1)
				assert.Equal(t, &param{optional: true, value: "bar"}, handler.params["foo"])
			})

			t.Run("should treat complexer params properly", func(t *testing.T) {
				var handler = parseHandler(`foo.bar(a, b?, x="a", y ,z, foo ?= "bar")`)
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Len(t, handler.params, 6)
				assert.Equal(t, &param{optional: false, value: ""}, handler.params["a"])
				assert.Equal(t, &param{optional: true, value: ""}, handler.params["b"])
				assert.Equal(t, &param{optional: false, value: "a"}, handler.params["x"])
				assert.Equal(t, &param{optional: false, value: ""}, handler.params["y"])
				assert.Equal(t, &param{optional: false, value: ""}, handler.params["z"])
				assert.Equal(t, &param{optional: true, value: "bar"}, handler.params["foo"])
			})

			t.Run("should treat escaped values properly", func(t *testing.T) {
				var handler = parseHandler(`foo.bar(foo?="\"bar")`)
				assert.Equal(t, "foo.bar", handler.handler)
				assert.Len(t, handler.params, 1)
				assert.Equal(t, &param{optional: true, value: `"bar`}, handler.params["foo"])
			})
		})
	})

	t.Run("API", func(t *testing.T) {
		var registry = NewRegistry()
		t.Run("Should create a new router", func(t *testing.T) {
			assert.NotNil(t, registry)
		})

		t.Run("Should allow controller registration", func(t *testing.T) {
			registry.HandleAny("page.view", testController)
		})

		t.Run("Should make path registration easy", func(t *testing.T) {
			_, err := registry.Route("/", `page.view(page="home")`) // hardcoded
			assert.NoError(t, err)
			_, err = registry.Route("/page/:page", `page.view(page)`) // extract via param
			assert.NoError(t, err)
			_, err = registry.Route("/homepage/:page", `page.view(page?="home2")`) // extract via param, default value
			assert.NoError(t, err)
			_, err = registry.Route("/page", `page.view(page="page")`) // extract from GET
			assert.NoError(t, err)
			_, err = registry.Route("/page2", `page.view(page?="page2")`) // extract from GET
			assert.NoError(t, err)
			_, err = registry.Route("/mustget", `page.view(page)`) // extract from GET
			assert.NoError(t, err)
		})

		t.Run("Should reverse routes properly", func(t *testing.T) {
			p, err := registry.Reverse("page.view", map[string]string{"page": "home"})
			assert.Equal(t, "/", p)
			assert.NoError(t, err)

			p, err = registry.Reverse("page.view", map[string]string{})
			assert.Equal(t, "/homepage/home2", p)
			assert.NoError(t, err)

			p, err = registry.Reverse("page.view", map[string]string{"page": "foo"})
			assert.Equal(t, "/page/foo", p)
			assert.NoError(t, err)
		})

		t.Run("Should match paths", func(t *testing.T) {
			controller, params := registry.match("/homepage/home2")
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "home2", params["page"])

			controller, params = registry.match("/")
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "home", params["page"])

			controller, params = registry.match("/page2")
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "page2", params["page"])

			controller, params = registry.match("/page")
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "page", params["page"])
		})

		t.Run("Should match HTTP Requests", func(t *testing.T) {
			request, _ := http.NewRequest("GET", "/page2", nil)
			controller, params, _ := registry.matchRequest(request)
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "page2", params["page"])

			request, _ = http.NewRequest("GET", "/page2?page=foo", nil)
			controller, params, _ = registry.matchRequest(request)
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "foo", params["page"])

			request, _ = http.NewRequest("GET", "/mustget", nil)
			controller, params, _ = registry.matchRequest(request)
			assert.Equal(t, handlerAction{}, controller)
			assert.Nil(t, params)

			request, _ = http.NewRequest("GET", "/mustget?page=foo", nil)
			controller, params, _ = registry.matchRequest(request)
			assert.NotNil(t, controller)
			assert.Len(t, params, 1)
			assert.Equal(t, "foo", params["page"])
		})

		t.Run("should render get requests if possible", func(t *testing.T) {
			registry.HandleAny("page.get", testController)
			_, err := registry.Route("/path_mustget", `page.get(page)`)
			assert.NoError(t, err)

			registry.HandleAny("page.get2", testController)
			_, err = registry.Route("/path_mustget2", `page.get2(page?="test")`)
			assert.NoError(t, err)

			path, err := registry.Reverse("page.get", map[string]string{"page": "test"})
			assert.NoError(t, err)
			assert.Equal(t, "/path_mustget?page=test", path)

			path, err = registry.Reverse("page.get2", map[string]string{"page": "test"})
			assert.NoError(t, err)
			assert.Equal(t, "/path_mustget2", path)

			path, err = registry.Reverse("page.get2", map[string]string{"page": "nottest"})
			assert.NoError(t, err)
			assert.Equal(t, "/path_mustget2?page=nottest", path)
		})
	})

	t.Run("Catchall", func(t *testing.T) {
		registry := NewRegistry()
		assert.NotNil(t, registry)
		registry.HandleAny("page.view", testController)
		_, err := registry.Route("/page/:page", `page.view(page)`)
		assert.NoError(t, err)
		_, err = registry.Route("/page2/:page", `page.view(page, *")`)
		assert.NoError(t, err)
		_, err = registry.Route("/page3/:page", `page.view`)
		assert.NoError(t, err)

		path, err := registry.Reverse("page.view", map[string]string{"page": "test"})
		assert.NoError(t, err)
		assert.Equal(t, "/page/test", path)

		path, err = registry.Reverse("page.view", map[string]string{"page": "test", "foo": "bar"})
		assert.NoError(t, err)
		assert.Equal(t, "/page2/test?foo=bar", path)

		path, err = registry.Reverse("page.view", map[string]string{"page": "test", "foo": "bar", "x": "y"})
		assert.NoError(t, err)
		assert.Equal(t, "/page2/test?foo=bar&x=y", path)
	})

	t.Run("Enforce Normalization", func(t *testing.T) {
		registry := NewRegistry()
		registry.HandleAny("page.view", testController)
		registry.HandleAny("page2.view", testController)
		route, err := registry.Route("/page/:page", `page.view(page)`)
		assert.NoError(t, err)
		route.Normalize("page")
		_, err = registry.Route("/page2/:page", `page2.view(page)`)
		assert.NoError(t, err)

		path, err := registry.Reverse("page.view", map[string]string{"page": "Test & 123 - test"})
		assert.NoError(t, err)
		assert.Equal(t, "/page/test-%26-123-test", path)

		path, err = registry.Reverse("page2.view", map[string]string{"page": "Test & 123 - test"})
		assert.NoError(t, err)
		assert.Equal(t, "/page2/Test+%26+123+-+test", path)
	})
}
