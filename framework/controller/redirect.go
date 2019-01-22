package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/web/responder"
)

// Redirect Default controller
type Redirect struct {
	Responder responder.RedirectAware `inject:""`
}

// Redirect `to` a controller, all other params are passed on
func (redirect *Redirect) Redirect(ctx context.Context, request *web.Request) web.Response {
	params := request.ParamAll()
	to := request.MustParam1("to")

	delete(params, "to")

	return redirect.Responder.Redirect(to, params)
}

// RedirectURL redirects to a url
func (redirect *Redirect) RedirectURL(ctx context.Context, request *web.Request) web.Response {
	return redirect.Responder.RedirectURL(request.MustParam1("url"))
}

// RedirectPermanent is the same as Redirect but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanent(ctx context.Context, request *web.Request) web.Response {
	params := request.ParamAll()
	to := request.MustParam1("to")

	delete(params, "to")

	return redirect.Responder.RedirectPermanent(to, params)
}

// RedirectPermanentURL is the same as RedirectURL but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanentURL(ctx context.Context, request *web.Request) web.Response {
	return redirect.Responder.RedirectPermanentURL(request.MustParam1("url"))
}
