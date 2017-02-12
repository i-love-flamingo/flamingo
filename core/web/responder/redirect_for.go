package responder

import (
	"flamingo/core"
	"flamingo/core/web"
	"net/http"
)

type RedirectFor struct {
	app *core.App
}

func (r *RedirectFor) SetApp(app *core.App) {
	r.app = app
}

func (r *RedirectFor) Response(name string, args ...string) web.Response {
	url := r.app.Url(name, args...)

	return web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url.String(),
	}
}
