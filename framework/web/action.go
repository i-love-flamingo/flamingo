package web

import (
	"context"
	"net/http"
)

type (
	// Action defines an explicit http action
	Action func(ctx context.Context, req *Request) Result

	// DataAction is a method called which does not return the web response itself, but data instead
	DataAction func(ctx context.Context, req *Request, callParams RequestParams) interface{}

	wrappedHTTPHandler struct {
		handler http.Handler
		request *Request
	}
)

// WrapHTTPHandler wraps an http.Handler to be used in the flamingo http package
func WrapHTTPHandler(handler http.Handler) Action {
	return func(ctx context.Context, req *Request) Result {
		return &wrappedHTTPHandler{
			handler: handler,
			request: req,
		}
	}
}

// WrapDataAction allows to register a data action for a HTTP method
func WrapDataAction(da DataAction) Action {
	return func(ctx context.Context, req *Request) Result {
		return &DataResponse{
			Data: da(ctx, req, req.Params),
		}
	}
}

func (h *wrappedHTTPHandler) Apply(ctx context.Context, rw http.ResponseWriter) error {
	h.handler.ServeHTTP(rw, h.request.Request().WithContext(ctx))
	return nil
}

func (a Action) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	result := a(r.Context(), &Request{request: *r, session: *SessionFromContext(r.Context())})
	if err := result.Apply(r.Context(), rw); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
