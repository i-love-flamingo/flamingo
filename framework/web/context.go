package web

import (
	"context"

	"go.opencensus.io/trace"
)

// RunWithDetachedContext returns a context which is detached from the original deadlines, timeouts & co
func RunWithDetachedContext(origCtx context.Context, fnc func(ctx context.Context)) {
	origCtx, span := trace.StartSpan(origCtx, "flamingo/detachedContext")
	defer span.End()

	request := RequestFromContext(origCtx)
	session := SessionFromContext(origCtx)
	if request != nil && session == nil {
		session = request.Session()
	}

	ctx := ContextWithRequest(trace.NewContext(context.Background(), span), request)
	ctx = ContextWithSession(ctx, session)

	fnc(ctx)
}
