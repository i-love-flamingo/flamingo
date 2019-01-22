package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// SessionFlashController takes care of supported flash messages
	SessionFlashController struct{}

	// FlashMessage contains a type and a printable message
	FlashMessage struct {
		Type    string
		Message interface{}
	}
)

func getMessages(r *web.Request) (messages []interface{}) {
	for _, flash := range r.Session().Flashes() {
		messages = append(messages, FlashMessage{Message: flash})
	}
	return
}

// Data Controller for sessionflashcontroller
func (sfc *SessionFlashController) Data(c context.Context, r *web.Request, _ web.RequestParams) interface{} {
	return getMessages(r)
}
