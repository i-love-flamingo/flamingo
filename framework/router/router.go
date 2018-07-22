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

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/profiler"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
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
	// RegistryProvider is called to retrieve registered routes
	RegistryProvider func() []Module

	// Router defines the basic Router which is used for holding a context-scoped setup
	// This includes DI resolving etc
	Router struct {
		base *url.URL

		Sessions               sessions.Store      `inject:",optional"` // Sessions storage, which are used to retrieve user-context session
		SessionName            string              `inject:"config:session.name"`
		ContextFactory         web.ContextFactory  `inject:""` // ContextFactory for new contexts
		ProfilerProvider       ProfilerProvider    `inject:""`
		EventRouterProvider    EventRouterProvider `inject:""`
		eventrouter            event.Router
		Injector               *dingo.Injector  `inject:""`
		RouterRegistryProvider RegistryProvider `inject:""`
		RouterRegistry         *Registry        `inject:""`
		NotFoundHandler        string           `inject:"config:flamingo.router.notfound"`
		ErrorHandler           string           `inject:"config:flamingo.router.error"`
		FilterProvider         FilterProvider   `inject:",optional"`
		filters                []Filter
	}

	// P is a shorthand for parameter
	P map[string]string

	// errorKey for context errors
	errorKey uint
)

// NewRouter factory
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

	if router.RouterRegistryProvider != nil {
		for _, m := range router.RouterRegistryProvider() {
			m.Routes(router.RouterRegistry)
		}
	}

	var routerroutes = make([]*Handler, len(router.RouterRegistry.routes))
	for k, v := range router.RouterRegistry.routes {
		routerroutes[k] = v
	}
	routes.routes = append(routes.routes, routerroutes...)

	// inject router instances
	// deprecated: only used for legacy controllers
	for name, ha := range router.RouterRegistry.handler {
		c := ha.legacyController
		if c != nil {
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
		}
		ha.legacyController = c
		routes.handler[name] = ha
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

// Base URL getter
func (router *Router) Base() *url.URL {
	return router.base
}

// SetBase for router
func (router *Router) SetBase(u *url.URL) {
	router.base = u
}

// TryURL is the same as URL below, but checks if the url is possible and returns an error
func (router *Router) TryURL(name string, params map[string]string) (u *url.URL, err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Println(p)
			if perr, ok := p.(error); ok {
				err = perr
			} else {
				err = fmt.Errorf("%v", p)
			}
		}
	}()

	u = router.URL(name, params)
	return
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
	resultURL, err = url.Parse(strings.TrimRight(router.base.Path, "/") + "/" + strings.TrimLeft(p, "/"))
	if err != nil {
		panic(err)
	}

	return resultURL
}

func (router *Router) recover(ctx context.Context, r *web.Request, rw http.ResponseWriter, err interface{}) {
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
		router.RouterRegistry.handler[router.ErrorHandler].any(context.WithValue(ctx, ERROR, errors.WithStack(e)), r).Apply(ctx, rw)
	} else if err, ok := err.(string); ok {
		router.RouterRegistry.handler[router.ErrorHandler].any(context.WithValue(ctx, ERROR, errors.New(err)), r).Apply(ctx, rw)
	} else {
		router.RouterRegistry.handler[router.ErrorHandler].any(ctx, r).Apply(ctx, rw)
	}
}

// ServeHTTP shadows the internal mux.Router's ServeHTTP to defer panic recoveries and logging.
// TODO simplify and merge with `handle`
func (router *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, span := trace.StartSpan(req.Context(), "router/ServeHTTP")

	// shadow the response writer
	rw = &web.VerboseResponseWriter{ResponseWriter: rw}
	req = req.WithContext(context.WithValue(req.Context(), "rw", rw))

	var s *sessions.Session
	var err error

	// initialize the session
	if router.Sessions != nil {
		_, span := trace.StartSpan(req.Context(), "router/sessions/get")
		s, err = router.Sessions.Get(req, router.SessionName)
		if err != nil {
			log.Println(err)
			_, span := trace.StartSpan(req.Context(), "router/sessions/new")
			s, err = router.Sessions.New(req, router.SessionName)
			if err != nil {
				log.Println(err)
			}
			span.End()
		}
		span.End()
	}

	// retrieve a new context
	ctx := router.ContextFactory(router.ProfilerProvider(), router.eventrouter, rw, req, s)

	// assign context to request
	req = req.WithContext(context.WithValue(req.Context(), web.CONTEXT, ctx))

	// dispatch OnRequest event, the request might be changed
	e := &OnRequestEvent{rw, req, ctx}
	router.eventrouter.Dispatch(ctx, e)
	req = e.Request

	span.End()

	done := ctx.Profile("matchRequest", req.RequestURI)
	_, span = trace.StartSpan(req.Context(), "router/matchRequest")
	controller, params, handler := router.RouterRegistry.matchRequest(req)

	ctx.LoadParams(params)
	if handler != nil {
		ctx.WithValue("HandlerName", handler.GetHandlerName())
	}
	ctx.WithValue("Handler", handlerdata{params, handler})
	done()
	span.End()

	defer ctx.Profile("request", req.RequestURI)()

	tracectx, span := trace.StartSpan(req.Context(), "router/request")
	req = req.WithContext(tracectx)
	defer span.End()

	webRequest := web.RequestFromRequest(req, s).WithVars(params)
	ctx.WithValue("__req", webRequest)

	chain := &FilterChain{
		Filters:    make([]Filter, len(router.filters)),
		Controller: controller,
	}
	copy(chain.Filters, router.filters)

	chain.Filters = append(chain.Filters, lastFilter(func(ctx context.Context, r *web.Request, rw http.ResponseWriter) web.Response {
		ctx, span := trace.StartSpan(ctx, "router/controller")
		defer span.End()

		// catch errors
		defer func() {
			if err := recover(); err != nil {
				router.recover(ctx, r, rw, err)
			}
			// fire finish event
			router.eventrouter.Dispatch(ctx, &OnFinishEvent{rw, req, err, web.ToContext(ctx)})
		}()

		var response web.Response

		if c, ok := controller.method[req.Method]; ok {
			response = c(ctx, webRequest)
		} else if controller.any != nil {
			response = controller.any(ctx, webRequest)
		} else {
			// deprecated: refactored in favor of proper controller actions
			if cc, ok := controller.legacyController.(GETController); ok && req.Method == http.MethodGet {
				response = cc.Get(web.ToContext(ctx))
			} else if cc, ok := controller.legacyController.(POSTController); ok && req.Method == http.MethodPost {
				response = cc.Post(web.ToContext(ctx))
			} else if cc, ok := controller.legacyController.(PUTController); ok && req.Method == http.MethodPut {
				response = cc.Put(web.ToContext(ctx))
			} else if cc, ok := controller.legacyController.(DELETEController); ok && req.Method == http.MethodDelete {
				response = cc.Delete(web.ToContext(ctx))
			} else if cc, ok := controller.legacyController.(HEADController); ok && req.Method == http.MethodHead {
				response = cc.Head(web.ToContext(ctx))
			} else {
				switch c := controller.legacyController.(type) {
				case DataController:
					response = &web.JSONResponse{Data: c.Data(web.ToContext(ctx))}

				case func(web.Context) web.Response:
					response = c(web.ToContext(ctx))

				case func(web.Context) interface{}:
					response = &web.JSONResponse{Data: c(web.ToContext(ctx))}

				case http.Handler:
					response = &web.ServeHTTPResponse{VerboseResponseWriter: rw.(*web.VerboseResponseWriter)}
					c.ServeHTTP(response.(http.ResponseWriter), req)

				default:
					response = router.RouterRegistry.handler[router.NotFoundHandler].any(context.WithValue(ctx, ERROR, errors.Errorf("legacy controller type unknown/unset: %T", c)), r)
				}
			}
		}

		if response, ok := response.(web.OnResponse); ok {
			response.OnResponse(ctx, webRequest, rw)
		}

		// fire response event
		router.eventrouter.Dispatch(ctx, &OnResponseEvent{controller, response, req, rw, web.ToContext(ctx)})

		return response
	}))

	response := chain.Next(ctx, webRequest, rw)

	if router.Sessions != nil {
		_, span := trace.StartSpan(ctx, "router/sessions/safe")
		if err := router.Sessions.Save(req, rw, ctx.Session()); err != nil {
			log.Println(err)
		}
		span.End()
	}

	if response != nil {
		_, span := trace.StartSpan(ctx, "router/responseApply")
		response.Apply(ctx, rw)
		span.End()
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
		if c, ok := c.legacyController.(DataController); ok {
			return router.Injector.GetInstance(c).(DataController).Data(getCtx)
		}
		if c, ok := c.legacyController.(func(web.Context) interface{}); ok {
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
