package voter

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

// todo: add custom type
const (
	AccessAbstained = iota
	AccessGranted
	AccessDenied
)

// SecurityVoter defines a common interface for voters who vote on security decisions
type SecurityVoter interface {
	Vote(context.Context, *web.Session, string, interface{}) int
}
