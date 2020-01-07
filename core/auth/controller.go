package auth

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

// controller manages login and callback requests
type controller struct {
	service *WebIdentityService
}

// Inject WebIdentityService dependency
func (c *controller) Inject(service *WebIdentityService) {
	c.service = service
}

// Callback is called e.g. for OIDC
func (c *controller) Callback(ctx context.Context, request *web.Request) web.Result {
	return c.service.callback(ctx, request)
}

// Login starts an authenticate for flow
func (c *controller) Login(ctx context.Context, request *web.Request) web.Result {
	return c.service.AuthenticateFor(request.Params["broker"], ctx, request)
}
