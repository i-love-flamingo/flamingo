package responder

import (
	"flamingo/core"
	"flamingo/core/web"
	"net/http"
)

type RedirectAware struct {
	app *core.App
}

func (r *RedirectAware) SetApp(app *core.App) {
	r.app = app
}

func (r *RedirectAware) Redirect(name string, args ...string) web.Response {
	url := r.app.Url(name, args...)

	return web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url.String(),
	}
}
