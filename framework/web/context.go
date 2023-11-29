package web

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// RunWithDetachedContext returns a context which is detached from the original deadlines, timeouts & co
func RunWithDetachedContext(origCtx context.Context, fnc func(ctx context.Context)) {
	origCtx, span := otel.Tracer("flamingo.me/opentelemetry").Start(origCtx, "flamingo/detachedContext")
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
