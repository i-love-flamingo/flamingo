package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

// Render controller
type Render struct {
	responder *web.Responder
}

// Inject *web.Responder
func (controller *Render) Inject(responder *web.Responder) {
	controller.responder = responder
}

// Render responder
func (controller *Render) Render(ctx context.Context, request *web.Request) web.Result {
	return controller.responder.Render(request.Params["tpl"], nil)
}
