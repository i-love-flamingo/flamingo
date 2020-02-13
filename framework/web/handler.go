package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
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

		sessionStore *SessionStore
		sessionName  string
		prefix       string
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

	session, err := h.sessionStore.LoadByRequest(ctx, httpRequest)
	if err != nil {
		h.logger.WithContext(ctx).Warn(err)
	}

	_, span = trace.StartSpan(ctx, "router/matchRequest")
	controller, params, handler := h.routerRegistry.matchRequest(httpRequest)

	if handler != nil {
		ctx, _ = tag.New(ctx, tag.Upsert(ControllerKey, handler.GetHandlerName()), tag.Insert(opencensus.KeyArea, "-"))
		httpRequest = httpRequest.WithContext(ctx)
		start := time.Now()
		defer func() {
			stats.Record(ctx, rt.M(time.Since(start).Nanoseconds()/1000000))
		}()
	}

	req := &Request{
		request: *httpRequest,
		session: Session{s: session.s, sessionSaveMode: session.sessionSaveMode},
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

	ctx, span = trace.StartSpan(ctx, "router/request")
	defer span.End()

	chain := &FilterChain{
		filters: h.filter,
		final: func(ctx context.Context, r *Request, rw http.ResponseWriter) (response Result) {
			ctx, span := trace.StartSpan(ctx, "router/controller")
			defer span.End()

			defer func() {
				if err := panicToError(recover()); err != nil {
					h.logger.WithContext(ctx).Error(err)
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
				err := errors.Errorf("action for method %q not found and no any fallback", req.Request().Method)
				h.logger.WithContext(ctx).Warn(err)
				response = h.routerRegistry.handler[FlamingoNotfound].any(context.WithValue(ctx, RouterError, err), r)
				span.SetStatus(trace.Status{Code: trace.StatusCodeNotFound, Message: "action not found"})
			}

			return response
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
			finishErr = err
			h.logger.WithContext(ctx).Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(rw, "%+v", err)
		}
	}
}
