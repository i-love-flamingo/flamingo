package session

import (
	"context"

	"github.com/gorilla/sessions"
)

type contextKey int

const sessionKey contextKey = iota

// Context saves the session in the context
func Context(ctx context.Context, session *sessions.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// FromContext retrieves the session from the context
func FromContext(ctx context.Context) (*sessions.Session, bool) {
	s, ok := ctx.Value(sessionKey).(*sessions.Session)
	return s, ok
}
