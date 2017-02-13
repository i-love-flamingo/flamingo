package responder

import (
	"flamingo/core"
	"flamingo/core/template"
	"flamingo/core/web"
	"net/http"
)

type RenderAware struct {
	app *core.App
}

func (r *RenderAware) SetApp(app *core.App) {
	r.app = app
}

func (r *RenderAware) Render(context web.Context, tpl string) web.Response {
	return web.ContentResponse{
		Status:      http.StatusOK,
		Body:        template.Render(r.app, tpl, nil),
		ContentType: "text/html",
	}
}
