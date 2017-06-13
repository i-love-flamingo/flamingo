package controller

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

// Redirect Default controller
type Redirect struct {
	Responder responder.RedirectAware `inject:""`
}

// Redirect `to` a controller, all other params are passed on
func (redirect *Redirect) Redirect(ctx web.Context) web.Response {
	var params = ctx.ParamAll()
	var to = ctx.MustParam1("to")

	delete(params, "to")

	return redirect.Responder.Redirect(to, params)
}

// RedirectURL redirects to a url
func (redirect *Redirect) RedirectURL(ctx web.Context) web.Response {
	return redirect.Responder.RedirectURL(ctx.MustParam1("url"))
}

// RedirectPermanent is the same as Redirect but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanent(ctx web.Context) web.Response {
	var params = ctx.ParamAll()
	var to = ctx.MustParam1("to")

	delete(params, "to")

	return redirect.Responder.RedirectPermanent(to, params)
}

// RedirectPermanentURL is the same as RedirectURL but with a HTTP permanent redirect
func (redirect *Redirect) RedirectPermanentURL(ctx web.Context) web.Response {
	return redirect.Responder.RedirectPermanentURL(ctx.MustParam1("url"))
}
