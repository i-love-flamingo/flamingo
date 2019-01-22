package voter

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

const (
	AccessAbstained = iota
	AccessGranted
	AccessDenied
)

type (
	SecurityVoter interface {
		Vote(context.Context, *web.Session, string, interface{}) int
	}
)
