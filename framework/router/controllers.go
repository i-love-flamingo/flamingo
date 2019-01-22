package router

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Action defines an explicit http action
	Action func(ctx context.Context, req *web.Request) web.Response

	// DataAction is a method called which does not return the web response itself, but data instead
	DataAction func(ctx context.Context, req *web.Request) interface{}

	// ControllerOption defines a type for Controller options
	// todo still usable?
	ControllerOption string

	// ControllerOptionAware is an interface for Controller which want to interact with filter
	// todo still usable?
	ControllerOptionAware interface {
		CheckOption(option ControllerOption) bool
	}
)

// HTTPAction wraps a default http.Handler to a flamingo router action
func HTTPAction(handler http.Handler) Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		r := &web.ServeHTTPResponse{
			VerboseResponseWriter: ctx.Value("rw").(*web.VerboseResponseWriter),
		}
		handler.ServeHTTP(r, req.Request())
		return r
	}
}
