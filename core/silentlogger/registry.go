package silentlogger

import (
	"context"
	"sync"

	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type LoggingContextRegistry struct {
	mu       sync.RWMutex
	registry map[string]*SilentContext
}

func (r *LoggingContextRegistry) Inject() *LoggingContextRegistry {
	r.registry = make(map[string]*SilentContext)

	return r
}

func (r *LoggingContextRegistry) Get(id string) *SilentContext {
	if r == nil || id == "" {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.registry[id]
}

func (r *LoggingContextRegistry) Notify(_ context.Context, event flamingo.Event) {
	switch typed := event.(type) {
	case *web.OnRequestEvent:
		ctx := typed.Request.Request().Context()
		span := trace.FromContext(ctx)

		r.mu.Lock()
		defer r.mu.Unlock()
		r.registry[span.SpanContext().TraceID.String()] = new(SilentContext)
	case *web.OnFinishEvent:
		ctx := typed.Request.Request().Context()
		span := trace.FromContext(ctx)

		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.registry, span.SpanContext().TraceID.String())
	}
}
