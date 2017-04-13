package responder

import (
	"flamingo/framework/web"
	"net/http"
)

// JSONAware allows pug_template rendering
type JSONAware struct{}

// JSON returns a web.ContentResponse with status 200
func (r *JSONAware) JSON(data interface{}) web.Response {
	return &web.JSONResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}
