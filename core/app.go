package core

import (
	"flamingo/core/context"
	"flamingo/core/web"
	"flamingo/core/web/responder"
	"log"
	"net/http"

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

	// App defines the basic multiplexer
	App struct {
		router *mux.Router
	}
)

// NewApp factory for web router
func NewApp() *App {
	return &App{}
}

// Router generates a http.Handler
func (r *App) Router(ctx *context.Context) *mux.Router {
	r.router = mux.NewRouter()

	for route, name := range ctx.Routes {
		log.Println("Register", name, "at", route)
		r.router.Handle(route, r.handle(ctx.Handler[name])).Name(name)
	}

	return r.router
}

func (r *App) handle(c Controller) http.Handler {
	if c, ok := c.(responder.RouterAware); ok {
		c.SetRouter(r.router)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
