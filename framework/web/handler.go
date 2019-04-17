package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
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

		sessionStore sessions.Store
		sessionName  string
		prefix       string
	}

	emptyResponseWriter struct{}
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

func (h *handler) getSession(ctx context.Context, httpRequest *http.Request) (gs *sessions.Session) {
	// initialize the session
	if h.sessionStore != nil {
		var span *trace.Span
		var err error

		ctx, span = trace.StartSpan(ctx, "router/sessions/get")
		gs, err = h.sessionStore.Get(httpRequest, h.sessionName)
		if err != nil {
			h.logger.WithContext(ctx).Warn(err)
			_, span := trace.StartSpan(ctx, "router/sessions/new")
			gs, err = h.sessionStore.New(httpRequest, h.sessionName)
			if err != nil {
				h.logger.WithContext(ctx).Warn(err)
			}
			span.End()
		}
		span.End()
	}

	return
}

func panicToError(p interface{}) error {
	if p == nil {
		return nil
	}

	var err error
	switch errIface := p.(type) {
	case error:
		err = errors.WithStack(errIface)
	case string:
		err = errors.New(errIface)
	default:
		err = errors.Errorf("router/controller: %+v", errIface)
	}
	return err
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, httpRequest *http.Request) {
	httpRequest.URL.Path = strings.TrimPrefix(httpRequest.URL.Path, h.prefix)

	ctx, span := trace.StartSpan(httpRequest.Context(), "router/ServeHTTP")
	defer span.End()

	gs := h.getSession(ctx, httpRequest)

	_, span = trace.StartSpan(ctx, "router/matchRequest")
	controller, params, handler := h.routerRegistry.matchRequest(httpRequest)

	if handler != nil {
		ctx, _ = tag.New(ctx, tag.Upsert(ControllerKey, handler.GetHandlerName()), tag.Upsert(opencensus.KeyArea, "-"))
		httpRequest = httpRequest.WithContext(ctx)
		start := time.Now()
		defer func() {
			stats.Record(ctx, rt.M(time.Since(start).Nanoseconds()/1000000))
		}()
	}

	req := &Request{
		request: *httpRequest,
		session: Session{
			s: gs,
		},
		Params: params,
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
				response = h.routerRegistry.handler[FlamingoNotfound].any(context.WithValue(ctx, RouterError, errors.Errorf("action for method %q not found and no any fallback", req.Request().Method)), r)
				span.SetStatus(trace.Status{Code: trace.StatusCodeNotFound, Message: "action not found"})
			}

			return response
		},
	}

	result := chain.Next(ctx, req, rw)

	if h.sessionStore != nil {
		ctx, span := trace.StartSpan(ctx, "router/sessions/save")
		if err := h.sessionStore.Save(req.Request(), rw, gs); err != nil {
			h.logger.WithContext(ctx).Warn(err)
		}
		span.End()
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
	if h.sessionStore != nil {
		ctx, span := trace.StartSpan(ctx, "router/sessions/persist")
		if err := h.sessionStore.Save(req.Request(), emptyResponseWriter{}, gs); err != nil {
			h.logger.WithContext(ctx).Warn(err)
		}
		span.End()
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
			finishErr = err
			h.logger.WithContext(ctx).Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(rw, "%+v", err)
		}
	}
}

// emptyResponseWriter to be able to properly persist sessions
func (emptyResponseWriter) Header() http.Header        { return http.Header{} }
func (emptyResponseWriter) Write([]byte) (int, error)  { return 0, io.ErrUnexpectedEOF }
func (emptyResponseWriter) WriteHeader(statusCode int) {}
