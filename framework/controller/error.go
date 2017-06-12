package controller

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

type (
	// Error controller
	Error struct {
		Responder responder.ErrorAware `inject:""`
	}
)

// Error responder
func (controller *Error) Error(ctx web.Context) web.Response {
	var err error
	if ctx.Value(router.ERROR) != nil {
		err = ctx.Value(router.ERROR).(error)
	}
	return controller.Responder.Error(ctx, err)
}

// NotFound responder
func (controller *Error) NotFound(ctx web.Context) web.Response {
	var err error
	if ctx.Value(router.ERROR) != nil {
		err = ctx.Value(router.ERROR).(error)
	}
	return controller.Responder.ErrorNotFound(ctx, err)
}
