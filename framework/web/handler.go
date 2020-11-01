package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type (
	handler struct {
		routerRegistry *RouterRegistry
		filter         []Filter

		eventRouter flamingo.EventRouter
		logger      flamingo.Logger

		sessionStore *SessionStore
		sessionName  string
		prefix       string
		responder    *Responder
	}

	panicError struct {
		err   error
		stack []byte
	}
)

var (
	rt = stats.Int64("flamingo/router/controller", "controller request times", stats.UnitMilliseconds)
	// ControllerKey exposes the current controller/handler key
	ControllerKey, _ = tag.NewKey("controller")

	// RouterError defines error value for issues appearing during routing process
	RouterError contextKeyType = "error"
)

func init() {
	if err := opencensus.View("flamingo/router/controller", rt, view.Distribution(100, 500, 1000, 2500, 5000, 10000), ControllerKey); err != nil {
		panic(err)
	}
}

func (e *panicError) Error() string {
	return e.err.Error()
}

func (e *panicError) String() string {
	return e.err.Error()
}

func (e *panicError) Unwrap() error {
	return e.err
}

func (e *panicError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, e.err.Error())
			_, _ = fmt.Fprintf(s, "\n%s", e.stack)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.err.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.err)
	}
}

func panicToError(p interface{}) error {
	if p == nil {
		return nil
	}

	var err error
	switch errIface := p.(type) {
	case error:
		//err = fmt.Errorf("controller panic: %w", errIface)
		err = &panicError{err: fmt.Errorf("controller panic: %w", errIface), stack: debug.Stack()}
	case string:
		err = &panicError{err: fmt.Errorf("controller panic: %s", errIface), stack: debug.Stack()}
	default:
		err = &panicError{err: fmt.Errorf("controller panic: %+v", errIface), stack: debug.Stack()}
	}
	return err
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, httpRequest *http.Request) {
	httpRequest.URL.Path = strings.TrimPrefix(httpRequest.URL.Path, h.prefix)

	ctx, span := trace.StartSpan(httpRequest.Context(), "router/ServeHTTP")
	defer span.End()

	session, err := h.sessionStore.LoadByRequest(ctx, httpRequest)
	if err != nil {
		h.logger.WithContext(ctx).Warn(err)
	}

	_, span = trace.StartSpan(ctx, "router/matchRequest")
	controller, params, handler := h.routerRegistry.matchRequest(httpRequest)

	var handlerName string
	if handler != nil {
		handlerName = handler.handler
		ctx, _ = tag.New(ctx, tag.Upsert(ControllerKey, handler.GetHandlerName()), tag.Insert(opencensus.KeyArea, "-"))
		httpRequest = httpRequest.WithContext(ctx)
		start := time.Now()
		defer func() {
			stats.Record(ctx, rt.M(time.Since(start).Nanoseconds()/1000000))
		}()
	}

	req := &Request{
		request:     *httpRequest,
		session:     Session{s: session.s, sessionSaveMode: session.sessionSaveMode},
		handlerName: handlerName,
		Params:      params,
	}
	ctx = ContextWithRequest(ContextWithSession(ctx, req.Session()), req)

	var finishErr error
	defer func() {
		// fire finish event
		h.eventRouter.Dispatch(ctx, &OnFinishEvent{OnRequestEvent{req, rw}, finishErr})
	}()

	h.eventRouter.Dispatch(ctx, &OnRequestEvent{req, rw})

	span.End() // router/matchRequest

	ctx, span = trace.StartSpan(ctx, "router/request")
	defer span.End()

	chain := &FilterChain{
		filters: h.filter,
		final: func(ctx context.Context, r *Request, rw http.ResponseWriter) (response Result) {
			ctx, span := trace.StartSpan(ctx, "router/controller")
			defer span.End()

			defer func() {
				if err := panicToError(recover()); err != nil {
					response = h.routerRegistry.handler[FlamingoError].any(context.WithValue(ctx, RouterError, err), r)
					span.SetStatus(trace.Status{Code: trace.StatusCodeAborted, Message: "controller panic"})
				}
			}()

			defer h.eventRouter.Dispatch(ctx, &OnResponseEvent{OnRequestEvent{req, rw}, response})

			if c, ok := controller.method[req.Request().Method]; ok && c != nil {
				response = c(ctx, r)
			} else if controller.any != nil {
				response = controller.any(ctx, r)
			} else {
				err := fmt.Errorf("action for method %q not found and no \"any\" fallback", req.Request().Method)
				response = h.routerRegistry.handler[FlamingoNotfound].any(context.WithValue(ctx, RouterError, err), r)
				span.SetStatus(trace.Status{Code: trace.StatusCodeNotFound, Message: "action not found"})
			}

			return h.responder.completeResult(response)
		},
	}

	result := chain.Next(ctx, req, rw)

	if header, err := h.sessionStore.Save(ctx, req.Session()); err == nil {
		AddHTTPHeader(rw.Header(), header)
	} else {
		h.logger.WithContext(ctx).Warn(err)
	}

	var finalErr error
	if result != nil {
		ctx, span := trace.StartSpan(ctx, "router/responseApply")

		func() {
			//catch panic in Apply only
			defer func() {
				if err := panicToError(recover()); err != nil {
					finalErr = err
				}
			}()
			finalErr = result.Apply(ctx, rw)
		}()

		span.End()
	}

	// ensure that the session has been saved in the backend
	if _, err := h.sessionStore.Save(ctx, req.Session()); err != nil {
		h.logger.WithContext(ctx).Warn(err)
	}

	for _, cb := range chain.postApply {
		cb(finalErr, result)
	}

	if finalErr != nil {
		finishErr = finalErr
		defer func() {
			if err := panicToError(recover()); err != nil {
				finishErr = err
				h.logger.WithContext(ctx).Error(err)
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(rw, "%+v", err)
			}
		}()

		if err := h.routerRegistry.handler[FlamingoError].any(context.WithValue(ctx, RouterError, finalErr), req).Apply(ctx, rw); err != nil {
			panic(err)
		}
	}
}
