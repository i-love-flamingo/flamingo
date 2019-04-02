package responder

import (
	"net/http"

	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

// RedirectAware allows a controller to issue a 302 redirect
type (
	// RedirectAware trait
	RedirectAware interface {
		Redirect(name string, args map[string]string) web.Response
		RedirectURL(url string) web.Response
		RedirectPermanent(name string, args map[string]string) web.Response
		RedirectPermanentURL(url string) web.Response
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

// RedirectURL returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectURL(url string) web.Response {
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

// RedirectPermanentURL returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectPermanentURL(url string) web.Response {
	return &web.RedirectResponse{
		Status:   http.StatusMovedPermanently,
		Location: url,
	}
}
