package web

import (
	"context"
	"errors"
	"net/http"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/profiler"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
)

//go:generate mockery -name "Context"

type (
	// ContextKey is used for context.WithValue
	ContextKey string

	// Context defines what a controller sees
	Context interface {
		context.Context
		profiler.Profiler
		Profiler() profiler.Profiler
		EventRouter() event.Router
		WithValue(key, value interface{}) Context

		LoadParams(p map[string]string)
		WithVars(vars map[string]string) Context

		Form(string) ([]string, error)
		MustForm(string) []string
		Form1(string) (string, error)
		MustForm1(string) string
		FormAll() map[string][]string
		Param1(string) (string, error)
		MustParam1(string) string
		ParamAll() map[string]string
		Query(string) ([]string, error)
		MustQuery(string) []string
		Query1(string) (string, error)
		MustQuery1(string) string
		QueryAll() map[string][]string
		Request() *http.Request

		ID() string

		Push(target string, opts *http.PushOptions) error

		Session() *sessions.Session
	}

	// ContextFactory creates a new context
	ContextFactory func(profiler profiler.Profiler, eventrouter event.Router, rw http.ResponseWriter, r *http.Request, session *sessions.Session) Context

	ctx struct {
		context.Context
		profiler    profiler.Profiler
		Eventrouter event.Router

		vars    map[string]string
		request *http.Request
		id      string
		writer  http.ResponseWriter
		pusher  http.Pusher
		session *sessions.Session
	}
)

const (
	// CONTEXT key for contexts
	CONTEXT ContextKey = "context"
	// ID for request id
	ID ContextKey = "ID"
)

// NewContext creates a new context with a unique ID
func NewContext() Context {
	id := uuid.NewV4().String()

	c := &ctx{
		id:       id,
		profiler: new(profiler.NullProfiler),
	}

	c.Context = context.WithValue(context.Background(), ID, c.id)

	return c
}

// ContextFromRequest returns a ctx enriched by Request Data
func ContextFromRequest(profiler profiler.Profiler, eventrouter event.Router, rw http.ResponseWriter, r *http.Request, session *sessions.Session) Context {
	id := uuid.NewV4().String()

	if session != nil {
		if oid, ok := session.Values["context.id"]; ok {
			id = oid.(string)
		}
	}

	c := &ctx{
		request:     r,
		id:          id,
		writer:      rw,
		session:     session,
		profiler:    profiler,
		Eventrouter: eventrouter,
	}

	if pusher, ok := rw.(http.Pusher); ok {
		c.pusher = pusher
	}

	c.Context = context.WithValue(r.Context(), ID, c.id)

	return c
}

// WithVars loads parameters
func (c *ctx) WithVars(vars map[string]string) Context {
	newctx := c.clone()

	newctx.vars = vars

	return newctx
}

func (c *ctx) clone() *ctx {
	var newctx = new(ctx)

	newctx.session = c.session
	newctx.request = c.request
	newctx.profiler = c.profiler
	newctx.Eventrouter = c.Eventrouter
	newctx.id = c.id
	newctx.Context = c.Context

	return newctx
}

// WithValue enriches the context value with a key-value pair
func (c *ctx) WithValue(key, value interface{}) Context {
	c.Context = context.WithValue(c.Context, key, value)
	return c
}

// LoadParams load request params
func (c *ctx) LoadParams(p map[string]string) {
	c.vars = p
}

// EventRouter returns the registered event router
func (c *ctx) EventRouter() event.Router {
	return c.Eventrouter
}

// Profiler returns the registered event router
func (c *ctx) Profiler() profiler.Profiler {
	return c.profiler
}

// Profile profiles
func (c *ctx) Profile(a, b string) profiler.ProfileFinishFunc {
	return c.profiler.Profile(a, b)
}

// Session returns the ctx Session
func (c *ctx) Session() *sessions.Session {
	return c.session
}

// Push pushes Assets via HTTP2 server push
func (c *ctx) Push(target string, opts *http.PushOptions) error {
	if c.pusher != nil {
		return c.pusher.Push(target, opts)
	}
	return nil
}

// ID returns the ctx Id (random Int)
func (c *ctx) ID() string {
	return c.id
}

// Form get POST value
func (c *ctx) Form(n string) ([]string, error) {
	if r, ok := c.FormAll()[n]; ok {
		return r, nil
	}
	return nil, errors.New("form value not found")
}

// MustForm panics if n is not found
func (c *ctx) MustForm(n string) []string {
	r, err := c.Form(n)
	if err != nil {
		panic(err)
	}
	return r
}

// Form1 get first POST value
func (c *ctx) Form1(n string) (string, error) {
	r, err := c.Form(n)
	if err != nil {
		return "", err
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return "", errors.New("form value not found")
}

// MustForm1 panics if n is not found
func (c *ctx) MustForm1(n string) string {
	r, err := c.Form1(n)
	if err != nil {
		panic(err)
	}
	return r
}

// FormAll get all POST values
func (c *ctx) FormAll() map[string][]string {
	c.Request().ParseForm()
	return c.Request().Form
}

// Param1 get first querystring param
func (c *ctx) Param1(n string) (string, error) {
	if r, ok := c.vars[n]; ok {
		return r, nil
	}
	return "", errors.New("param " + n + " not found")
}

// MustParam1 panics if n is not found
func (c *ctx) MustParam1(n string) string {
	r, err := c.Param1(n)
	if err != nil {
		panic(err)
	}
	return r
}

// ParamAll get all querystring params
func (c *ctx) ParamAll() map[string]string {
	return c.vars
}

// Query looks up Raw Query map for Param
func (c *ctx) Query(n string) ([]string, error) {
	if r, ok := c.QueryAll()[n]; ok {
		return r, nil
	}
	return nil, errors.New("query values not found")
}

// MustQuery panics if n is not found
func (c *ctx) MustQuery(n string) []string {
	r, err := c.Query(n)
	if err != nil {
		panic(err)
	}
	return r
}

// Query1 looks up Raw Query map for First Param
func (c *ctx) Query1(n string) (string, error) {
	r, err := c.Query(n)
	if err != nil {
		return "", err
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return "", errors.New("query value not found")
}

// MustQuery1 panics if n is not found
func (c *ctx) MustQuery1(n string) string {
	r, err := c.Query1(n)
	if err != nil {
		panic(err)
	}
	return r
}

// QueryAll returns a Map of the Raw Query
func (c *ctx) QueryAll() map[string][]string {
	if c.request == nil {
		return nil
	}
	return c.request.URL.Query()
}

// Request get the context's request
func (c *ctx) Request() *http.Request {
	return c.request
}
