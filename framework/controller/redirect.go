package controller

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/framework/web"
)

// Redirect Default controller
type Redirect struct {
	responder *web.Responder
}

// Inject *web.Responder
func (redirect *Redirect) Inject(responder *web.Responder) {
	redirect.responder = responder
}

// Redirect `to` a controller, all other params are passed on
func (redirect *Redirect) Redirect(ctx context.Context, request *web.Request) web.Result {
	to := request.Params["to"]
	delete(request.Params, "to")

	return redirect.responder.RouteRedirect(to, request.Params)
}

// RedirectURL redirects to a url
func (redirect *Redirect) RedirectURL(ctx context.Context, request *web.Request) web.Result {
	target, err := url.Parse(request.Params["url"])
	if err != nil {
		return redirect.responder.ServerError(err)
	}
	return redirect.responder.URLRedirect(target)
}

// RedirectPermanent is the same as Redirect but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanent(ctx context.Context, request *web.Request) web.Result {
	to := request.Params["to"]

	delete(request.Params, "to")

	return redirect.responder.RouteRedirect(to, request.Params).Permanent()
}

// RedirectPermanentURL is the same as RedirectURL but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanentURL(ctx context.Context, request *web.Request) web.Result {
	target, err := url.Parse(request.Params["url"])
	if err != nil {
		return redirect.responder.ServerError(err)
	}
	return redirect.responder.URLRedirect(target).Permanent()
}
