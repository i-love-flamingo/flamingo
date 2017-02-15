package responder

import (
	"flamingo/core/core/app"
	"flamingo/core/core/app/web"
	"net/http"
)

type RedirectAware struct {
	App *app.App `inject:""`
}

func (r *RedirectAware) Redirect(name string, args ...string) web.Response {
	url := r.App.Url(name, args...)

	return web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url.String(),
	}
}
