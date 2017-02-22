package router

import (
	"encoding/json"
	"flamingo/core/flamingo/context"
	"flamingo/core/flamingo/service_container"
	"flamingo/core/flamingo/web"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type (
	// Controller defines a web controller
	// it is an interface{} as it can be served by multiple possible controllers,
	// such as generic GET/POST controller, http.Handler, handler-functions, etc.
	Controller interface{}

	// GETController is implemented by controllers which have a Get method
	GETController interface {
		// Get is called for GET-Requests
		Get(web.Context) web.Response
	}

	// POSTController is implemented by controllers which have a Post method
	POSTController interface {
		// Post is called for POST-Requests
		Post(web.Context) web.Response
	}

	// DataController is a controller used to retrieve data, such as user-information, basket
	// etc.
	// By default this will be handled by templates, but there is an out-of-the-box support
	// for JSON requests via /_flamingo/json/{name}, as well as their own route if defined.
	DataController interface {
		// Data is called for data requests
		Data(web.Context) interface{}
	}

	// DataHandler behaves the same as DataController, but just for direct callbacks
	DataHandler func(web.Context) interface{}

	// Router defines the basic Router which is used for holding a context-scoped setup
	// This includes DI resolving etc
	Router struct {
		router   *mux.Router
		routes   map[string]string
		handler  map[string]interface{}
		base     *url.URL
		Logger   *log.Logger `inject:""`
		Sessions sessions.Store
	}
)

// CreateRouter factory for Router
// CreateRouter creates the new flamingo router, set's up handlers and routes and resolved the DI
func CreateRouter(ctx *context.Context, serviceContainer *service_container.ServiceContainer) *Router {
	a := &Router{
		Sessions: sessions.NewCookieStore([]byte("something-very-secret")),
	}

	serviceContainer.Register(a)
	serviceContainer.Resolve()

	// bootstrap
	a.router = mux.NewRouter()
	a.routes = make(map[string]string)
	a.handler = make(map[string]interface{})
	a.base, _ = url.Parse("scheme://" + ctx.BaseUrl)

	// set up routes
	for p, name := range serviceContainer.Routes {
		a.routes[name] = p
	}

	for p, name := range ctx.Routes {
		a.routes[name] = p
	}

	// set up handlers
	for name, handler := range serviceContainer.Handler {
		a.handler[name] = handler
	}

	for name, handler := range ctx.Handler {
		a.handler[name] = handler
	}

	known := make(map[string]bool)

	for name, handler := range a.handler {
		if known[name] {
			continue
		}
		known[name] = true
		route, ok := a.routes[name]
		if !ok {
			continue
		}
		a.Logger.Println("Register", name, "at", route)
		a.router.Handle(route, a.handle(handler)).Name(name)
	}

	return a
}

// Url helps resolving URL's by it's name
// Example:
// 	flamingo.Url("cms.page.view", "name", "Home")
// results in
// 	/baseurl/cms/Home
//
func (router *Router) Url(name string, params ...string) *url.URL {
	if router.router.Get(name) == nil {
		panic("route " + name + " not found")
	}
	u, err := router.router.Get(name).URL(params...)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(router.base.Path, u.Path)
	return u
}

// ServeHTTP shadows the internal mux.Router's ServeHTTP to defer panic recoveries and logging
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w = &ResponseWriter{ResponseWriter: w}
	start := time.Now()
	defer func() {
		var err interface{}
		if err = recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintln(err)))
			w.Write(debug.Stack())
		}
		w.(*ResponseWriter).Log(router.Logger, time.Since(start), req, err)
	}()

	router.router.ServeHTTP(w, req)
}

func (router *Router) handle(c Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s, _ := router.Sessions.Get(req, "aial")

		ctx := web.ContextFromRequest(w, req, s)

		var response web.Response

		switch c := c.(type) {
		case GETController:
			if req.Method == http.MethodGet {
				response = c.Get(ctx)
			}

		case POSTController:
			if req.Method == http.MethodPost {
				response = c.Post(ctx)
			}

		case func(web.Context) web.Response:
			response = c(ctx)

		case DataController:
			response = web.JsonResponse{Data: c.(DataController).Data(ctx)}

		case func(web.Context) interface{}:
			response = web.JsonResponse{Data: c(ctx)}

		case http.Handler:
			c.ServeHTTP(w, req)
			return

		default:
			w.WriteHeader(404)
			w.Write([]byte("404 page not found (no handler)"))
			return
		}

		router.Sessions.Save(req, w, ctx.Session())

		response.Apply(w)
	})
}

// Get is the ServeHTTP's equivalent for DataController and DataHandler
func (router *Router) Get(handler string, ctx web.Context) interface{} {
	if c, ok := router.handler[handler]; ok {
		if c, ok := c.(DataController); ok {
			return c.Data(ctx)
		}
		if c, ok := c.(func(web.Context) interface{}); ok {
			return c(ctx)
		}
		panic("not a data controller")
	} else { // mock...
		data, err := ioutil.ReadFile("frontend/src/mocks/" + handler + ".json")
		if err == nil {
			var res interface{}
			json.Unmarshal(data, &res)
			return res
		} else {
			panic(err)
		}
	}
	panic("not a handler: " + handler)
}
