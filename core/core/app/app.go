package app

import (
	"flamingo/core/core/app/context"
	"flamingo/core/core/app/web"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/color"
)

type (
	// Controller defines a web controller
	Controller interface{}

	GETController interface {
		Get(c web.Context) web.Response
	}

	POSTController interface {
		Post(c web.Context) web.Response
	}

	DataController interface {
		Data(c web.Context) interface{}
	}

	AppAwareInterface interface {
		SetApp(*App)
	}

	FixRoute struct {
		Handler string
		Params  map[string]string
	}

	// App defines the basic multiplexer
	App struct {
		router    *mux.Router
		routes    map[string]string
		handler   map[string]interface{}
		fixroutes map[string]FixRoute
		Debug     bool
		base      *url.URL
		log       *log.Logger
	}

	ResponseWriter struct {
		http.ResponseWriter
		status int
		size   int
	}
)

func (r *ResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *ResponseWriter) Write(data []byte) (int, error) {
	l, e := r.ResponseWriter.Write(data)
	r.size += l
	return l, e
}

func (r *ResponseWriter) WriteHeader(h int) {
	r.status = h
	r.ResponseWriter.WriteHeader(h)
}

// New factory for web router
func New(ctx *context.Context) *App {
	a := &App{
	//		fixroutes: make(map[string]FixRoute),
	}

	a.router = mux.NewRouter()
	a.routes = ctx.Routes
	a.handler = ctx.Handler
	a.base, _ = url.Parse("scheme://" + ctx.BaseUrl)
	a.log = log.New(os.Stdout, "["+ctx.Name+"] ", 0)

	for route, name := range ctx.Routes {
		a.log.Println("Register", name, "at", route)
		if _, ok := ctx.Handler[name]; !ok {
			panic("no handler for" + name)
		}
		a.router.Handle(route, a.handle(ctx.Handler[name])).Name(name)
	}

	return a
}

func fixid(handler string, params ...string) string {
	return handler + "!!" + strings.Join(params, "!!!")
}

// Router generates a http.Handler
func (r *App) Router() *mux.Router {
	return r.router
}

func (r *App) Url(name string, params ...string) *url.URL {
	u, err := r.router.Get(name).URL(params...)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(r.base.Path, u.Path)
	return u
}

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w = &ResponseWriter{ResponseWriter: w}
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(500)
			if a.Debug {
				w.Write(debug.Stack())
			}
		}
		if a.Debug {
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

			a.log.Printf(cp("%03d | % 8s | % 15s | %s"), ww.status, req.Method, time.Since(start), req.RequestURI)
		}
	}()

	a.router.ServeHTTP(w, req)
}

func (r *App) handle(c Controller) http.Handler {
	if c, ok := c.(AppAwareInterface); ok {
		c.SetApp(r)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := web.ContextFromRequest(w, req)

		if req.Method == http.MethodGet {
			if c, ok := c.(GETController); ok {
				c.Get(ctx).Apply(w)
				return
			}
		} else if req.Method == http.MethodPost {
			if c, ok := c.(POSTController); ok {
				c.Post(ctx).Apply(w)
				return
			}
		}

		if c, ok := c.(http.Handler); ok {
			c.ServeHTTP(w, req)
			return
		}

		w.WriteHeader(404)
		w.Write([]byte("404 page not found (no handler)"))
		//panic("cannot serve " + req.RequestURI)
	})
}

func (a *App) Get(ctx web.Context, handler string) interface{} {
	if c, ok := a.handler[handler]; ok {
		if c, ok := c.(DataController); ok {
			return c.Data(ctx)
		}
		panic("not a data controller")
	}
	panic("not a handler")
}
