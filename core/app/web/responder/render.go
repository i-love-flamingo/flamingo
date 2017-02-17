package responder

import (
	"flamingo/core/app"
	"flamingo/core/app/web"
	"flamingo/core/packages/pug-template"
	"net/http"
)

// RenderAware allows pug-template rendering
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
