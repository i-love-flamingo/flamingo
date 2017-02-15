package responder

import (
	"flamingo/core/core/app"
	"flamingo/core/core/app/web"
	"flamingo/core/core/template"
	"net/http"
)

type RenderAware struct {
	App *app.App `inject:""`
}

func (r *RenderAware) Render(context web.Context, tpl string, data interface{}) web.Response {
	return web.ContentResponse{
		Status:      http.StatusOK,
		Body:        template.Render(r.App, context, tpl, data),
		ContentType: "text/html",
	}
}
