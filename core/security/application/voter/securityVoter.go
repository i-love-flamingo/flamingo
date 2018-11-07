package voter

import (
	"context"

	"github.com/gorilla/sessions"
)

const (
	AccessAbstained = iota
	AccessGranted
	AccessDenied
)

type (
	SecurityVoter interface {
		Vote(context.Context, *sessions.Session, string, interface{}) int
	}
)
