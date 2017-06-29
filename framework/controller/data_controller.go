package controller

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

type (
	// DataController registers a route to allow external tools/ajax to retrieve data Handler
	DataController struct {
		responder.JSONAware `inject:""`
		Router              *router.Router `inject:""`
	}

	// SessionFlashController takes care of supported flash messages
	SessionFlashController struct{}

	// FlashMessage contains a type and a printable message
	FlashMessage struct {
		Type    string
		Message interface{}
	}
)

// Get Handler registered at /_flamingo/json/{Handler} and return's the call to Get()
func (gc *DataController) Get(c web.Context) web.Response {
	return gc.JSON(gc.Router.Get(c.MustParam1("handler"), c))
}

func getMessages(c web.Context, typ string) (messages []interface{}) {
	for _, flash := range c.Session().Flashes(typ) {
		messages = append(messages, FlashMessage{Type: typ, Message: flash})
	}
	return
}

// Data Controller for sessionflashcontroller
func (sfc *SessionFlashController) Data(c web.Context) interface{} {
	var messages []interface{}

	messages = append(messages, getMessages(c, web.ERROR)...)
	messages = append(messages, getMessages(c, web.WARNING)...)
	messages = append(messages, getMessages(c, web.INFO)...)

	return messages
}
