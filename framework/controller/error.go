package controller

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/web"
)

// Error controller
type Error struct {
	responder *web.Responder
}

// Inject *web.Responder
func (controller *Error) Inject(responder *web.Responder) {
	controller.responder = responder
}

// Error responder
func (controller *Error) Error(ctx context.Context, request *web.Request) web.Result {
	var err error
	if ctx.Value(web.RouterError) != nil {
		err = ctx.Value(web.RouterError).(error)
	} else {
		err = errors.New("no error found in provided context")
	}
	return controller.responder.ServerError(err)
}

// NotFound responder
func (controller *Error) NotFound(ctx context.Context, request *web.Request) web.Result {
	var err error
	if ctx.Value(web.RouterError) != nil {
		err = ctx.Value(web.RouterError).(error)
	} else {
		err = errors.New("no error found in provided context")
	}
	return controller.responder.NotFound(err)
}
