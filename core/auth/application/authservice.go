package application

import (
	"context"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//Authinfo - generic Authinfo interface that should be used
	Authinfo interface {
		IsAuthenticated(ctx context.Context,r *web.Request) bool
		GetIdendity(ctx context.Context,r *web.Request) (*domain.Idendity, error)
	}

)


