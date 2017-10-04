package responder

import (
	"net/http"

	"go.aoe.com/flamingo/framework/web"
)

type (
	// JSONAware controller trait
	JSONAware interface {
		JSON(data interface{}) web.Response
	}

	// FlamingoJSONAware allows pug_template rendering
	FlamingoJSONAware struct{}
)

var _ JSONAware = &FlamingoJSONAware{}

// JSON returns a web.ContentResponse with status 200
func (r *FlamingoJSONAware) JSON(data interface{}) web.Response {
	return &web.JSONResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}
