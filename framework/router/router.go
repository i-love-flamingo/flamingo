package router

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/profiler"
	"go.aoe.com/flamingo/framework/web"

	"fmt"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

const (
	// FlamingoError is the controller name for errors
	FlamingoError = "flamingo.error"
	// FlamingoNotfound is the controller name for 404 notfound
	FlamingoNotfound = "flamingo.notfound"
	// ERROR is used to bind errors to contexts
	ERROR ErrorKey = iota
)

type (
	ProfilerProvider    func() profiler.Profiler // ProfilerProvider for profiler injection
	EventRouterProvider func() event.Router      // EventRouterProvider for event injection

	// Router defines the basic Router which is used for holding a context-scoped setup
	// This includes DI resolving etc
	Router struct {
		base *url.URL

		Sessions            sessions.Store      `inject:",optional"` // Sessions storage, which are used to retrieve user-context session
		SessionName         string              `inject:"config:session.name"`
		ContextFactory      web.ContextFactory  `inject:""` // ContextFactory for new contexts
		ProfilerProvider    ProfilerProvider    `inject:""`
		EventRouterProvider EventRouterProvider `inject:""`
		eventrouter         event.Router
		Injector            *dingo.Injector `inject:""`
		RouterRegistry      *Registry       `inject:""`
		NotFoundHandler     string          `inject:"config:flamingo.router.notfound"`
		ErrorHandler        string          `inject:"config:flamingo.router.error"`
	}

	// P is a shorthand for parameter
	P map[string]string

	// ErrorKey for context errors
	ErrorKey uint
)

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return new(Router)
}

// SetBase for router
func (router *Router) SetBase(u *url.URL) {
	router.base = u
}

// Init the router
func (router *Router) Init(routingConfig *config.Area) *Router {
	if router.base == nil {
		router.base, _ = url.Parse("http://host")
	}

	// Make sure to not taint the global router registry
	var routes = NewRegistry()

	// build routes
	for _, route := range routingConfig.Routes {
		routes.Route(route.Path, route.Controller)
		if route.Name != "" {
			routes.Alias(route.Name, route.Controller)
		}
	}

	var routerroutes = make([]*Handler, len(router.RouterRegistry.routes))
	for k, v := range router.RouterRegistry.routes {
		routerroutes[k] = v
	}
	routes.routes = append(routes.routes, routerroutes...)

	// inject router instances
	for name, c := range router.RouterRegistry.handler {
		switch c.(type) {
		case http.Handler:
		case func(web.Context) web.Response:
		case func(web.Context) interface{}:
		case GETController, POSTController, HEADController, PUTController, DELETEController, DataController:
			c = router.Injector.GetInstance(reflect.TypeOf(c))
		default:
			var rv = reflect.ValueOf(c)
			if !rv.IsValid() {
				panic(fmt.Sprintf("Invalid Controller bound! %s: %#v", name, c))
			}
			// Check if we have a Receiver Function of the type
			// func(c Controller, ctx web.Context) web.Response
			// If so, we instantiate c Controller and convert it to
			// c.func(ctx web.Context) web.Response
			if rv.Type().Kind() == reflect.Func &&
				rv.Type().NumIn() == 2 &&
				rv.Type().NumOut() == 1 &&
				rv.Type().In(1).AssignableTo(reflect.TypeOf((*web.Context)(nil)).Elem()) &&
				rv.Type().Out(0).AssignableTo(reflect.TypeOf((*web.Response)(nil)).Elem()) {
				var ci = reflect.ValueOf(router.Injector.GetInstance(rv.Type().In(0).Elem()))
				c = func(ctx web.Context) web.Response {
					return rv.Call([]reflect.Value{ci, reflect.ValueOf(ctx)})[0].Interface().(web.Response)
				}
			}
		}
		routes.handler[name] = c
	}

	router.RouterRegistry = routes

	router.eventrouter = router.EventRouterProvider()

	return router
}

// URL helps resolving URL's by it's name.
func (router *Router) URL(name string, params map[string]string) *url.URL {
	var resultURL = new(url.URL)

	parts := strings.SplitN(name, "?", 2)
	name = parts[0]

	if len(parts) == 2 {
		var query, _ = url.ParseQuery(parts[1])
		resultURL.RawQuery = query.Encode()
	}

	p, err := router.RouterRegistry.Reverse(name, params)
	if err != nil {
		panic(err)
	}

	resultURL.Path = router.base.Path + p

	return resultURL
}

func (router *Router) recover(ctx web.Context, rw http.ResponseWriter, err interface{}) {
	defer func() {
		if err := recover(); err != nil {
			// bad bad recover
			rw.WriteHeader(http.StatusInternalServerError)
			if err, ok := err.(error); ok {
				fmt.Fprintf(rw, "%+v", errors.WithStack(err))
				return
			}
			fmt.Fprintf(rw, "%+v", errors.Errorf("%+v", err))
		}
	}()

	if e, ok := err.(error); ok {
		log.Printf("%+v", errors.WithStack(e))
		router.RouterRegistry.handler[router.ErrorHandler].(func(web.Context) web.Response)(ctx.WithValue(ERROR, errors.WithStack(e))).Apply(ctx, rw)
	} else if err, ok := err.(string); ok {
		router.RouterRegistry.handler[router.ErrorHandler].(func(web.Context) web.Response)(ctx.WithValue(ERROR, errors.New(err))).Apply(ctx, rw)
	} else {
		router.RouterRegistry.handler[router.ErrorHandler].(func(web.Context) web.Response)(ctx).Apply(ctx, rw)
	}
}

// ServeHTTP shadows the internal mux.Router's ServeHTTP to defer panic recoveries and logging.
// TODO simplify and merge with `handle`
func (router *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// shadow the response writer
	rw = &VerboseResponseWriter{ResponseWriter: rw}

	var s *sessions.Session
	var err error

	// initialize the session
	if router.Sessions != nil {
		s, err = router.Sessions.Get(req, router.SessionName)
		if err != nil {
			log.Println(err)
			s, err = router.Sessions.New(req, router.SessionName)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// retrieve a new context
	var ctx = router.ContextFactory(router.ProfilerProvider(), router.eventrouter, rw, req, s)

	// assign context to request
	req = req.WithContext(context.WithValue(req.Context(), web.CONTEXT, ctx))

	// dispatch OnRequest event, the request might be changed
	var e = &OnRequestEvent{rw, req, ctx}
	router.eventrouter.Dispatch(e)
	req = e.Request

	// catch errors
	defer func() {
		if err := recover(); err != nil {
			router.recover(ctx, rw, err)
		}
		// fire finish event
		router.eventrouter.Dispatch(&OnFinishEvent{rw, req, err, ctx})
	}()

	var controller, params, handler = router.RouterRegistry.MatchRequest(req)
	ctx.LoadParams(params)
	if controller == nil {
		controller = router.RouterRegistry.handler[router.NotFoundHandler]
	}
	ctx.WithValue("Handler", handlerdata{params, handler})
	router.handle(controller).ServeHTTP(rw, req)
}

// handle sets the controller for a router which handles a Request.
func (router *Router) handle(c Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context().Value(web.CONTEXT).(web.Context) // get Request context

		defer ctx.Profile("request", req.RequestURI)()

		var response web.Response

		if cc, ok := c.(GETController); ok && req.Method == http.MethodGet {
			response = cc.Get(ctx)
		} else if cc, ok := c.(POSTController); ok && req.Method == http.MethodPost {
			response = cc.Post(ctx)
		} else if cc, ok := c.(PUTController); ok && req.Method == http.MethodPut {
			response = cc.Put(ctx)
		} else if cc, ok := c.(DELETEController); ok && req.Method == http.MethodDelete {
			response = cc.Delete(ctx)
		} else if cc, ok := c.(HEADController); ok && req.Method == http.MethodHead {
			response = cc.Head(ctx)
		} else {
			switch c := c.(type) {
			case DataController:
				response = &web.JSONResponse{Data: c.Data(ctx)}

			case func(web.Context) web.Response:
				response = c(ctx)

			case func(web.Context) interface{}:
				response = &web.JSONResponse{Data: c(ctx)}

			case http.Handler:
				c.ServeHTTP(w, req)

			default:
				response = router.RouterRegistry.handler[router.ErrorHandler].(func(web.Context) web.Response)(ctx)
			}
		}

		if response, ok := response.(web.OnResponse); ok {
			response.OnResponse(ctx, w)
		}

		// fire response event
		router.eventrouter.Dispatch(&OnResponseEvent{c, response, req, w, ctx})

		if router.Sessions != nil {
			if err := router.Sessions.Save(req, w, ctx.Session()); err != nil {
				log.Println(err)
			}
		}

		if response != nil {
			response.Apply(ctx, w)
		}
	})
}

// Get is the ServeHTTP's equivalent for DataController and DataHandler.
// TODO refactor
func (router *Router) Get(handler string, ctx web.Context, params ...map[interface{}]interface{}) interface{} {
	defer ctx.Profile("get", handler)()

	// reformat data to map[string]string, just as in normal request vars would look like
	// dataController might be called via Ajax (instead of right via template) so this should be unified
	vars := reformatParams(ctx, params...)
	getCtx := ctx.WithVars(vars)

	if c, ok := router.RouterRegistry.handler[handler]; ok {
		if c, ok := c.(DataController); ok {
			return router.Injector.GetInstance(c).(DataController).Data(getCtx)
		}
		if c, ok := c.(func(web.Context) interface{}); ok {
			return c(getCtx)
		}
		panic(errors.Errorf("%q is not a data controller", handler))
	}
	panic(errors.Errorf("data controller %q not found", handler))
}

func reformatParams(ctx web.Context, params ...map[interface{}]interface{}) map[string]string {
	vars := make(map[string]string)
	for k, v := range ctx.ParamAll() {
		vars[k] = v
	}

	if len(params) == 1 {
		for k, v := range params[0] {
			if k, ok := k.(string); ok {
				switch v := v.(type) {
				case string:
					vars[k] = v
				case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
					vars[k] = strconv.Itoa(int(reflect.ValueOf(v).Int()))
				case float32:
					vars[k] = strconv.FormatFloat(float64(v), 'f', -1, 32)
				case float64:
					vars[k] = strconv.FormatFloat(v, 'f', -1, 64)
				}
			}
		}
	}
	return vars
}
