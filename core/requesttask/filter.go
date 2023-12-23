package requesttask

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	filter struct{}
	rKey   string
)

const wg rKey = "requestTaskWg"

// Do runs a background task in the current request scope
func Do(ctx context.Context, r *web.Request, task func(ctx context.Context, r *web.Request)) {
	if err := TryDo(ctx, r, task); err != nil {
		task(ctx, r)
	}
}

// TryDo tries to schedule an async task in the background
func TryDo(ctx context.Context, r *web.Request, task func(ctx context.Context, r *web.Request)) error {
	wg, _ := r.Values.Load(wg)
	if wg, ok := wg.(*sync.WaitGroup); ok {
		wg.Add(1)

		go func() {
			ctx, span := otel.Tracer("flamingo.me/opentelemetry").Start(ctx, "requestTask")
			task(ctx, r)
			span.End()
			wg.Done()
		}()

		return nil
	}

	return errors.New("the current request is unable to schedule background tasks")
}

// Filter waits for running tasks to finish before the request processing is done
func (f *filter) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, fc *web.FilterChain) web.Result {
	r.Values.Store(wg, new(sync.WaitGroup))
	response := fc.Next(ctx, r, w)

	// wait for possible tasks to finish
	wg, _ := r.Values.Load(wg)
	if wg, ok := wg.(*sync.WaitGroup); ok {
		wg.Wait()
	}

	return response
}
