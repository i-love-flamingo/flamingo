package responder

import (
	"flamingo/core/core/app"
	"flamingo/core/core/app/web"
	"flamingo/core/core/template"
	"net/http"
)

type RenderAware struct {
	app *app.App
}

func (r *RenderAware) SetApp(app *app.App) {
	r.app = app
}

func (r *RenderAware) Render(context web.Context, tpl string) web.Response {
	return web.ContentResponse{
		Status:      http.StatusOK,
		Body:        template.Render(r.app, context, tpl, nil),
		ContentType: "text/html",
	}
}
