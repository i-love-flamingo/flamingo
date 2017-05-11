package router

import (
	"context"
	"encoding/json"
	configcontext "flamingo/framework/context"
	"flamingo/framework/dingo"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/web"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
	"strings"

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
		router            *mux.Router
		hardroutes        map[string]configcontext.Route
		hardroutesreverse map[string]configcontext.Route
		base              *url.URL

		Sessions            sessions.Store           `inject:""` // Sessions storage, which are used to retrieve user-context session
		SessionName         string                   `inject:"config:session.name"`
		ContextFactory      web.ContextFactory       `inject:""` // ContextFactory for new contexts
		ProfilerProvider    func() profiler.Profiler `inject:""`
		EventRouterProvider func() event.Router      `inject:""`
		Injector            *dingo.Injector          `inject:""`
		RouterRegistry      *RouterRegistry          `inject:""`
	}
)

func NewRouter() *Router {
	router := &Router{
		router:            mux.NewRouter(),
		hardroutes:        make(map[string]configcontext.Route),
		hardroutesreverse: make(map[string]configcontext.Route),
	}

	return router
}

func (router *Router) Init(routingConfig *configcontext.RoutingConfig) *Router {
	router.base, _ = url.Parse("scheme://" + routingConfig.BaseURL)

	for _, route := range routingConfig.Routes {
		if route.Args == nil {
			router.RouterRegistry.routes[route.Controller] = route.Path
		} else {
			router.hardroutes[route.Path] = route
			p := make([]string, len(route.Args)*2)
			for k, v := range route.Args {
				p = append(p, k, v)
			}
			router.hardroutesreverse[route.Controller+strings.Join(p, "!!")] = route
		}
	}

	for name, handler := range router.RouterRegistry.handler {
		if route, ok := router.RouterRegistry.routes[name]; ok {
			router.router.Handle(route, router.handle(handler)).Name(name)
		}
	}

	return router
}

// URL helps resolving URL's by it's name.
// Example:
//     flamingo.URL("cms.page.view", "name", "Home")
// results in
//     /baseurl/cms/Home
func (router *Router) URL(name string, params ...string) *url.URL {
	var resultURL *url.URL

	query := url.Values{}
	parts := strings.SplitN(name, "?", 2)
	name = parts[0]
	if len(parts) == 2 {
		query, _ = url.ParseQuery(parts[1])
	}

	if route, ok := router.hardroutesreverse[name+`!!!!`+strings.Join(params, "!!")]; ok {
		resultURL, _ = url.Parse(route.Path)
	} else {
		resultURL = router.url(name, params...)
	}

	resultURL.Path = path.Join(router.base.Path, resultURL.Path)
	resultURL.RawQuery = query.Encode()

	return resultURL
}

// url builds a URL for a Router.
func (router *Router) url(name string, params ...string) *url.URL {
	if router.router.Get(name) == nil {
		//panic("route " + name + " not found")
		return &url.URL{
			Fragment: name + "::" + strings.Join(params, ":"),
		}
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

	// get the session
	s, _ := router.Sessions.Get(req, router.SessionName)

	// retrieve a new context
	var ctx = router.ContextFactory(router.ProfilerProvider(), router.EventRouterProvider(), w, req, s)

	// assign context to request
	req = req.WithContext(context.WithValue(req.Context(), web.CONTEXT, ctx))

	// dispatch OnRequest event, the request might be changed
	var e = &OnRequestEvent{w, req}
	ctx.EventRouter().Dispatch(e)
	req = e.Request

	defer func() {
		var err interface{}
		if err = recover(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintln(err)))
			w.Write(debug.Stack())
		}
		// fire finish event
		ctx.EventRouter().Dispatch(&OnFinishEvent{w, req, err})
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
			response = router.Injector.GetInstance(cc).(GETController).Get(ctx)
		} else if cc, ok := c.(POSTController); ok && req.Method == http.MethodPost {
			response = router.Injector.GetInstance(cc).(POSTController).Post(ctx)
		} else {
			switch c := c.(type) {
			case func(web.Context) web.Response:
				response = c(ctx)

			case DataController:
				response = &web.JSONResponse{Data: router.Injector.GetInstance(c).(DataController).Data(ctx)}

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

		// fire response event
		ctx.EventRouter().Dispatch(&OnResponseEvent{c, response, req, w})

		router.Sessions.Save(req, w, ctx.Session())

		response.Apply(ctx, w)
	})
}

// Get is the ServeHTTP's equivalent for DataController and DataHandler.
func (router *Router) Get(handler string, ctx web.Context) interface{} {
	defer ctx.Profile("get", handler)()

	if c, ok := router.RouterRegistry.handler[handler]; ok {
		if c, ok := c.(DataController); ok {
			return router.Injector.GetInstance(c).(DataController).Data(ctx)
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

func (router *Router) GetHardRoutes() map[string]configcontext.Route {
	return router.hardroutes
}
