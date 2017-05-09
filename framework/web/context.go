package web

import (
	"context"
	"encoding/json"
	"errors"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type (
	// ContextKey is used for context.WithValue
	ContextKey string

	// Context defines what a controller sees
	Context interface {
		context.Context
		profiler.Profiler
		Profiler() profiler.Profiler
		EventRouter() event.Router

		LoadVars(r *http.Request)
		Form(string) ([]string, error)
		Form1(string) (string, error)
		FormAll() map[string][]string
		Param1(string) (string, error)
		ParamAll() map[string]string
		Query(string) ([]string, error)
		Query1(string) (string, error)
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
	CONTEXT ContextKey = "context"
	ID                 = "ID"
)

// ContextFromRequest returns a ctx enriched by Request Data
func ContextFromRequest(profiler profiler.Profiler, eventrouter event.Router, rw http.ResponseWriter, r *http.Request, session *sessions.Session) Context {
	c := &ctx{
		request:     r,
		id:          strconv.Itoa(rand.Int()),
		writer:      rw,
		session:     session,
		profiler:    profiler,
		Eventrouter: eventrouter,
	}

	if pusher, ok := rw.(http.Pusher); ok {
		c.pusher = pusher
	}

	c.Context = context.WithValue(r.Context(), ID, c.id)

	if r.Header.Get("Content-Type") == "application/json" {
		b, _ := ioutil.ReadAll(r.Body)
		var data map[string]string
		json.Unmarshal(b, &data)

		r.PostForm = make(url.Values)
		for k, v := range data {
			r.PostForm[k] = []string{v}
		}
	}

	return c
}

// LoadVars load request vars from mux.Vars
func (c *ctx) LoadVars(r *http.Request) {
	c.vars = mux.Vars(r)
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

// Form1 get first POST value
func (c *ctx) Form1(n string) (string, error) {
	r, err := c.Form(n)
	if err {
		return "", err
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return "", errors.New("form value not found")
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
	return nil, errors.New("param value not found")
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

// Query1 looks up Raw Query map for First Param
func (c *ctx) Query1(n string) (string, error) {
	r, err := c.Query(n)
	if err {
		return "", err
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, errors.New("query value not found")
}

// QueryAll returns a Map of the Raw Query
func (c *ctx) QueryAll() map[string][]string {
	return c.request.URL.Query()
}

// Request get the context's request
func (c *ctx) Request() *http.Request {
	return c.request
}
