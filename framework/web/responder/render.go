package responder

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"strings"

	"flamingo.me/flamingo/v3/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/v3/framework/router"
	"flamingo.me/flamingo/v3/framework/template"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// RenderAware controller trait
	RenderAware interface {
		Render(context context.Context, tpl string, data interface{}) web.Response
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
func (r *FlamingoRenderAware) Render(context context.Context, tpl string, data interface{}) (response web.Response) {
	statusCode := http.StatusOK
	if r.httpStatusCode > 0 {
		statusCode = r.httpStatusCode
	}

	//if d, err := context.Query1("debugdata"); err == nil && d != "" {
	//	return &web.JSONResponse{
	//		Data: pugjs.Convert(data),
	//	}
	//}

	if req, ok := web.FromContext(context); ok && r.Engine != nil {
		partialRenderer, ok := r.Engine.(template.PartialEngine)
		if partials := req.Request().Header.Get("X-Partial"); partials != "" && ok {
			content, err := partialRenderer.RenderPartials(context, tpl, data, strings.Split(partials, ","))
			body, err := json.Marshal(map[string]interface{}{"partials": content, "data": new(web.GetPartialDataFunc).Func(context).(func() map[string]interface{})()})
			if err != nil {
				panic(err)
			}
			return &web.ContentResponse{
				BasicResponse: web.BasicResponse{
					Status: statusCode,
				},
				Body:        bytes.NewReader(body),
				ContentType: "application/json; charset=utf-8",
			}
		}
	}

	if r.Engine != nil {
		body, err := r.Engine.Render(context, tpl, data)
		if err != nil {
			panic(err)
		}
		response = &web.ContentResponse{
			BasicResponse: web.BasicResponse{
				Status: statusCode,
			},
			Body:        body,
			ContentType: "text/html; charset=utf-8",
		}
	} else {
		body, err := json.Marshal(pugjs.Convert(data))
		if err != nil {
			panic(err)
		}
		response = &web.ContentResponse{
			BasicResponse: web.BasicResponse{
				Status: statusCode,
			},
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
