package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/dgrijalva/jwt-go"
	"go.elastic.co/apm/v2"
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
	// RouterError defines error value for issues appearing during routing process
	RouterError contextKeyType = "error"
)

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

	span, ctx := apm.StartSpan(httpRequest.Context(), "router/serveHTTP", "http")
	defer span.End()

	session, err := h.sessionStore.LoadByRequest(ctx, httpRequest)
	if err != nil {
		e := apm.CaptureError(ctx, err)
		e.SetSpan(span)
		e.Send()
		h.logger.WithContext(ctx).Warn(err)
	}

	matchSpan, _ := apm.StartSpan(ctx, "router/matchRequest", "http")
	controller, params, handler := h.routerRegistry.matchRequest(httpRequest)

	var handlerName string
	if handler != nil {
		handlerName = handler.handler
		httpRequest = httpRequest.WithContext(ctx)
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

	matchSpan.End()

	span.Context.SetLabel("handler", handlerName)
	span.Context.SetLabel("username", h.username(httpRequest))
	for k, v := range params {
		span.Context.SetLabel(k, v)
	}
	for k, v := range httpRequest.URL.Query() {
		span.Context.SetLabel(k, v)
	}

	chain := &FilterChain{
		filters: h.filter,
		final: func(ctx context.Context, r *Request, rw http.ResponseWriter) (response Result) {
			childSpan, childCtx := apm.StartSpan(ctx, "router/controller", "http")
			defer childSpan.End()

			defer func() {
				if err := panicToError(recover()); err != nil {
					e := apm.CaptureError(ctx, err)
					e.SetSpan(childSpan)
					e.Send()
					response = h.routerRegistry.handler[FlamingoError].any(context.WithValue(ctx, RouterError, err), r)
					span.Context.SetHTTPStatusCode(http.StatusInternalServerError)
				}
			}()

			defer h.eventRouter.Dispatch(childCtx, &OnResponseEvent{OnRequestEvent{req, rw}, response})

			if c, ok := controller.method[req.Request().Method]; ok && c != nil {
				response = c(childCtx, r)
			} else if controller.any != nil {
				response = controller.any(childCtx, r)
			} else {
				err := fmt.Errorf("action for method %q not found and no \"any\" fallback", req.Request().Method)
				e := apm.CaptureError(ctx, err)
				e.SetSpan(childSpan)
				e.Send()
				response = h.routerRegistry.handler[FlamingoNotfound].any(context.WithValue(childCtx, RouterError, err), r)
				span.Context.SetHTTPStatusCode(http.StatusNotFound)
			}

			return h.responder.completeResult(response)
		},
	}

	result := chain.Next(ctx, req, rw)

	if header, err := h.sessionStore.Save(ctx, req.Session()); err == nil {
		AddHTTPHeader(rw.Header(), header)
	} else {
		e := apm.CaptureError(ctx, err)
		e.SetSpan(span)
		e.Send()
		h.logger.WithContext(ctx).Warn(err)
	}

	var finalErr error
	if result != nil {
		applySpan, newCtx := apm.StartSpan(ctx, "router/applyResponse", "http")

		func() {
			//catch panic in Apply only
			defer func() {
				if err := panicToError(recover()); err != nil {
					e := apm.CaptureError(ctx, err)
					e.SetSpan(applySpan)
					e.Send()
					finalErr = err
				}
			}()
			finalErr = result.Apply(newCtx, rw)
		}()

		applySpan.End()
	}

	// ensure that the session has been saved in the backend
	if _, err := h.sessionStore.Save(ctx, req.Session()); err != nil {
		e := apm.CaptureError(ctx, err)
		e.SetSpan(span)
		e.Send()
		finalErr = err
		h.logger.WithContext(ctx).Warn(err)
	}

	for _, cb := range chain.postApply {
		cb(finalErr, result)
	}

	if finalErr != nil {
		finishErr = finalErr
		defer func() {
			if err := panicToError(recover()); err != nil {
				e := apm.CaptureError(ctx, err)
				e.SetSpan(span)
				e.Send()
				finalErr = err
				finishErr = err
				h.logger.WithContext(ctx).Error(err)
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(rw, "%+v", err)
			}
		}()

		if err := h.routerRegistry.handler[FlamingoError].any(context.WithValue(ctx, RouterError, finalErr), req).Apply(ctx, rw); err != nil {
			e := apm.CaptureError(ctx, err)
			e.SetSpan(span)
			e.Send()
			finalErr = err
			panic(err)
		}
	}
}

func (h *handler) username(httpRequest *http.Request) string {
	parser := &jwt.Parser{}

	auth := httpRequest.Header.Get("Authorization")
	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return ""
	}

	token, _, err := parser.ParseUnverified(authParts[1], jwt.MapClaims{})
	if err != nil {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ""
	}

	return fmt.Sprint(claims["preferred_username"])
}
