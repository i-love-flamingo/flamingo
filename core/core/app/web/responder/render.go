package responder

import (
	"flamingo/core/core/app"
	"flamingo/core/core/app/web"
	"flamingo/core/core/template"
	"net/http"
)

// RenderAware allows template rendering
type RenderAware struct {
	App *app.App `inject:""`
}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *RenderAware) Render(context web.Context, tpl string, data interface{}) web.Response {
	return web.ContentResponse{
		Status:      http.StatusOK,
		Body:        template.Render(r.App, context, tpl, data),
		ContentType: "text/html",
	}
}
