package controller

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

type (
	// DataController registers a route to allow external tools/ajax to retrieve data Handler
	DataController struct {
		Router               *router.Router `inject:""`
		*responder.JSONAware `inject:""`
	}
)

// Get Handler registered at /_flamingo/json/{Handler} and return's the call to Get()
func (gc *DataController) Get(c web.Context) web.Response {
	return gc.JSON(gc.Router.Get(c.Param1("Handler"), c))
}
