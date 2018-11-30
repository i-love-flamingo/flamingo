package interfaces

import (
	"context"

	"github.com/pkg/errors"

	"flamingo.me/flamingo/core/csrfPreventionFilter/application"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	CsrfMiddleware struct {
		responder *web.Responder
		service   application.Service
	}
)

func (m *CsrfMiddleware) Inject(r *web.Responder, s application.Service) {
	m.responder = r
	m.service = s
}

func (m *CsrfMiddleware) Secured(action router.Action) router.Action {
	return func(ctx context.Context, r *web.Request) web.Response {
		if !m.service.IsValid(r) {
			return m.responder.Forbidden(errors.New("csrf_token is not valid"))
		}

		return action(ctx, r)
	}
}
