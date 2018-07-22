package router

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/framework/web"
)

type (
	// Action defines an explicit http action
	Action func(ctx context.Context, req *web.Request) web.Response

	// DataAction is a method called which does not return the web response itself, but data instead
	DataAction func(ctx context.Context, req *web.Request) interface{}

	// Controller defines a web Controller
	// it is an interface{} as it can be served by multiple possible controllers,
	// such as generic GET/POST Controller, http.Handler, Handler-functions, etc.
	//
	// deprecated: superseeded by method based handling
	Controller interface{}

	// ControllerOption defines a type for Controller options
	ControllerOption string

	// ControllerOptionAware is an interface for Controller which want to interact with filter
	ControllerOptionAware interface {
		CheckOption(option ControllerOption) bool
	}

	// GETController is implemented by controllers which have a Get method
	// deprecated: use method based handler
	GETController interface {
		// Get is called for GET-Requests
		Get(web.Context) web.Response
	}

	// POSTController is implemented by controllers which have a Post method
	// deprecated: use method based handler
	POSTController interface {
		// Post is called for POST-Requests
		Post(web.Context) web.Response
	}

	// PUTController is implemented by controllers which have a Put method
	// deprecated: use method based handler
	PUTController interface {
		// Put is called for PUT-Requests
		Put(web.Context) web.Response
	}

	// DELETEController is implemented by controllers which have a Delete method
	// deprecated: use method based handler
	DELETEController interface {
		// Delete is called for DELETE-Requests
		Delete(web.Context) web.Response
	}

	// HEADController is implemented by controllers which have a Head method
	// deprecated: use method based handler
	HEADController interface {
		// Head is called for HEAD-Requests
		Head(web.Context) web.Response
	}

	// DataController is a Controller used to retrieve data, such as user-information, basket
	// etc.
	// By default this will be handled by templates, but there is an out-of-the-box support
	// for JSON requests via /_flamingo/json/{name}, as well as their own route if defined.
	DataController interface {
		// Data is called for data requests
		Data(web.Context) interface{}
	}

	// DataHandler behaves the same as DataController, but just for direct callbacks
	DataHandler func(web.Context) interface{}
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
