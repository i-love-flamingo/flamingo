package interfaces

import (
	"context"
	"net/url"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//Authservice - generic Authservice interface that should be used
	Authservice interface {
		Authenticate(ctx context.Context, returnURL *url.URL) (web.Result, error)
		IsAuthenticated(ctx context.Context,r *web.Request) bool
		GetIdendity(ctx context.Context,r *web.Request) (domain.Idendity, error)
	}
)


