package responder

import (
	"bytes"
	"encoding/json"
	"flamingo/core/pugtemplate/pugjs"
	"flamingo/framework/router"
	"flamingo/framework/template"
	"flamingo/framework/web"
	"io"
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
func (r *FlamingoRenderAware) Render(context web.Context, tpl string, data interface{}) web.Response {
	if d, err := context.Query1("debugdata"); err == nil && d != "" {
		return &web.JSONResponse{
			Data: pugjs.Convert(data),
		}
	}

	var body io.Reader
	var err error
	if r.Engine != nil {
		body, err = r.Engine.Render(context, tpl, data)
	} else {
		var b []byte
		b, err = json.Marshal(data)
		body = bytes.NewBuffer(b)
	}
	if err != nil {
		panic(err)
	}
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        body,
		ContentType: "text/html; charset=utf-8",
	}
}
