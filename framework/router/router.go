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
	"time"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/opencensus"
	"flamingo.me/flamingo/framework/session"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
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
		EventRouterProvider    EventRouterProvider `inject:""`
		eventrouter            event.Router
		Injector               *dingo.Injector  `inject:""`
		RouterRegistryProvider RegistryProvider `inject:""`
		RouterRegistry         *Registry        `inject:""`
		RouterTimeout          float64          `inject:"config:flamingo.router.timeout"`
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

var (
	rt = stats.Int64("flamingo/router/controller", "controller request times", stats.UnitMilliseconds)
	// ControllerKey exposes the current controller/handler key
	ControllerKey, _ = tag.NewKey("controller")
)

func init() {
	opencensus.View("flamingo/router/controller", rt, view.Distribution(100, 500, 1000, 2500, 5000, 10000), ControllerKey)
}

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

	for name, ha := range router.RouterRegistry.handler {
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

	// TODO: this is deprecated
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
func (router *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO simplify and merge with `handle`
	serveCtx, span := trace.StartSpan(req.Context(), "router/ServeHTTP")

	// shadow the response writer
	rw = &web.VerboseResponseWriter{ResponseWriter: rw}

	deadlineContext, cancelFunc := context.WithTimeout(
		context.WithValue(req.Context(), "rw", rw),
		time.Duration(router.RouterTimeout)*time.Millisecond,
	)
	req = req.WithContext(deadlineContext)

	defer cancelFunc()

	var s *sessions.Session
	var err error

	// initialize the session
	if router.Sessions != nil {
		ctx, span := trace.StartSpan(serveCtx, "router/sessions/get")
		s, err = router.Sessions.Get(req, router.SessionName)
		if err != nil {
			log.Println(err)
			_, span := trace.StartSpan(ctx, "router/sessions/new")
			s, err = router.Sessions.New(req, router.SessionName)
			if err != nil {
				log.Println(err)
			}
			span.End()
		}
		span.End()
	}

	req = req.WithContext(session.Context(req.Context(), s))

	// retrieve a new context
	ctx := req.Context()

	// dispatch OnRequest event, the request might be changed
	e := &OnRequestEvent{rw, req}
	router.eventrouter.Dispatch(ctx, e)
	req = e.Request

	span.End()

	_, span = trace.StartSpan(req.Context(), "router/matchRequest")
	controller, params, handler := router.RouterRegistry.matchRequest(req)

	if handler != nil {
		ctx, _ = tag.New(req.Context(), tag.Upsert(ControllerKey, handler.GetHandlerName()))
		req = req.WithContext(ctx)
		start := time.Now()
		defer func() {
			stats.Record(req.Context(), rt.M(time.Since(start).Nanoseconds()/1000000))
		}()
	}

	span.End()

	ctx, span = trace.StartSpan(req.Context(), "router/request")
	req = req.WithContext(ctx)
	defer span.End()

	webRequest := web.RequestFromRequest(req, s).WithVars(params)

	chain := &FilterChain{
		Filters: make([]Filter, len(router.filters)),
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
			router.eventrouter.Dispatch(ctx, &OnFinishEvent{rw, req, err})
		}()

		var response web.Response

		if c, ok := controller.method[req.Method]; ok && c != nil {
			response = c(ctx, webRequest)
		} else if controller.any != nil {
			response = controller.any(ctx, webRequest)
		} else {
			response = router.RouterRegistry.handler[router.NotFoundHandler].any(context.WithValue(ctx, ERROR, errors.Errorf("action for method %q not found and no any fallback", req.Method)), r)
		}

		if response, ok := response.(web.OnResponse); ok {
			response.OnResponse(ctx, webRequest, rw)
		}

		// fire response event
		router.eventrouter.Dispatch(ctx, &OnResponseEvent{response, req, rw})

		return response
	}))

	response := chain.Next(web.Context_(ctx, webRequest), webRequest, rw)

	if router.Sessions != nil {
		_, span := trace.StartSpan(ctx, "router/sessions/safe")
		if err := router.Sessions.Save(req, rw, webRequest.Session()); err != nil {
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

func dataParams(r *web.Request, params map[interface{}]interface{}) map[string]string {
	vars := make(map[string]string)
	for k, v := range r.ParamAll() {
		vars[k] = v
	}

	for k, v := range params {
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

	return vars
}

// Data calls a flamingo data controller
func (router *Router) Data(ctx context.Context, handler string, params map[interface{}]interface{}) interface{} {
	ctx, span := trace.StartSpan(ctx, "flamingo/router/data")
	span.Annotate(nil, handler)
	defer span.End()

	r, ok := web.FromContext(ctx)
	if !ok {
		r = web.RequestFromRequest(nil, sessions.NewSession(router.Sessions, "-"))
	}

	r.LoadParams(dataParams(r, params))

	if c, ok := router.RouterRegistry.handler[handler]; ok {
		if c.data != nil {
			return c.data(ctx, r)
		}
		panic(errors.Errorf("%q is not a data Controller", handler))
	}
	panic(errors.Errorf("data Controller %q not found", handler))
}
