package responder

import (
	"flamingo/core/web"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	RouterAware interface {
		SetRouter(r *mux.Router)
	}
)

type RedirectFor struct {
	RouterAware
	router *mux.Router
}

func (r *RedirectFor) SetRouter(router *mux.Router) {
	r.router = router
}

func (r *RedirectFor) Response(name string, args ...string) web.Response {
	url, _ := r.router.Get(name).URL(args...)

	return web.RedirectResponse{
		Status:   http.StatusFound,
		Location: url.String(),
	}
}
