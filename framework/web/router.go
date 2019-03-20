package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/config"

	"flamingo.me/dingo"
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
	eventRouterProvider func() flamingo.EventRouter
	filterProvider      func() []Filter
	registryProvider    func() []RoutesModule

	// Router defines the basic Router which is used for holding a context-scoped setup
	// A request is handled as follows:
	// the filter chain is called
	// -> within the filter chain, as the last action, the controller is called
	// possible errors:
	// - error result: normal handling as a response/result
	// - panic:
	// - apply error:
	// - apply panic:
	Router struct {
		base                   *url.URL
		sessionStore           sessions.Store
		sessionName            string
		eventRouterProvider    eventRouterProvider
		eventrouter            flamingo.EventRouter
		injector               *dingo.Injector
		routerRegistryProvider registryProvider
		routerRegistry         *RouterRegistry
		routerTimeout          float64
		notFoundHandler        string
		errorHandler           string
		filterProvider         filterProvider
		filters                []Filter
		logger                 flamingo.Logger
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
	eventRouterProvider eventRouterProvider,
	injector *dingo.Injector,
	routerRegistryProvider registryProvider,
	routerRegistry *RouterRegistry,
	logger flamingo.Logger,
	cfg *struct {
		SessionName     string         `inject:"config:session.name"`
		SessionStore    sessions.Store `inject:",optional"`
		RouterTimeout   float64        `inject:"config:flamingo.router.timeout"`
		NotFoundHandler string         `inject:"config:flamingo.router.notfound"`
		ErrorHandler    string         `inject:"config:flamingo.router.error"`
		FilterProvider  filterProvider `inject:",optional"`
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
	router.logger = logger
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

func (router *Router) getSession(ctx context.Context, httpRequest *http.Request) (gs *sessions.Session) {
	// initialize the session
	if router.sessionStore != nil {
		var span *trace.Span
		var err error

		ctx, span = trace.StartSpan(ctx, "router/sessions/get")
		gs, err = router.sessionStore.Get(httpRequest, router.sessionName)
		if err != nil {
			router.logger.WithContext(ctx).Warn(err)
			_, span := trace.StartSpan(ctx, "router/sessions/new")
			gs, err = router.sessionStore.New(httpRequest, router.sessionName)
			if err != nil {
				router.logger.WithContext(ctx).Warn(err)
			}
			span.End()
		}
		span.End()
	}

	return
}

func panicToError(p interface{}) error {
	if p == nil {
		return nil
	}

	var err error
	switch errIface := p.(type) {
	case error:
		err = errors.WithStack(errIface)
	case string:
		err = errors.New(errIface)
	default:
		err = errors.Errorf("router/controller: %+v", errIface)
	}
	return err
}

// ServeHTTP shadows the internal mux.Router's ServeHTTP to defer panic recoveries and logging.
func (router *Router) ServeHTTP(rw http.ResponseWriter, httpRequest *http.Request) {
	var err error

	ctx, span := trace.StartSpan(httpRequest.Context(), "router/ServeHTTP")
	defer span.End()

	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithTimeout(
		ctx,
		time.Duration(router.routerTimeout)*time.Millisecond, // todo how about changing this?
	)
	defer cancelFunc()

	gs := router.getSession(ctx, httpRequest)

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

	defer func() {
		// fire finish event
		router.eventrouter.Dispatch(ctx, &OnFinishEvent{OnRequestEvent{req, rw}, err})
	}()

	e := &OnRequestEvent{req, rw}
	router.eventrouter.Dispatch(ctx, e)

	span.End() // router/matchRequest

	ctx, span = trace.StartSpan(ctx, "router/request")
	defer span.End()

	chain := &FilterChain{
		filters: router.filters,
		final: func(ctx context.Context, r *Request, rw http.ResponseWriter) (response Result) {
			ctx, span := trace.StartSpan(ctx, "router/controller")
			defer span.End()

			defer func() {
				if err := panicToError(recover()); err != nil {
					response = router.routerRegistry.handler[router.errorHandler].any(context.WithValue(ctx, ERROR, err), r)
					span.SetStatus(trace.Status{Code: trace.StatusCodeAborted, Message: "controller panic"})
				}
			}()

			defer router.eventrouter.Dispatch(ctx, &OnResponseEvent{OnRequestEvent{req, rw}, response})

			if c, ok := controller.method[req.Request().Method]; ok && c != nil {
				response = c(ctx, r)
			} else if controller.any != nil {
				response = controller.any(ctx, r)
			} else {
				response = router.routerRegistry.handler[router.notFoundHandler].any(context.WithValue(ctx, ERROR, errors.Errorf("action for method %q not found and no any fallback", req.Request().Method)), r)
				span.SetStatus(trace.Status{Code: trace.StatusCodeNotFound, Message: "action not found"})
			}

			return response
		},
	}

	result := chain.Next(ctx, req, rw)

	if router.sessionStore != nil {
		_, span := trace.StartSpan(ctx, "router/sessions/save")
		if err := router.sessionStore.Save(req.Request(), rw, gs); err != nil {
			router.logger.WithContext(ctx).Warn(err)
		}
		span.End()
	}

	var finalErr error
	if result != nil {
		_, span := trace.StartSpan(ctx, "router/responseApply")

		func() {
			defer func() {
				if err := panicToError(recover()); err != nil {
					finalErr = err
				}
			}()
			finalErr = result.Apply(ctx, rw)
		}()

		span.End()
	}

	// ensure that the session has been saved in the backend
	if router.sessionStore != nil {
		_, span := trace.StartSpan(ctx, "router/sessions/persist")
		if err := router.sessionStore.Save(req.Request(), emptyResponseWriter{}, gs); err != nil {
			router.logger.WithContext(ctx).Warn(err)
		}
		span.End()
	}

	for _, cb := range chain.postApply {
		cb(finalErr, result)
	}

	if finalErr != nil {
		func() {
			defer func() {
				if err := panicToError(recover()); err != nil {
					router.logger.WithContext(ctx).Error(err)
					rw.WriteHeader(http.StatusInternalServerError)
					_, _ = fmt.Fprintf(rw, "%+v", err)
				}
			}()

			if err := router.routerRegistry.handler[router.errorHandler].any(context.WithValue(ctx, ERROR, finalErr), req).Apply(ctx, rw); err != nil {
				router.logger.WithContext(ctx).Error(err)
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(rw, "%+v", err)
			}
		}()
	}
}

func dataParams(params map[interface{}]interface{}) RequestParams {
	vars := make(map[string]string, len(params))

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
			return c.data(ctx, r, dataParams(params))
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
