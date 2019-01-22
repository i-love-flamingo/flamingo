package session

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type contextKey int

const sessionKey contextKey = iota

// Context saves the session in the context
func Context(ctx context.Context, session *web.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// FromContext retrieves the session from the context
func FromContext(ctx context.Context) (*web.Session, bool) {
	s, ok := ctx.Value(sessionKey).(*web.Session)
	return s, ok
}
