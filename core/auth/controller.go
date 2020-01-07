package auth

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

// controller manages login and callback requests
type controller struct {
	service   *WebIdentityService
	responder *web.Responder
}

// Inject WebIdentityService dependency
func (c *controller) Inject(service *WebIdentityService, responder *web.Responder) {
	c.service = service
	c.responder = responder
}

// Callback is called e.g. for OIDC
func (c *controller) Callback(ctx context.Context, request *web.Request) web.Result {
	return c.service.callback(ctx, request)
}

// Login starts an authenticate for flow
func (c *controller) Login(ctx context.Context, request *web.Request) web.Result {
	return c.service.AuthenticateFor(request.Params["broker"], ctx, request)
}

// LogoutAll removes all identities
func (c *controller) LogoutAll(ctx context.Context, request *web.Request) web.Result {
	c.service.Logout(ctx, request)
	return c.responder.RouteRedirect("", nil)
}

// Logout removes one identity
func (c *controller) Logout(ctx context.Context, request *web.Request) web.Result {
	c.service.LogoutFor(request.Params["broker"], ctx, request)
	return c.responder.RouteRedirect("", nil)
}
