package responder

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"net/http"
)

// RedirectAware allows a controller to issue a 302 redirect
type (
	// RedirectAware trait
	RedirectAware interface {
		Redirect(name string, args map[string]string) web.Response
		RedirectUrl(url string) web.Response
		RedirectPermanent(name string, args map[string]string) web.Response
		RedirectPermanentUrl(url string) web.Response
	}

	// FlamingoRedirectAware flamingo's redirect aware
	FlamingoRedirectAware struct {
		Router *router.Router `inject:""`
	}
)

var _ RedirectAware = &FlamingoRedirectAware{}

// Redirect returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) Redirect(name string, args map[string]string) web.Response {
	u := r.Router.URL(name, args)

	return &web.RedirectResponse{
		Status:   http.StatusFound,
		Location: u.String(),
	}
}

// RedirectUrl returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectUrl(url string) web.Response {
	return &web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url,
	}
}

// RedirectPermanent returns a web.RedirectPermanentResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectPermanent(name string, args map[string]string) web.Response {
	u := r.Router.URL(name, args)

	return &web.RedirectResponse{
		Status:   http.StatusMovedPermanently,
		Location: u.String(),
	}
}

// RedirectPermantentUrl returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectPermanentUrl(url string) web.Response {
	return &web.RedirectResponse{
		Status:   http.StatusMovedPermanently,
		Location: url,
	}
}
