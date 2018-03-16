package responder

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/template"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// RenderAware controller trait
	RenderAware interface {
		Render(context web.Context, tpl string, data interface{}) web.Response
		WithStatusCode(code int) RenderAware
	}

	// FlamingoRenderAware allows pug_template rendering
	FlamingoRenderAware struct {
		Router         *router.Router  `inject:""`
		Engine         template.Engine `inject:",optional"`
		httpStatusCode int
	}
)

var _ RenderAware = &FlamingoRenderAware{}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *FlamingoRenderAware) Render(context web.Context, tpl string, data interface{}) (response web.Response) {
	statusCode := http.StatusOK
	if r.httpStatusCode > 0 {
		statusCode = r.httpStatusCode
	}

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
			Status:      statusCode,
			Body:        body,
			ContentType: "text/html; charset=utf-8",
		}
	} else {
		body, err := json.Marshal(pugjs.Convert(data))
		if err != nil {
			panic(err)
		}
		response = &web.ContentResponse{
			Status:      statusCode,
			Body:        bytes.NewReader(body),
			ContentType: "application/json; charset=utf-8",
		}
	}
	return
}

// WithStatusCode sets the HTTP status code for the response
func (r *FlamingoRenderAware) WithStatusCode(code int) RenderAware {
	r.httpStatusCode = code
	return r
}
