package web

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"go.opencensus.io/trace"
)

type (
	// ReverseRouter allows to retrieve urls for controller
	ReverseRouter interface {
		// Relative returns a root-relative URL, starting with `/`
		// if to starts with "/" it will be used as the target, instead of resolving the URL
		Relative(to string, params map[string]string) (*url.URL, error)
		// Absolute returns an absolute URL, with scheme and host.
		// It takes the request to construct as many information as possible
		// if to starts with "/" it will be used as the target, instead of resolving the URL
		Absolute(r *Request, to string, params map[string]string) (*url.URL, error)
	}

	filterProvider    func() []Filter
	routesProvider    func() []RoutesModule
	responderProvider func() *Responder

	// Router represents actual implementation of ReverseRouter interface
	Router struct {
		base              *url.URL
		external          *url.URL
		eventRouter       flamingo.EventRouter
		filterProvider    filterProvider
		routesProvider    routesProvider
		logger            flamingo.Logger
		routerRegistry    *RouterRegistry
		configArea        *config.Area
		sessionStore      *SessionStore
		sessionName       string
		responderProvider responderProvider
	}
)

const (
	// FlamingoError is the Controller name for errors
	FlamingoError = "flamingo.error"
	// FlamingoNotfound is the Controller name for 404 notfound
	FlamingoNotfound = "flamingo.notfound"
)

// Inject dependencies
func (r *Router) Inject(
	cfg *struct {
		// base url configuration
		Scheme      string `inject:"config:flamingo.router.scheme,optional"`
		Host        string `inject:"config:flamingo.router.host,optional"`
		Path        string `inject:"config:flamingo.router.path,optional"`
		External    string `inject:"config:flamingo.router.external,optional"`
		SessionName string `inject:"config:flamingo.session.name,optional"`
	},
	sessionStore *SessionStore,
	eventRouter flamingo.EventRouter,
	filterProvider filterProvider,
	routesProvider routesProvider,
	logger flamingo.Logger,
	configArea *config.Area,
	responderProvider responderProvider,
) {
	r.base = &url.URL{
		Scheme: cfg.Scheme,
		Host:   cfg.Host,
		Path:   path.Join("/", cfg.Path, "/"),
	}

	if e, err := url.Parse(cfg.External); cfg.External != "" && err == nil {
		r.external = e
	} else if cfg.External != "" {
		r.logger.Warn("External URL Error: ", err)
	}

	r.eventRouter = eventRouter
	r.filterProvider = filterProvider
	r.routesProvider = routesProvider
	r.logger = logger
	r.configArea = configArea
	r.sessionStore = sessionStore
	r.sessionName = "flamingo"
	if cfg.SessionName != "" {
		r.sessionName = cfg.SessionName
	}
	r.responderProvider = responderProvider
}

// Handler creates and returns new instance of http.Handler interface
func (r *Router) Handler() http.Handler {
	r.routerRegistry = NewRegistry()

	if r.configArea != nil {
		for _, route := range r.configArea.Routes {
			r.routerRegistry.MustRoute(route.Path, route.Controller)
			if route.Name != "" {
				r.routerRegistry.Alias(route.Name, route.Controller)
			}
		}
	}

	for _, m := range r.routesProvider() {
		m.Routes(r.routerRegistry)
	}

	for _, handler := range r.routerRegistry.routes {
		if _, ok := r.routerRegistry.handler[handler.handler]; !ok {
			panic(fmt.Errorf("the handler %q has no controller, registered for path %q", handler.handler, handler.path.path))
		}
	}

	if r.responderProvider == nil {
		r.responderProvider = func() *Responder { return new(Responder) }
	}

	return &handler{
		routerRegistry: r.routerRegistry,
		filter:         r.filterProvider(),
		eventRouter:    r.eventRouter,
		logger:         r.logger.WithField(flamingo.LogKeyModule, "web").WithField(flamingo.LogKeyCategory, "handler"),
		sessionStore:   r.sessionStore,
		sessionName:    r.sessionName,
		prefix:         strings.TrimRight(r.Base().Path, "/"),
		responder:      r.responderProvider(),
	}
}

// ListenAndServe starts flamingo server
func (r *Router) ListenAndServe(addr string) error {
	r.eventRouter.Dispatch(context.Background(), &flamingo.ServerStartEvent{Port: addr})
	defer r.eventRouter.Dispatch(context.Background(), &flamingo.ServerShutdownEvent{})

	return http.ListenAndServe(addr, r.Handler())
}

// Base returns full base urls, containing scheme, domain and base path
func (r *Router) Base() *url.URL {
	if r.base == nil {
		return new(url.URL)
	}
	return r.base
}

// URL returns returns a root-relative URL, starting with `/`
// Deprecated: use Relative instead
func (r *Router) URL(to string, params map[string]string) (*url.URL, error) {
	return r.Relative(to, params)
}

func (r *Router) relative(to string, params map[string]string) (string, error) {
	if to == "" {
		return "", nil
	}

	if to[0] == '/' {
		return to, nil
	}

	p, err := r.routerRegistry.Reverse(to, params)
	return strings.TrimLeft(p, "/"), err
}

// Relative returns a root-relative URL, starting with `/`
func (r *Router) Relative(to string, params map[string]string) (*url.URL, error) {
	if to == "" {
		relativePath := r.Base().Path
		if r.external != nil {
			relativePath = r.external.Path
		}

		return &url.URL{
			Path: path.Join(relativePath),
		}, nil
	}

	p, err := r.relative(to, params)
	if err != nil {
		return nil, err
	}

	basePath := r.Base().Path
	if r.external != nil {
		basePath = r.external.Path
	}

	return url.Parse(path.Join("/", basePath, p, "/"))
}

// Absolute returns an absolute URL, with scheme and host.
// It takes the request to construct as many information as possible
func (r *Router) Absolute(req *Request, to string, params map[string]string) (*url.URL, error) {
	if r.external != nil {
		e := *r.external
		p, err := r.relative(to, params)
		if err != nil {
			return nil, err
		}
		e.Path = path.Join(e.Path, p)
		return &e, nil
	}

	scheme := r.Base().Scheme
	host := r.Base().Host

	if scheme == "" {
		if req != nil && req.request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	if host == "" && req != nil {
		host = req.request.Host
	}

	u, err := r.Relative(to, params)
	if err != nil {
		return u, err
	}

	u.Scheme = scheme
	u.Host = host
	return u, nil
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
func (r *Router) Data(ctx context.Context, handler string, params map[interface{}]interface{}) interface{} {
	ctx, span := trace.StartSpan(ctx, "flamingo/router/data")
	span.Annotate(nil, handler)
	defer span.End()

	req := RequestFromContext(ctx)

	if c, ok := r.routerRegistry.handler[handler]; ok {
		if c.data != nil {
			return c.data(ctx, req, dataParams(params))
		}
		err := fmt.Errorf("%q is not a data Controller", handler)
		r.logger.Error(err)
		panic(err)
	}
	err := fmt.Errorf("data Controller %q not found", handler)
	r.logger.Error(err)
	panic(err)
}
