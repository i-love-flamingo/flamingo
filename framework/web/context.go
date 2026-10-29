package web

import (
	"context"

	"go.elastic.co/apm/v2"
)

// RunWithDetachedContext returns a context which is detached from the original deadlines, timeouts & co
func RunWithDetachedContext(origCtx context.Context, fnc func(ctx context.Context)) {
	span, origCtx := apm.StartSpan(origCtx, "flamingo/detachedContext", "custom")
	defer span.End()

	request := RequestFromContext(origCtx)
	session := SessionFromContext(origCtx)
	if request != nil && session == nil {
		session = request.Session()
	}

	ctx := ContextWithRequest(apm.ContextWithSpan(origCtx, span), request)
	ctx = ContextWithSession(ctx, session)

	fnc(ctx)
}
