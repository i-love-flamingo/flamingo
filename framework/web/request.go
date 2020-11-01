package web

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type (
	// Request object stores the actual HTTP Request, Session, Params and attached Values
	Request struct {
		request     http.Request
		session     Session
		handlerName string
		Params      RequestParams
		Values      sync.Map
	}

	// RequestParams store string->string values for request data
	RequestParams map[string]string

	contextKeyType string
)

const contextRequest contextKeyType = "request"

var (
	// ErrFormNotFound is returned for unknown form values
	ErrFormNotFound = errors.New("form value not found")

	// ErrQueryNotFound is returned for unknown URL query parameters
	ErrQueryNotFound = errors.New("query value not found")
)

// CreateRequest creates a new request, with optional http.Request and Session.
// If any variable is nil it is ignored, otherwise it is copied into the new Request.
func CreateRequest(r *http.Request, s *Session) *Request {
	req := new(Request)
	if r != nil {
		req.request = *r
	} else {
		r, _ := http.NewRequest(http.MethodGet, "", nil)
		req.request = *r
	}
	if s != nil {
		req.session.s = s.s
	} else {
		req.session = *EmptySession()
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

// Form get POST value
func (r *Request) Form(name string) ([]string, error) {
	f, err := r.FormAll()
	if err != nil {
		return nil, err
	}
	val, ok := f[name]
	if !ok {
		return nil, ErrFormNotFound
	}
	return val, nil
}

// Form1 get first POST value
func (r *Request) Form1(name string) (string, error) {
	f, err := r.Form(name)
	if err != nil {
		return "", err
	}
	if len(f) > 0 {
		return f[0], nil
	}

	return "", ErrFormNotFound
}

// FormAll get all POST values
func (r *Request) FormAll() (map[string][]string, error) {
	err := r.Request().ParseForm()
	return r.Request().Form, err
}

// Query looks up Raw Query map for Param
func (r *Request) Query(name string) ([]string, error) {
	if f, ok := r.QueryAll()[name]; ok {
		return f, nil
	}

	return nil, ErrQueryNotFound
}

// Query1 looks up Raw Query map for First Param
func (r *Request) Query1(name string) (string, error) {
	if f, err := r.Query(name); err == nil && len(f) > 0 {
		return f[0], nil
	}
	return "", ErrQueryNotFound
}

// QueryAll returns a Map of the Raw Query
func (r *Request) QueryAll() url.Values {
	return r.request.URL.Query()
}

// HandlerName returns a Name of found handler
func (r *Request) HandlerName() string {
	return r.handlerName
}

// HasHandler checks if there is a handler for request
func (r *Request) HasHandler() bool {
	return r.handlerName != ""
}
