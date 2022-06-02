package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/opentelemetry"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gorilla/securecookie"
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
	// ControllerKey exposes the current controller/handler key
	ControllerKey, _ = baggage.NewKeyProperty("controller")

	// RouterError defines error value for issues appearing during routing process
	RouterError contextKeyType = "error"
	rtHistogram syncint64.Histogram
)

func init() {
	var err error
	rtHistogram, err = opentelemetry.GetMeter().SyncInt64().Histogram("flamingo/router/controller",
		instrument.WithDescription("controller request times"), instrument.WithUnit(unit.Milliseconds))
	if err != nil {
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
	ctx, span := opentelemetry.GetTracer().Start(httpRequest.Context(), "router/ServeHTTP")
	defer span.End()

	session, err := h.sessionStore.LoadByRequest(ctx, httpRequest)
	if err != nil {
		if cookieErr, ok := err.(securecookie.Error); ok && (cookieErr.IsUsage() || cookieErr.IsDecode()) {
			h.logger.WithContext(ctx).Debug(err)
		} else {
			h.logger.WithContext(ctx).Warn(err)
		}
	}

	_, span = opentelemetry.GetTracer().Start(ctx, "router/matchRequest")
	controller, params, handler := h.routerRegistry.matchRequest(httpRequest)

	if handler != nil {
		bagg := baggage.FromContext(ctx)
		ctrBaggage, _ := baggage.NewMember(ControllerKey.Key(), handler.GetHandlerName())
		areaBaggage, _ := baggage.NewMember(opentelemetry.KeyArea.String(), "-")
		bagg, _ = bagg.SetMember(ctrBaggage)
		afterDeletionBagg := bagg.DeleteMember(areaBaggage.Key())
		if afterDeletionBagg.Len() == bagg.Len() {
			bagg, _ = bagg.SetMember(areaBaggage)
		}
		ctx = baggage.ContextWithBaggage(ctx, bagg)
		httpRequest = httpRequest.WithContext(ctx)
		start := time.Now()
		defer func() {
			rtHistogram.Record(ctx, time.Since(start).Nanoseconds()/1000000)
		}()
	}

	req := &Request{
		request: *httpRequest,
		session: Session{s: session.s, sessionSaveMode: session.sessionSaveMode},
		Handler: handler,
		Params:  params,
	}
	ctx = ContextWithRequest(ContextWithSession(ctx, req.Session()), req)

	var finishErr error
	defer func() {
		// fire finish event
		h.eventRouter.Dispatch(ctx, &OnFinishEvent{OnRequestEvent{req, rw}, finishErr})
	}()

	h.eventRouter.Dispatch(ctx, &OnRequestEvent{req, rw})

	span.End() // router/matchRequest

	ctx, span = opentelemetry.GetTracer().Start(ctx, "router/request")
	defer span.End()

	chain := &FilterChain{
		filters: h.filter,
		final: func(ctx context.Context, r *Request, rw http.ResponseWriter) (response Result) {
			ctx, span := opentelemetry.GetTracer().Start(ctx, "router/controller")
			defer span.End()

			defer func() {
				if err := panicToError(recover()); err != nil {
					response = h.routerRegistry.handler[FlamingoError].any(context.WithValue(ctx, RouterError, err), r)
					span.SetStatus(codes.Error, "controller panic")
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
				span.SetStatus(codes.Error, "action not found")
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
		ctx, span := opentelemetry.GetTracer().Start(ctx, "router/responseApply")

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
				if errors.Is(err, context.Canceled) {
					h.logger.WithContext(ctx).Debug(err)
				} else {
					h.logger.WithContext(ctx).Error(err)
				}
				finishErr = err
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(rw, "%+v", err)
			}
		}()

		if err := h.routerRegistry.handler[FlamingoError].any(context.WithValue(ctx, RouterError, finalErr), req).Apply(ctx, rw); err != nil {
			panic(err)
		}
	}
}
