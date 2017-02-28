package responder

import (
	"flamingo/core/flamingo/router"
	"flamingo/core/flamingo/web"
	"flamingo/core/template"
	"net/http"
)

// RenderAware allows pug_template rendering
type RenderAware struct {
	Router *router.Router  `inject:""`
	Engine template.Engine `inject:""`
}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *RenderAware) Render(context web.Context, tpl string, data interface{}) web.Response {
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        r.Engine.Render(context, tpl, data),
		ContentType: "text/html",
	}
}
