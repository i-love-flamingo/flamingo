package responder

import (
	"net/http"

	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

// RedirectAware allows a controller to issue a 302 redirect
type (
	// RedirectAware trait
	RedirectAware interface {
		Redirect(name string, args map[string]string) web.Redirect
		RedirectURL(url string) web.Redirect
		RedirectPermanent(name string, args map[string]string) web.Redirect
		RedirectPermanentURL(url string) web.Redirect
	}

	// FlamingoRedirectAware flamingo's redirect aware
	FlamingoRedirectAware struct {
		Router *router.Router `inject:""`
	}
)

var _ RedirectAware = &FlamingoRedirectAware{}

// Redirect returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) Redirect(name string, args map[string]string) web.Redirect {
	u := r.Router.URL(name, args)

	return &web.RedirectResponse{
		BasicResponse: web.BasicResponse{
			Status: http.StatusFound,
		},
		Location: u.String(),
	}
}

// RedirectURL returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectURL(url string) web.Redirect {
	return &web.RedirectResponse{
		BasicResponse: web.BasicResponse{
			Status: http.StatusFound,
		},
		Location: url,
	}
}

// RedirectPermanent returns a web.RedirectPermanentResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectPermanent(name string, args map[string]string) web.Redirect {
	u := r.Router.URL(name, args)

	return &web.RedirectResponse{
		BasicResponse: web.BasicResponse{
			Status: http.StatusMovedPermanently,
		},
		Location: u.String(),
	}
}

// RedirectPermanentURL returns a web.RedirectResponse with the proper URL
func (r *FlamingoRedirectAware) RedirectPermanentURL(url string) web.Redirect {
	return &web.RedirectResponse{
		BasicResponse: web.BasicResponse{
			Status: http.StatusMovedPermanently,
		},
		Location: url,
	}
}
