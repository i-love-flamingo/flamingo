package web

import (
	"context"

	"go.opencensus.io/trace"
)

// RunWithDetachedContext returns a context which is detached from the original deadlines, timeouts & co
func RunWithDetachedContext(origCtx context.Context, fnc func(ctx context.Context)) {
	request := RequestFromContext(origCtx)
	session := SessionFromContext(origCtx)
	if request != nil && session == nil {
		session = request.Session()
	}
	ctx := ContextWithRequest(context.Background(), request)
	ctx = ContextWithSession(ctx, session)

	ctx, span := trace.StartSpanWithRemoteParent(ctx, "flamingo/detachedContext", trace.FromContext(origCtx).SpanContext())
	defer span.End()

	fnc(ctx)
}
