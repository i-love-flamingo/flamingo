package web

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type (
	// Request defines a web request
	Request struct {
		request *http.Request
		vars    map[string]string
		session *sessions.Session
	}
)

var (
	// ErrFormNotFound is triggered for missing form values
	ErrFormNotFound = errors.New("form value not found")
	// ErrParamNotFound for missing router params
	ErrParamNotFound = errors.New("param value not found")
	// ErrQueryNotFound for missing query params
	ErrQueryNotFound = errors.New("query value not found")
)

// RequestFromRequest wraps a http Request
func RequestFromRequest(r *http.Request, session *sessions.Session) *Request {
	return &Request{
		request: r,
		session: session,
	}
}

// WithVars loads parameters
func (r *Request) WithVars(vars map[string]string) *Request {
	request := r.clone()

	request.vars = vars

	return request
}

func (r *Request) clone() *Request {
	return &Request{
		session: r.session,
		request: r.request,
	}
}

// LoadParams load request params
func (r *Request) LoadParams(p map[string]string) *Request {
	r.vars = p
	return r
}

// Session returns the ctx Session
func (r *Request) Session() *sessions.Session {
	return r.session
}

// Form get POST value
func (r *Request) Form(n string) ([]string, bool) {
	f, ok := r.FormAll()[n]
	return f, ok
}

// MustForm panics if n is not found
func (r *Request) MustForm(n string) []string {
	f, ok := r.Form(n)
	if !ok {
		panic(ErrFormNotFound)
	}
	return f
}

// Form1 get first POST value
func (r *Request) Form1(n string) (string, bool) {
	f, ok := r.Form(n)
	if !ok {
		return "", false
	}

	if len(f) > 0 {
		return f[0], true
	}

	return "", false
}

// MustForm1 panics if n is not found
func (r *Request) MustForm1(n string) string {
	f, ok := r.Form1(n)
	if !ok {
		panic(ErrFormNotFound)
	}
	return f
}

// FormAll get all POST values
func (r *Request) FormAll() map[string][]string {
	r.Request().ParseForm()
	return r.Request().Form
}

// Param1 get first querystring param
func (r *Request) Param1(n string) (string, bool) {
	if r, ok := r.vars[n]; ok {
		return r, true
	}
	return "", false
}

// MustParam1 panics if n is not found
func (r *Request) MustParam1(n string) string {
	f, ok := r.Param1(n)
	if !ok {
		panic(ErrParamNotFound)
	}
	return f
}

// ParamAll get all querystring params
func (r *Request) ParamAll() map[string]string {
	return r.vars
}

// Query looks up Raw Query map for Param
func (r *Request) Query(n string) ([]string, bool) {
	f, ok := r.QueryAll()[n]
	return f, ok
}

// MustQuery panics if n is not found
func (r *Request) MustQuery(n string) []string {
	f, ok := r.Query(n)
	if !ok {
		panic(ErrQueryNotFound)
	}
	return f
}

// Query1 looks up Raw Query map for First Param
func (r *Request) Query1(n string) (string, bool) {
	f, ok := r.Query(n)
	if !ok {
		return "", false
	}
	if len(f) > 0 {
		return f[0], true
	}
	return "", false
}

// MustQuery1 panics if n is not found
func (r *Request) MustQuery1(n string) string {
	f, ok := r.Query1(n)
	if !ok {
		panic(ErrQueryNotFound)
	}
	return f
}

// QueryAll returns a Map of the Raw Query
func (r *Request) QueryAll() map[string][]string {
	if r.request == nil {
		return nil
	}
	return r.request.URL.Query()
}

// Request get the requests request
func (r *Request) Request() *http.Request {
	return r.request
}
