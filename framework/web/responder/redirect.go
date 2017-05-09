package responder

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"net/http"
)

// RedirectAware allows a controller to issue a 302 redirect
type RedirectAware struct {
	Router *router.Router `inject:""`
}

// Redirect returns a web.RedirectResponse with the proper URL
func (r *RedirectAware) Redirect(name string, args ...string) web.Response {
	u := r.Router.URL(name, args...)

	return &web.RedirectResponse{
		Status:   http.StatusFound,
		Location: u.String(),
	}
}

// RedirectUrl returns a web.RedirectResponse with the proper URL
func (r *RedirectAware) RedirectUrl(url string) web.Response {
	return &web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url,
	}
}

// RedirectPermanent returns a web.RedirectPermanentResponse with the proper URL
func (r *RedirectAware) RedirectPermanent(name string, args ...string) web.Response {
	u := r.Router.URL(name, args...)

	return &web.RedirectResponse{
		Status:   http.StatusMovedPermanently,
		Location: u.String(),
	}
}

// RedirectPermantentUrl returns a web.RedirectResponse with the proper URL
func (r *RedirectAware) RedirectPermanentUrl(url string) web.Response {
	return &web.RedirectResponse{
		Status:   http.StatusMovedPermanently,
		Location: url,
	}
}
