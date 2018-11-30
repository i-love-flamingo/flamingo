package interfaces

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"flamingo.me/flamingo/core/csrfPreventionFilter/application"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	CsrfFilter struct {
		responder *web.Responder
		service   application.Service
	}
)

func (f *CsrfFilter) Inject(r *web.Responder, s application.Service) {
	f.responder = r
	f.service = s
}

func (f *CsrfFilter) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	if !f.service.IsValid(r) {
		return f.responder.Forbidden(errors.New("csrf_token is not valid"))
	}

	return chain.Next(ctx, r, w)
}
