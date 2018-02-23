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

	"github.com/gojuno/minimock"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.aoe.com/flamingo/framework/event"
	eventMocks "go.aoe.com/flamingo/framework/event/mocks"
	"go.aoe.com/flamingo/framework/profiler"
	profilerMocks "go.aoe.com/flamingo/framework/profiler/mocks"
	"go.aoe.com/flamingo/framework/web"
	webMocks "go.aoe.com/flamingo/framework/web/mocks"
)

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

	profilerMock := new(profilerMocks.Profiler)
	router.ProfilerProvider = func() profiler.Profiler { return profilerMock }

	registry := NewRegistry()
	registry.Route("/test", "test")
	registry.Handle("test", func(context web.Context) web.Response { return nil })
	router.RouterRegistry = registry

	contextMock := new(webMocks.Context)
	router.ContextFactory = func(p profiler.Profiler, e event.Router, rw http.ResponseWriter, r *http.Request, session *sessions.Session) web.Context {
		contextMock.On("Profile", "matchRequest", "/test").Return(profiler.ProfileFinishFunc(func() {}))
		contextMock.On("Profile", "request", "/test").Return(profiler.ProfileFinishFunc(func() {}))
		contextMock.On("LoadParams", mock.Anything)
		contextMock.On(
			"WithValue",
			mock.MatchedBy(func(key interface{}) bool { return key.(string) == "Handler" }),
			mock.Anything,
		).Return(contextMock)
		contextMock.On(
			"WithValue",
			mock.MatchedBy(func(key interface{}) bool { return key.(string) == "HandlerName" }),
			mock.Anything,
		).Return(contextMock)
		return contextMock
	}
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

	profilerMock := NewProfilerMock(t)
	routerVar.ProfilerProvider = func() profiler.Profiler { return profilerMock }

	registry := NewRegistry()
	registry.Route("/test", "test")
	registry.Handle("test", func(context web.Context) web.Response { return nil })

	routerVar.RouterRegistry = registry

	tester := minimock.NewController(t)

	contextMock := NewContextMock(tester)
	contextMock.ProfileFunc = func(p string, p1 string) (r profiler.ProfileFinishFunc) { return profiler.ProfileFinishFunc(func() {}) }
	contextMock.LoadParamsMock.Expect(map[string]string{}).Return()
	contextMock.WithValueFunc = func(p interface{}, p1 interface{}) (r web.Context) { return nil }
	routerVar.ContextFactory = func(profiler profiler.Profiler, eventrouter event.Router, rw http.ResponseWriter, r *http.Request, session *sessions.Session) web.Context {
		return contextMock
	}

	eventRouter := NewRouterMock(t)
	eventRouter.DispatchFunc = func(ctx context.Context, p event.Event) {}
	routerVar.eventrouter = eventRouter

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

	assert.True(t, contextMock.AllMocksCalled())
}
