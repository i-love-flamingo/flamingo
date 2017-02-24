package controller

import (
	"flamingo/core/flamingo/router"
	"flamingo/core/flamingo/web"
	"flamingo/core/flamingo/web/responder"
)

type (
	// GetController registers a route to allow external tools/ajax to retrieve data Handler
	DataController struct {
		Router               *router.Router `inject:""`
		*responder.JsonAware `inject:""`
	}
)

// Get Handler registered at /_flamingo/json/{Handler} and return's the call to Get()
func (gc *DataController) Get(c web.Context) web.Response {
	return gc.Json(gc.Router.Get(c.Param1("Handler"), c))
}
