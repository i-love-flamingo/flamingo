package responder

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/template"
	"net/http"
)

// RenderAware allows pug_template rendering
type RenderAware struct {
	Router *router.Router  `inject:""`
	Engine template.Engine `inject:""`
}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *RenderAware) Render(context web.Context, tpl string, data interface{}) *web.ContentResponse {
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        r.Engine.Render(context, tpl, data),
		ContentType: "text/html; charset=utf-8",
	}
}
