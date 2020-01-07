package auth

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type Controller struct {
	service *WebIdentityService
}

func (c *Controller) Inject(service *WebIdentityService) {
	c.service = service
}

func (c *Controller) Callback(ctx context.Context, request *web.Request) web.Result {
	return c.service.callback(ctx, request)
}

func (c *Controller) Login(ctx context.Context, request *web.Request) web.Result {
	return c.service.AuthenticateFor(request.Params["broker"], ctx, request)
}
