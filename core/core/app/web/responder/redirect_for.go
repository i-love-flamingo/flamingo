package responder

import (
	"flamingo/core/core/app"
	"goaoe/core/web"
	"net/http"
)

type RedirectAware struct {
	app *app.App
}

func (r *RedirectAware) SetApp(app *app.App) {
	r.app = app
}

func (r *RedirectAware) Redirect(name string, args ...string) web.Response {
	url := r.app.Url(name, args...)

	return web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url.String(),
	}
}
