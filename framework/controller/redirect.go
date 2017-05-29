package controller

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

// Redirect Default controller
type Redirect struct {
	Responder *responder.RedirectAware `inject:""`
}

// Redirect `to` a controller, all other params are passed on
func (redirect *Redirect) Redirect(ctx web.Context) web.Response {
	var params = ctx.ParamAll()
	var to = ctx.MustParam1("to")

	delete(params, "to")

	return redirect.Responder.Redirect(to, params)
}

// RedirectUrl redirects to a url
func (redirect *Redirect) RedirectUrl(ctx web.Context) web.Response {
	return redirect.Responder.RedirectUrl(ctx.MustParam1("url"))
}

// RedirectPermanent is the same as Redirect but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanent(ctx web.Context) web.Response {
	var params = ctx.ParamAll()
	var to = ctx.MustParam1("to")

	delete(params, "to")

	return redirect.Responder.RedirectPermanent(to, params)
}

// RedirectPermanentUrl is the same as RedirectUrl but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanentUrl(ctx web.Context) web.Response {
	return redirect.Responder.RedirectPermanentUrl(ctx.MustParam1("url"))
}
