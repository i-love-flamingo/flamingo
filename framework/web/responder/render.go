package responder

import (
	"flamingo/framework/router"
	"flamingo/framework/template"
	"flamingo/framework/web"
	"net/http"
)

type (
	// RenderAware controller trait
	RenderAware interface {
		Render(context web.Context, tpl string, data interface{}) web.Response
	}

	// FlamingoRenderAware allows pug_template rendering
	FlamingoRenderAware struct {
		Router *router.Router  `inject:""`
		Engine template.Engine `inject:""`
	}
)

var _ RenderAware = &FlamingoRenderAware{}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *FlamingoRenderAware) Render(context web.Context, tpl string, data interface{}) web.Response {
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        r.Engine.Render(context, tpl, data),
		ContentType: "text/html; charset=utf-8",
	}
}
