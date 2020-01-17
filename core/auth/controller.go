package auth

import (
	"context"
	"errors"

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
	if resp := c.service.callback(ctx, request); resp != nil {
		return resp
	}
	return c.responder.NotFound(errors.New("broker callback found"))
}

// Login starts an authenticate for flow
func (c *controller) Login(ctx context.Context, request *web.Request) web.Result {
	redirecturl, ok := request.Params["redirecturl"]
	if !ok || redirecturl == "" {
		redirecturl = request.Request().Referer()
	}
	request.Params["redirecturl"] = redirecturl

	if resp := c.service.AuthenticateFor(ctx, request.Params["broker"], request); resp != nil {
		return resp
	}
	return c.responder.NotFound(errors.New("broker login found"))
}

// LogoutAll removes all identities
func (c *controller) LogoutAll(ctx context.Context, request *web.Request) web.Result {
	c.service.Logout(ctx, request)
	return c.responder.RouteRedirect("", nil)
}

// Logout removes one identity
func (c *controller) Logout(ctx context.Context, request *web.Request) web.Result {
	c.service.LogoutFor(ctx, request.Params["broker"], request)
	return c.responder.RouteRedirect("", nil)
}
