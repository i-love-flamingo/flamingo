package responder

import (
	"net/http"

	"flamingo.me/flamingo/framework/web"
)

type (
	// JSONAware controller trait
	JSONAware interface {
		JSON(data interface{}) web.Response
		JSONError(data interface{}, statusCode int) web.Response
	}

	// FlamingoJSONAware allows pug_template rendering
	FlamingoJSONAware struct{}
)

var _ JSONAware = &FlamingoJSONAware{}

// JSON returns a web.ContentResponse with status 200
func (r *FlamingoJSONAware) JSON(data interface{}) web.Response {
	return &web.JSONResponse{
		BasicResponse: web.BasicResponse{
			Status: http.StatusOK,
		},
		Data: data,
	}
}

// JSONError returns a web.ContentResponse with status given
func (r *FlamingoJSONAware) JSONError(data interface{}, statusCode int) web.Response {
	return &web.JSONResponse{
		BasicResponse: web.BasicResponse{
			Status: statusCode,
		},
		Data: data,
	}
}
