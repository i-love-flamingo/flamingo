package responder

import (
	"flamingo/core/web"
	"net/http"

	"github.com/gorilla/mux"
)

type RenderTemplate struct {
	router *mux.Router
}

func (r *RenderTemplate) SetRouter(router *mux.Router) {
	r.router = router
}

func (r *RenderTemplate) RenderResponse(context web.Context, tpl string) web.Response {
	return web.ContentResponse{
		Status:      http.StatusOK,
		Body:        web.Render(context, r.router, tpl, nil),
		ContentType: "text/html",
	}
}
