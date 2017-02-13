package web

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type (
	// Context c
	Context interface {
		Form(string) []string
		Form1(string) string
		FormAll() map[string][]string
		Param1(string) string
		ParamAll() map[string]string
		Query(string) []string
		Query1(string) string
		QueryAll() map[string][]string
		Request() *http.Request

		ID() string

		Push(target string, opts *http.PushOptions) error
	}

	ctx struct {
		vars    map[string]string
		request *http.Request
		debug   bool
		id      string
		writer  http.ResponseWriter
		pusher  http.Pusher
	}
)

func ContextFromRequest(rw http.ResponseWriter, r *http.Request) *ctx {
	c := new(ctx)
	c.vars = mux.Vars(r)
	c.request = r
	c.id = strconv.Itoa(rand.Int())
	c.writer = rw
	pusher, ok := rw.(http.Pusher)
	if ok {
		c.pusher = pusher
	}

	c.debug = true

	return c
}

func (c *ctx) Push(target string, opts *http.PushOptions) error {
	if c.pusher != nil {
		return c.pusher.Push(target, opts)
	}
	return nil
}

func (c *ctx) ID() string {
	return c.id
}

// Form get POST values
func (c *ctx) Form(n string) []string {
	return c.FormAll()[n]
}

// Form1 get first POST value
func (c *ctx) Form1(n string) string {
	if len(c.Form(n)) < 1 {
		return ""
	}
	return c.Form(n)[0]
}

// FormAll get all POST values
func (c *ctx) FormAll() map[string][]string {
	c.Request().ParseForm()
	return c.Request().Form
}

// Param get querystring param
func (c *ctx) Param1(n string) string {
	return c.vars[n]
}

// Params get all querystring params
func (c *ctx) ParamAll() map[string]string {
	return c.vars
}

func (c *ctx) Query(n string) []string {
	return c.QueryAll()[n]
}

func (c *ctx) Query1(n string) string {
	return c.Query(n)[0]
}

func (c *ctx) QueryAll() map[string][]string {
	return c.request.URL.Query()
}

// Request get the context's request
func (c *ctx) Request() *http.Request {
	return c.request
}
