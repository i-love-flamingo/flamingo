package router

import (
	"context"
	"encoding/json"
	configcontext "flamingo/core/flamingo/context"
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/flamingo/web"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	RouterRegister = "router.register"
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
		router            *mux.Router
		routes            map[string]string
		handler           map[string]Controller
		hardroutes        map[string]configcontext.Route
		hardroutesreverse map[string]configcontext.Route
		base              *url.URL

		Logger           *log.Logger        `inject:""` // Logger is a default logger for now
		Sessions         sessions.Store     `inject:""` // Sessions storage, which are used to retrieve user-context session
		ServiceContainer *di.Container      `inject:""` // ServiceContainer holder
		ContextFactory   web.ContextFactory `inject:""` // ContextFactory for new contexts
	}
)

// NewCookieStore because vendor-folder are hard...
func NewCookieStore(secret []byte) *sessions.CookieStore {
	return sessions.NewCookieStore(secret)
}

// CreateRouter factory for Routers.
// Creates the new flamingo router, set's up handlers and routes and resolved the DI.
// BUG(bastian.ike) hardroutesreverse style is borked
func CreateRouter(ctx *configcontext.Context, serviceContainer *di.Container) *Router {
	router := &Router{
		router:            mux.NewRouter(),
		routes:            make(map[string]string),
		hardroutes:        make(map[string]configcontext.Route),
		hardroutesreverse: make(map[string]configcontext.Route),
		handler:           make(map[string]Controller),
	}

	router.base, _ = url.Parse("scheme://" + ctx.BaseURL)

	serviceContainer.Register(router)
	serviceContainer.Resolve(router)

	for _, route := range ctx.Routes {
		if route.Args == nil {
			router.routes[route.Controller] = route.Path
		} else {
			router.hardroutes[route.Path] = route
			p := make([]string, len(route.Args)*2)
			for k, v := range route.Args {
				p = append(p, k, v)
			}
			router.hardroutesreverse[route.Controller+strings.Join(p, "!!")] = route
		}
	}

	for name, handler := range router.handler {
		if route, ok := router.routes[name]; ok {
			router.router.Handle(route, router.handle(handler)).Name(name)
		}
	}

	return router
}

// Handle registers the controller for a named route
func (router *Router) Handle(name string, controller Controller) {
	router.handler[name] = controller
	router.ServiceContainer.Resolve(controller)
}

// Router registers the path for a named route
func (router *Router) Route(path, name string) {
	router.routes[name] = path
}

// CompilerPass to get router register functions
func (router *Router) CompilerPass(c *di.Container) {
	for _, registerFunc := range c.GetTagged(RouterRegister) {
		registerFunc.Value.(func(r *Router))(router)
	}
}

// URL helps resolving URL's by it's name.
// Example:
//     flamingo.URL("cms.page.view", "name", "Home")
// results in
//     /baseurl/cms/Home
func (router *Router) URL(name string, params ...string) *url.URL {
	var resultURL *url.URL
	if route, ok := router.hardroutesreverse[name+`!!!!`+strings.Join(params, "!!")]; ok {
		resultURL, _ = url.Parse(route.Path)
	} else {
		resultURL = router.url(name, params...)
	}

	resultURL.Path = path.Join(router.base.Path, resultURL.Path)

	return resultURL
}

// url builds a URL for a Router.
func (router *Router) url(name string, params ...string) *url.URL {
	if router.router.Get(name) == nil {
		panic("route " + name + " not found")
	}

	resultURL, err := router.router.Get(name).URL(params...)

	if err != nil {
		panic(err)
	}

	return resultURL
}

// ServeHTTP shadows the internal mux.Router's ServeHTTP to defer panic recoveries and logging.
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// shadow the response writer
	w = &VerboseResponseWriter{ResponseWriter: w}

	// get the AKL session
	// TODO DI
	s, _ := router.Sessions.Get(req, "akl")

	// retrieve a new context
	var ctx = router.ContextFactory(w, req, s)

	// assign context to request
	req = req.WithContext(context.WithValue(req.Context(), web.CONTEXT, ctx))

	// dispatch OnRequest event, the request might be changed
	var event = &OnRequestEvent{w, req}
	ctx.EventRouter().Dispatch(REQUEST, event)
	req = event.Request

	defer func() {
		var err interface{}
		if err = recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintln(err)))
			w.Write(debug.Stack())
		}
		// fire finish event
		ctx.EventRouter().Dispatch(FINISH, &OnFinishEvent{w, req, err})
	}()

	if route, ok := router.hardroutes[req.URL.Path]; ok {
		p := make([]string, len(route.Args)*2)
		for k, v := range route.Args {
			p = append(p, k, v)
		}
		req.URL = router.url(route.Controller, p...)
	}

	router.router.ServeHTTP(w, req)
}

// handle sets the controller for a router which handles a Request.
func (router *Router) handle(c Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context().Value(web.CONTEXT).(web.Context) // get Request context
		ctx.LoadVars(req)                                     // LoadVars, because MuxVars has not resolved them in ServeHTTP

		defer ctx.Profile("request", req.RequestURI)()

		var response web.Response

		if cc, ok := c.(GETController); ok && req.Method == http.MethodGet {
			response = cc.Get(ctx)
		} else if cc, ok := c.(POSTController); ok && req.Method == http.MethodPost {
			response = cc.Post(ctx)
		} else {
			switch c := c.(type) {
			case func(web.Context) web.Response:
				response = c(ctx)

			case DataController:
				response = &web.JSONResponse{Data: c.(DataController).Data(ctx)}

			case func(web.Context) interface{}:
				response = &web.JSONResponse{Data: c(ctx)}

			case http.Handler:
				c.ServeHTTP(w, req)
				return

			default:
				w.WriteHeader(404)
				w.Write([]byte("404 page not found (no handler)"))
				return
			}
		}

		router.Sessions.Save(req, w, ctx.Session())

		// fire response event
		ctx.EventRouter().Dispatch(RESPONSE, &OnResponseEvent{c, response})

		response.Apply(w)
	})
}

// Get is the ServeHTTP's equivalent for DataController and DataHandler.
func (router *Router) Get(handler string, ctx web.Context) interface{} {
	defer ctx.Profile("get", handler)()

	if c, ok := router.handler[handler]; ok {
		if c, ok := c.(DataController); ok {
			return c.Data(ctx)
		}
		if c, ok := c.(func(web.Context) interface{}); ok {
			return c(ctx)
		}

		panic("not a data controller")
	} else { // mock...
		defer ctx.Profile("fallback", handler)
		data, err := ioutil.ReadFile("frontend/src/mocks/" + handler + ".json")
		if err == nil {
			var res interface{}
			json.Unmarshal(data, &res)
			return res
		}
		panic(err)
	}
}
