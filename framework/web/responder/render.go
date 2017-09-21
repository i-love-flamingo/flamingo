package responder

import (
	"bytes"
	"encoding/json"
	"flamingo/core/pugtemplate/pugjs"
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
		Engine template.Engine `inject:",optional"`
	}
)

var _ RenderAware = &FlamingoRenderAware{}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *FlamingoRenderAware) Render(context web.Context, tpl string, data interface{}) (response web.Response) {
	if d, err := context.Query1("debugdata"); err == nil && d != "" {
		return &web.JSONResponse{
			Data: pugjs.Convert(data),
		}
	}

	if r.Engine != nil {
		body, err := r.Engine.Render(context, tpl, data)
		if err != nil {
			panic(err)
		}
		response = &web.ContentResponse{
			Status:      http.StatusOK,
			Body:        body,
			ContentType: "text/html; charset=utf-8",
		}
	} else {
		body, err := json.Marshal(pugjs.Convert(data))
		if err != nil {
			panic(err)
		}
		response = &web.ContentResponse{
			Status:      http.StatusOK,
			Body:        bytes.NewReader(body),
			ContentType: "application/json; charset=utf-8",
		}
	}
	return
}
