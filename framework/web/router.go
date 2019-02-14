package web

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
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
	// deprecated fix
	ERROR errorKey = iota
)

type (
	// EventRouterProvider for event injection
	EventRouterProvider func() flamingo.EventRouter
	// FilterProvider for filter injection
	FilterProvider func() []Filter
	// RegistryProvider is called to retrieve registered routes
	RegistryProvider func() []RoutesModule

	// Router defines the basic Router which is used for holding a context-scoped setup
	// This includes DI resolving etc
	Router struct {
		base *url.URL

		sessionStore           sessions.Store
		sessionName            string
		eventRouterProvider    EventRouterProvider
		eventrouter            flamingo.EventRouter
		injector               *dingo.Injector
		routerRegistryProvider RegistryProvider
		routerRegistry         *RouterRegistry
		routerTimeout          float64
		notFoundHandler        string
		errorHandler           string
		filterProvider         FilterProvider
		filters                []Filter
	}

	// errorKey for context errors
	errorKey uint

	emptyResponseWriter struct{}
)

var (
	rt = stats.Int64("flamingo/router/controller", "controller request times", stats.UnitMilliseconds)
	// ControllerKey exposes the current controller/handler key
	ControllerKey, _ = tag.NewKey("controller")
)

func init() {
	if err := opencensus.View("flamingo/router/controller", rt, view.Distribution(100, 500, 1000, 2500, 5000, 10000), ControllerKey); err != nil {
		panic(err)
	}
}

// Inject dependencies
func (router *Router) Inject(
	eventRouterProvider EventRouterProvider,
	injector *dingo.Injector,
	routerRegistryProvider RegistryProvider,
	routerRegistry *RouterRegistry,
	cfg *struct {
		SessionName     string         `inject:"config:session.name"`
		SessionStore    sessions.Store `inject:",optional"`
		RouterTimeout   float64        `inject:"config:flamingo.router.timeout"`
		NotFoundHandler string         `inject:"config:flamingo.router.notfound"`
		ErrorHandler    string         `inject:"config:flamingo.router.error"`
		FilterProvider  FilterProvider `inject:",optional"`
	},
) *Router {
	router.eventRouterProvider = eventRouterProvider
	router.injector = injector
	router.routerRegistryProvider = routerRegistryProvider
	router.routerRegistry = routerRegistry
	router.sessionName = cfg.SessionName
	router.sessionStore = cfg.SessionStore
	router.routerTimeout = cfg.RouterTimeout
	router.notFoundHandler = cfg.NotFoundHandler
	router.errorHandler = cfg.ErrorHandler
	router.filterProvider = cfg.FilterProvider
	return router
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

	if router.routerRegistryProvider != nil {
		for _, m := range router.routerRegistryProvider() {
			m.Routes(router.routerRegistry)
		}
	}

	var routerroutes = make([]*Handler, len(router.routerRegistry.routes))
	for k, v := range router.routerRegistry.routes {
		routerroutes[k] = v
	}
	routes.routes = append(routes.routes, routerroutes...)

	for name, ha := range router.routerRegistry.handler {
		routes.handler[name] = ha
	}

	for _, handler := range routes.routes {
		if _, ok := routes.handler[handler.handler]; !ok {
			panic(errors.Errorf("The handler %q has no controller, registered for path %q", handler.handler, handler.path.path))
		}
	}
	router.routerRegistry = routes

	router.eventrouter = router.eventRouterProvider()
	router.filters = router.filterProvider()

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

// URL helps resolving URL's by it's name.
func (router *Router) URL(name string, params map[string]string) (*url.URL, error) {
	p, err := router.routerRegistry.Reverse(name, params)
	if err != nil {
		return nil, err
	}
	return url.Parse(strings.TrimRight(router.base.Path, "/") + "/" + strings.TrimLeft(p, "/"))
}

func (router *Router) recover(ctx context.Context, r *Request, rw http.ResponseWriter, err interface{}) {
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
		router.routerRegistry.handler[router.errorHandler].any(context.WithValue(ctx, ERROR, errors.WithStack(e)), r).Apply(ctx, rw)
	} else if err, ok := err.(string); ok {
		router.routerRegistry.handler[router.errorHandler].any(context.WithValue(ctx, ERROR, errors.New(err)), r).Apply(ctx, rw)
	} else {
		router.routerRegistry.handler[router.errorHandler].any(ctx, r).Apply(ctx, rw)
	}
}

// ServeHTTP shadows the internal mux.Router's ServeHTTP to defer panic recoveries and logging.
func (router *Router) ServeHTTP(rw http.ResponseWriter, httpRequest *http.Request) {
	// TODO simplify and merge with `handle`
	ctx, span := trace.StartSpan(httpRequest.Context(), "router/ServeHTTP")

	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithTimeout(
		ctx,
		time.Duration(router.routerTimeout)*time.Millisecond, // todo how about changing this?
	)
	defer cancelFunc()

	var gs *sessions.Session
	var err error

	// initialize the session
	if router.sessionStore != nil {
		var span *trace.Span
		ctx, span = trace.StartSpan(ctx, "router/sessions/get")
		gs, err = router.sessionStore.Get(httpRequest, router.sessionName)
		if err != nil {
			log.Println(err)
			_, span := trace.StartSpan(ctx, "router/sessions/new")
			gs, err = router.sessionStore.New(httpRequest, router.sessionName)
			if err != nil {
				log.Println(err)
			}
			span.End()
		}
		span.End()
	}

	// dispatch OnRequest event, the request might be changed
	//req = e.Request

	span.End()

	ctx, span = trace.StartSpan(ctx, "router/matchRequest")
	controller, params, handler := router.routerRegistry.matchRequest(httpRequest)

	if handler != nil {
		ctx, _ = tag.New(ctx, tag.Upsert(ControllerKey, handler.GetHandlerName()), tag.Upsert(opencensus.KeyArea, "-"))
		httpRequest = httpRequest.WithContext(ctx)
		start := time.Now()
		defer func() {
			stats.Record(ctx, rt.M(time.Since(start).Nanoseconds()/1000000))
		}()
	}

	req := &Request{
		request: *httpRequest,
		session: Session{
			s: gs,
		},
		Params: params,
	}
	ctx = ContextWithRequest(ContextWithSession(ctx, req.Session()), req)

	e := &OnRequestEvent{req, rw}
	router.eventrouter.Dispatch(ctx, e)

	span.End()

	ctx, span = trace.StartSpan(ctx, "router/request")
	defer span.End()

	chain := &FilterChain{
		Filters: make([]Filter, 0, len(router.filters)+1),
	}
	copy(chain.Filters, router.filters)

	chain.Filters = append(chain.Filters, lastFilter(func(ctx context.Context, r *Request, rw http.ResponseWriter) Result {
		ctx, span := trace.StartSpan(ctx, "router/controller")
		defer span.End()

		// catch errors
		defer func() {
			if err := recover(); err != nil {
				router.recover(ctx, r, rw, err)
			}
			// fire finish event
			router.eventrouter.Dispatch(ctx, &OnFinishEvent{OnRequestEvent{req, rw}, err})
		}()

		var response Result

		if c, ok := controller.method[req.Request().Method]; ok && c != nil {
			response = c(ctx, r)
		} else if controller.any != nil {
			response = controller.any(ctx, r)
		} else {
			response = router.routerRegistry.handler[router.notFoundHandler].any(context.WithValue(ctx, ERROR, errors.Errorf("action for method %q not found and no any fallback", req.Request().Method)), r)
		}

		if response, ok := response.(onResponse); ok {
			response.OnResponse(ctx, r, rw)
		}

		// fire response event
		router.eventrouter.Dispatch(ctx, &OnResponseEvent{OnRequestEvent{req, rw}, response})

		return response
	}))

	response := chain.Next(ctx, req, rw)

	if router.sessionStore != nil {
		_, span := trace.StartSpan(ctx, "router/sessions/save")
		if err := router.sessionStore.Save(req.Request(), rw, gs); err != nil {
			log.Println(err)
		}
		span.End()
	}

	if response != nil {
		_, span := trace.StartSpan(ctx, "router/responseApply")
		if err := response.Apply(ctx, rw); err != nil {
			panic(err) // bail out?
		}
		span.End()
	}

	// ensure that the session has been saved in the backend
	if router.sessionStore != nil {
		_, span := trace.StartSpan(ctx, "router/sessions/persist")
		if err := router.sessionStore.Save(req.Request(), emptyResponseWriter{}, gs); err != nil {
			log.Println(err)
		}
		span.End()
	}
}

func dataParams(rParams RequestParams, params map[interface{}]interface{}) RequestParams {
	vars := make(map[string]string)
	for k, v := range rParams {
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

	r := RequestFromContext(ctx)

	if c, ok := router.routerRegistry.handler[handler]; ok {
		if c.data != nil {
			return c.data(ctx, r, dataParams(r.Params, params))
		}
		panic(errors.Errorf("%q is not a data Controller", handler))
	}
	panic(errors.Errorf("data Controller %q not found", handler))
}

func (emptyResponseWriter) Header() http.Header {
	return http.Header{}
}

func (emptyResponseWriter) Write([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (emptyResponseWriter) WriteHeader(statusCode int) {}
