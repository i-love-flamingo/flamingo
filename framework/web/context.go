package web

import (
	"context"

	"flamingo.me/flamingo/v3/framework/opentelemetry"

	"go.opentelemetry.io/otel/trace"
)

// RunWithDetachedContext returns a context which is detached from the original deadlines, timeouts & co
func RunWithDetachedContext(origCtx context.Context, fnc func(ctx context.Context)) {
	origCtx, span := opentelemetry.GetTracer().Start(origCtx, "flamingo/detachedContext")
	defer span.End()

	request := RequestFromContext(origCtx)
	session := SessionFromContext(origCtx)
	if request != nil && session == nil {
		session = request.Session()
	}

	ctx := ContextWithRequest(trace.ContextWithSpan(context.Background(), span), request)
	ctx = ContextWithSession(ctx, session)

	fnc(ctx)
}
