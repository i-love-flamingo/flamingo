package controller

import (
	"context"

	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// Render controller
	Render struct {
		Responder responder.RenderAware `inject:""`
	}
)

// Render responder
func (controller *Render) Render(ctx context.Context, request *web.Request) web.Response {
	return controller.Responder.Render(ctx, request.MustParam1("tpl"), nil)
}
