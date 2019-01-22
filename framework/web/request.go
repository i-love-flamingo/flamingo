package web

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

type (
	// Request object stores the actual HTTP Request, Session, Params and attached Values
	Request struct {
		request http.Request
		session Session
		Params  RequestParams
		Values  sync.Map
	}

	// RequestParams store string->string values for request data
	RequestParams map[string]string

	contextKeyType string
)

const (
	contextRequest contextKeyType = "request"
)

// CreateRequest creates a new request, with optional http.Request and Session.
// If any variable is nil it is ignored, otherwise it is copied into the new Request.
func CreateRequest(r *http.Request, s *Session) *Request {
	req := new(Request)
	if r != nil {
		req.request = *r
	}
	if s != nil {
		req.session.s = s.s
	}
	req.Params = make(RequestParams)
	return req
}

// RequestFromContext retrieves the request from the context, if available
func RequestFromContext(ctx context.Context) *Request {
	req, _ := ctx.Value(contextRequest).(*Request)
	return req
}

// ContextWithRequest stores the request in a new context
func ContextWithRequest(ctx context.Context, r *Request) context.Context {
	return context.WithValue(ctx, contextRequest, r)
}

// Request getter
func (r *Request) Request() *http.Request {
	return &r.request
}

// Session getter
func (r *Request) Session() *Session {
	return &r.session
}

// RemoteAddress get the requests real remote address
func (r *Request) RemoteAddress() []string {
	var remoteAddress []string

	forwardFor := strings.TrimSpace(r.request.Header.Get("X-Forwarded-For"))
	ips := strings.Split(forwardFor, ",")
	if len(forwardFor) > 0 {
		for _, ip := range ips {
			remoteAddress = append(remoteAddress, strings.TrimSpace(ip))
		}

	}

	remoteAddress = append(remoteAddress, strings.TrimSpace(r.request.RemoteAddr))

	return remoteAddress
}
