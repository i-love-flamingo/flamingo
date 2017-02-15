package app

import (
	"encoding/json"
	"flamingo/core/core/app/web"
	"io/ioutil"
)

type (
	// DataController is a controller used to retrieve data, such as user-information, basket
	// etc.
	// By default this will be handled by templates, but there is an out-of-the-box support
	// for JSON requests via /_flamingo/json/{name}, as well as their own route if defined.
	DataController interface {
		// Data is called for data requests
		Data(web.Context) interface{}
	}

	// DataHandler behaves the same as DataController, but just for direct callbacks
	DataHandler func(web.Context) interface{}

	// GetController registers a route to allow external tools/ajax to retrieve data handler
	GetController struct {
		App *App `inject:""`
	}
)

// Get is the ServeHTTP's equivalent for DataController and DataHandler
func (a *App) Get(handler string, ctx web.Context) interface{} {
	if c, ok := a.handler[handler]; ok {
		if c, ok := c.(DataController); ok {
			return c.Data(ctx)
		}
		if c, ok := c.(func(web.Context) interface{}); ok {
			return c(ctx)
		}
		panic("not a data controller")
	} else if a.Debug { // mock...
		data, err := ioutil.ReadFile("frontend/src/mocks/" + handler + ".json")
		if err == nil {
			var res interface{}
			json.Unmarshal(data, &res)
			return res
		}
	}
	panic("not a handler: " + handler)
}

// GetHandler is registered at /_flamingo/json/{handler} and return's the call to Get()
func (gc *GetController) GetHandler(c web.Context) web.Response {
	return web.JsonResponse{
		Data: gc.App.Get(c.Param1("handler"), c),
	}
}
