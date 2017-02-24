package responder

import (
	"flamingo/core/flamingo/web"
	"net/http"
)

// JsonAware allows pug_template rendering
type JsonAware struct{}

// Json returns a web.ContentResponse with status 200
func (r *JsonAware) Json(data interface{}) web.Response {
	return web.JsonResponse{
		Status: http.StatusOK,
		Data:   data,
	}
}
