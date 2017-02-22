package responder

import (
	"flamingo/core/flamingo/web"
	"net/http"
)

// RenderAware allows pug_template rendering
type JsonAware struct{}

// Render returns a web.ContentResponse with status 200 and ContentType text/html
func (r *JsonAware) Json(data interface{}) web.Response {
	return web.JsonResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}
