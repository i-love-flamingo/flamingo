package auth

import (
	"context"
	"errors"
	"net/url"

	"flamingo.me/flamingo/v3/framework/web"
)

// controller manages login and callback requests
type controller struct {
	service       *WebIdentityService
	responder     *web.Responder
	reverseRouter web.ReverseRouter
}

// Inject WebIdentityService dependency
func (c *controller) Inject(service *WebIdentityService, responder *web.Responder, reverseRouter web.ReverseRouter) {
	c.service = service
	c.responder = responder
	c.reverseRouter = reverseRouter
}

// Callback is called e.g. for OIDC
func (c *controller) Callback(ctx context.Context, request *web.Request) web.Result {
	if resp := c.service.callback(ctx, request); resp != nil {
		return resp
	}
	return c.responder.NotFound(errors.New("broker for callback not found"))
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
	return c.responder.NotFound(errors.New("broker for login not found"))
}

// LogoutAll removes all identities
func (c *controller) LogoutAll(ctx context.Context, request *web.Request) web.Result {
	return c.service.Logout(ctx, request, nil)
}

// Logout removes one identity
func (c *controller) Logout(ctx context.Context, request *web.Request) web.Result {
	return c.service.LogoutFor(ctx, request.Params["broker"], request, nil)
}

// LogoutCallback redirects to the next upcoming redirect url
func (c *controller) LogoutCallback(ctx context.Context, request *web.Request) web.Result {
	redirects := c.service.getLogoutRedirects(request)
	if len(redirects) == 0 {
		if postRedirect, ok := request.Session().Load("core.auth.logoutredirect"); ok {
			request.Session().Delete("core.auth.logoutredirect")
			return c.responder.URLRedirect(postRedirect.(*url.URL))
		}
		return c.responder.RouteRedirect("", nil)
	}
	next := redirects[0]
	c.service.storeLogoutRedirects(request, redirects[1:])
	if postRedirect, err := c.reverseRouter.Absolute(request, "core.auth.logoutCallback", nil); err == nil {
		query := next.Query()
		query.Set("post_logout_redirect_uri", postRedirect.String())
		next.RawQuery = query.Encode()
	}
	return c.responder.URLRedirect(next)
}
