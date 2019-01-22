package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/router"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/web/responder"
	"github.com/pkg/errors"
)

type (
	// Error controller
	Error struct {
		Responder responder.ErrorAware `inject:""`
	}
)

// Error responder
func (controller *Error) Error(ctx context.Context, request *web.Request) web.Response {
	var err error
	if ctx.Value(router.ERROR) != nil {
		err = ctx.Value(router.ERROR).(error)
	} else {
		err = errors.New("no error found in provided context")
	}
	return controller.Responder.Error(ctx, err)
}

// NotFound responder
func (controller *Error) NotFound(ctx context.Context, request *web.Request) web.Response {
	var err error
	if ctx.Value(router.ERROR) != nil {
		err = ctx.Value(router.ERROR).(error)
	} else {
		err = errors.New("no error found in provided context")
	}
	return controller.Responder.ErrorNotFound(ctx, err)
}
