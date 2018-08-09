package controller

import (
	"context"

	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// DataController registers a route to allow external tools/ajax to retrieve data Handler
	DataController struct {
		responder.JSONAware
		router *router.Router
	}

	// SessionFlashController takes care of supported flash messages
	SessionFlashController struct{}

	// FlashMessage contains a type and a printable message
	FlashMessage struct {
		Type    string
		Message interface{}
	}
)

func (gc *DataController) Inject(aware responder.JSONAware, router *router.Router) {
	gc.JSONAware = aware
	gc.router = router
}

// Get Handler registered at /_flamingo/json/{Handler} and return's the call to Get()
func (gc *DataController) Get(c context.Context, r *web.Request) web.Response {
	return gc.JSON(gc.router.Data(c, r.MustParam1("handler"), nil))
}

func getMessages(r *web.Request, typ string) (messages []interface{}) {
	for _, flash := range r.Session().Flashes(typ) {
		messages = append(messages, FlashMessage{Type: typ, Message: flash})
	}
	return
}

// Data Controller for sessionflashcontroller
func (sfc *SessionFlashController) Data(c context.Context, r *web.Request) interface{} {
	var messages []interface{}

	messages = append(messages, getMessages(r, web.ERROR)...)
	messages = append(messages, getMessages(r, web.WARNING)...)
	messages = append(messages, getMessages(r, web.INFO)...)

	return messages
}
