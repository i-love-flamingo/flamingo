package flamingo

import (
	"flamingo/core/flamingo/context"
	"flamingo/core/flamingo/web"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/labstack/gommon/color"
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

	// Router defines the basic Router which is used for holding a context-scoped setup
	// This includes DI resolving etc
	Router struct {
		router   *mux.Router
		routes   map[string]string
		handler  map[string]interface{}
		Debug    bool
		base     *url.URL
		Logger   *log.Logger `inject:""`
		Sessions sessions.Store
	}

	// ResponseWriter shadows http.ResponseWriter and tracks written bytes and result status for logging
	ResponseWriter struct {
		http.ResponseWriter
		status int
		size   int
	}
)

// Writes calls http.ResponseWriter.Write and records the written bytes
func (r *ResponseWriter) Write(data []byte) (int, error) {
	l, e := r.ResponseWriter.Write(data)
	r.size += l
	return l, e
}

// WriteHeader call http.ResponseWriter.WriteHeader and records the status code
func (r *ResponseWriter) WriteHeader(h int) {
	r.status = h
	r.ResponseWriter.WriteHeader(h)
}

func (r *ResponseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := r.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return nil
}

// New factory for Router
// New creates the new flamingo, set's up handlers and routes and resolved the DI
func New(ctx *context.Context, serviceContainer *ServiceContainer) *Router {
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
	for p, name := range serviceContainer.routes {
		a.routes[name] = p
	}

	for p, name := range ctx.Routes {
		a.routes[name] = p
	}

	// set up handlers
	for name, handler := range serviceContainer.handler {
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

// Router returns the http.Handler
func (router *Router) Router() *mux.Router {
	return router.router
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
		extra := ""

		if err := recover(); err != nil {
			w.WriteHeader(500)
			if router.Debug {
				extra += fmt.Sprintf(`| Error: %s`, err)
				w.Write([]byte(fmt.Sprintln(err)))
				w.Write(debug.Stack())
			}
		}
		if router.Debug {
			ww := w.(*ResponseWriter)
			var cp func(msg interface{}, styles ...string) string
			switch {
			case ww.status >= 200 && ww.status < 300:
				cp = color.Green
			case ww.status >= 300 && ww.status < 400:
				cp = color.Blue
			case ww.status >= 400 && ww.status < 500:
				cp = color.Yellow
			case ww.status >= 500 && ww.status < 600:
				cp = color.Red
			default:
				cp = color.Black
			}

			if ww.Header().Get("Location") != "" {
				extra += "-> " + ww.Header().Get("Location")
			}
			router.Logger.Printf(cp("%03d | %-8s | % 15s | % 6d byte | %s %s"), ww.status, req.Method, time.Since(start), ww.size, req.RequestURI, extra)
		}
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
