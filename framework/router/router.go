package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/profiler"
	"go.aoe.com/flamingo/framework/web"
)

const (
	// FlamingoError is the Controller name for errors
	FlamingoError = "flamingo.error"
	// FlamingoNotfound is the Controller name for 404 notfound
	FlamingoNotfound = "flamingo.notfound"

	// ERROR is used to bind errors to contexts
	ERROR errorKey = iota
)

type (
	// ProfilerProvider for profiler injection
	ProfilerProvider func() profiler.Profiler
	// EventRouterProvider for event injection
	EventRouterProvider func() event.Router
	// FilterProvider for filter injection
	FilterProvider func() []Filter

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
		FilterProvider      FilterProvider  `inject:",optional"`
		filters             []Filter
	}

	// P is a shorthand for parameter
	P map[string]string

	// errorKey for context errors
	errorKey uint
)

func NewRouter() *Router {
	return new(Router)
}

// Init the router
func (router *Router) Init(routingConfig *config.Area) *Router {
	if router.base == nil {
		router.base, _ = url.Parse("http://host")
	}

	// Make sure to not taint the global router registry
	routes := NewRegistry()

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
		case http.Handler, func(web.Context) web.Response, func(web.Context) interface{}:
			break

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

	for _, handler := range routes.routes {
		if _, ok := routes.handler[handler.handler]; !ok {
			panic(errors.Errorf("The handler %q has no controller, registered for path %q", handler.handler, handler.path.path))
		}
	}
	router.RouterRegistry = routes

	router.eventrouter = router.EventRouterProvider()
	router.filters = router.FilterProvider()

	return router
}

func (router *Router) Base() *url.URL {
	return router.base
}

// SetBase for router
func (router *Router) SetBase(u *url.URL) {
	router.base = u
}

// URL helps resolving URL's by it's name.
func (router *Router) URL(name string, params map[string]string) *url.URL {
	var resultURL = new(url.URL)

	// todo: this is deprecated
	parts := strings.SplitN(name, "?", 2)
	name = parts[0]

	if len(parts) == 2 {
		log.Println("the usage of `?` in url(...) is deprecated")
		var query, _ = url.ParseQuery(parts[1])
		resultURL.RawQuery = query.Encode()
	}

	p, err := router.RouterRegistry.Reverse(name, params)
	if err != nil {
		panic(err)
	}
	resultURL, err = url.Parse(router.base.Path + p)
	if err != nil {
		panic(err)
	}

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
	rw = &web.VerboseResponseWriter{ResponseWriter: rw}

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
	ctx := router.ContextFactory(router.ProfilerProvider(), router.eventrouter, rw, req, s)

	// assign context to request
	req = req.WithContext(context.WithValue(req.Context(), web.CONTEXT, ctx))

	// dispatch OnRequest event, the request might be changed
	e := &OnRequestEvent{rw, req, ctx}
	router.eventrouter.Dispatch(e)
	req = e.Request

	done := ctx.Profile("matchRequest", req.RequestURI)
	controller, params, handler := router.RouterRegistry.MatchRequest(req)
	ctx.LoadParams(params)
	if controller == nil {
		controller = router.RouterRegistry.handler[router.NotFoundHandler]
	}
	if handler != nil {
		ctx.WithValue("HandlerName", handler.GetHandlerName())
	}
	ctx.WithValue("Handler", handlerdata{params, handler})
	done()

	defer ctx.Profile("request", req.RequestURI)()

	chain := &FilterChain{
		Filters:    make([]Filter, len(router.filters)),
		Controller: controller,
	}
	copy(chain.Filters, router.filters)

	chain.Filters = append(chain.Filters, lastFilter(func(ctx web.Context, rw http.ResponseWriter) web.Response {
		// catch errors
		defer func() {
			if err := recover(); err != nil {
				router.recover(ctx, rw, err)
			}
			// fire finish event
			router.eventrouter.Dispatch(&OnFinishEvent{rw, req, err, ctx})
		}()

		var response web.Response

		if cc, ok := controller.(GETController); ok && req.Method == http.MethodGet {
			response = cc.Get(ctx)
		} else if cc, ok := controller.(POSTController); ok && req.Method == http.MethodPost {
			response = cc.Post(ctx)
		} else if cc, ok := controller.(PUTController); ok && req.Method == http.MethodPut {
			response = cc.Put(ctx)
		} else if cc, ok := controller.(DELETEController); ok && req.Method == http.MethodDelete {
			response = cc.Delete(ctx)
		} else if cc, ok := controller.(HEADController); ok && req.Method == http.MethodHead {
			response = cc.Head(ctx)
		} else {
			switch c := controller.(type) {
			case DataController:
				response = &web.JSONResponse{Data: c.Data(ctx)}

			case func(web.Context) web.Response:
				response = c(ctx)

			case func(web.Context) interface{}:
				response = &web.JSONResponse{Data: c(ctx)}

			case http.Handler:
				response = &web.ServeHTTPResponse{VerboseResponseWriter: rw.(*web.VerboseResponseWriter)}
				c.ServeHTTP(response.(http.ResponseWriter), req)

			default:
				response = router.RouterRegistry.handler[router.ErrorHandler].(func(web.Context) web.Response)(ctx)
			}
		}

		if response, ok := response.(web.OnResponse); ok {
			response.OnResponse(ctx, rw)
		}

		// fire response event
		router.eventrouter.Dispatch(&OnResponseEvent{controller, response, req, rw, ctx})

		if router.Sessions != nil {
			if err := router.Sessions.Save(req, rw, ctx.Session()); err != nil {
				log.Println(err)
			}
		}

		return response
	}))

	response := chain.Next(ctx, rw)

	if response != nil {
		response.Apply(ctx, rw)
	}
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
		panic(errors.Errorf("%q is not a data Controller", handler))
	}
	panic(errors.Errorf("data Controller %q not found", handler))
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
