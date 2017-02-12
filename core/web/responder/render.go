package responder

import (
	"flamingo/core"
	"flamingo/core/template"
	"flamingo/core/web"
	"net/http"
)

type RenderTemplate struct {
	app *core.App
}

func (r *RenderTemplate) SetApp(app *core.App) {
	r.app = app
}

func (r *RenderTemplate) RenderResponse(context web.Context, tpl string) web.Response {
	return web.ContentResponse{
		Status:      http.StatusOK,
		Body:        template.Render(context, r.app, tpl, nil),
		ContentType: "text/html",
	}
}
