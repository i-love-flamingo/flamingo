package core

import (
	"flamingo/core/context"
	"flamingo/core/web"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"os"

	"time"

	"github.com/gorilla/mux"
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
)

// NewApp factory for web router
func NewApp(ctx *context.Context) *App {
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
	a.router.ServeHTTP(w, req)
}

func (r *App) handle(c Controller) http.Handler {
	if c, ok := c.(AppAwareInterface); ok {
		c.SetApp(r)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(500)
			}
			r.log.Println(req.RequestURI, time.Since(start))
		}()

		ctx := web.ContextFromRequest(req)

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

		panic("cannot serve " + req.RequestURI)
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
