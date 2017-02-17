package responder

import (
	"flamingo/core/app"
	"flamingo/core/app/web"
	"net/http"
)

// RedirectAware allows a controller to issue a 302 redirect
type RedirectAware struct {
	App *app.App `inject:""`
}

// Redirect returns a web.RedirectResponse with the proper URL
func (r *RedirectAware) Redirect(name string, args ...string) web.Response {
	url := r.App.Url(name, args...)

	return web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url.String(),
	}
}
