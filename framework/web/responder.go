package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"strings"

	"flamingo.me/flamingo/framework/template"
)

type (
	// ReverseRouter for RouteRedirect Responses
	ReverseRouter interface {
		URL(name string, params map[string]string) *url.URL
	}

	// Responder generates responses
	Responder struct {
		engine template.Engine
		router ReverseRouter
	}

	// HTTPResponse contains a status and a body
	HTTPResponse struct {
		Status uint
		Body   io.Reader
		Header http.Header
		hooks  *[]ResponseHook
	}

	// RouteRedirectResponse redirects to a certain route
	RouteRedirectResponse struct {
		HTTPResponse
		To     string
		Data   map[string]string
		router ReverseRouter
	}

	// URLRedirectResponse redirects to a certain URL
	URLRedirectResponse struct {
		HTTPResponse
		URL *url.URL
	}

	// DataResponse returns a response containing data, e.g. as JSON
	DataResponse struct {
		HTTPResponse
		Data interface{}
	}

	// RenderResponse renders data
	RenderResponse struct {
		DataResponse
		Template string
		engine   template.Engine
	}

	// ServerErrorResponse returns a server error, by default http 500
	ServerErrorResponse struct {
		HTTPResponse
		Error error
	}
)

// Inject Responder dependencies
func (r *Responder) Inject(engine template.Engine, router ReverseRouter) *Responder {
	r.engine = engine
	r.router = router
	return r
}

var _ Response = HTTPResponse{}

// Apply response
func (r HTTPResponse) Apply(c context.Context, w http.ResponseWriter) error {
	w.WriteHeader(int(r.Status))
	for name, vals := range r.Header {
		for _, val := range vals {
			w.Header().Add(name, val)
		}
	}
	_, err := io.Copy(w, r.Body)
	return err
}

// GetStatus returns the HTTP status
func (r HTTPResponse) GetStatus() int {
	return int(r.Status)
}

// GetContentLength returns the bodies content length
func (HTTPResponse) GetContentLength() int {
	return 0
}

// Hook helper
// deprecated: to be removed
func (r HTTPResponse) Hook(hooks ...ResponseHook) Response {
	*r.hooks = append(*r.hooks, hooks...)
	return r
}

// HTTP Response generator
func (r *Responder) HTTP(status uint, body io.Reader) HTTPResponse {
	return HTTPResponse{
		Status: status,
		Body:   body,
	}
}

// RouteRedirect generator
func (r *Responder) RouteRedirect(to string, data map[string]string) RouteRedirectResponse {
	return RouteRedirectResponse{
		To:     to,
		Data:   data,
		router: r.router,
		HTTPResponse: HTTPResponse{
			Status: http.StatusSeeOther,
		},
	}
}

// Apply response
func (r RouteRedirectResponse) Apply(c context.Context, w http.ResponseWriter) error {
	to := r.router.URL(r.To, r.Data)
	w.Header().Set("Location", to.String())
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r RouteRedirectResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// URLRedirect returns a 303 redirect to a given URL
func (r *Responder) URLRedirect(url *url.URL) URLRedirectResponse {
	return URLRedirectResponse{
		URL: url,
		HTTPResponse: HTTPResponse{
			Status: http.StatusSeeOther,
		},
	}
}

// Apply response
func (r URLRedirectResponse) Apply(c context.Context, w http.ResponseWriter) error {
	w.Header().Set("Location", r.URL.String())
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r URLRedirectResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// Data returns a data response which can be serialized
func (r *Responder) Data(data interface{}) DataResponse {
	return DataResponse{
		Data: data,
		HTTPResponse: HTTPResponse{
			Status: http.StatusOK,
		},
	}
}

// Apply response
func (r DataResponse) Apply(c context.Context, w http.ResponseWriter) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(r.Data); err != nil {
		return err
	}
	r.HTTPResponse.Header.Set("Content-Type", "application/json")
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r DataResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// Render creates a render response, with the supplied template and data
func (r *Responder) Render(tpl string, data interface{}) RenderResponse {
	return RenderResponse{
		Template:     tpl,
		engine:       r.engine,
		DataResponse: r.Data(data),
	}
}

// Apply response
func (r RenderResponse) Apply(c context.Context, w http.ResponseWriter) error {
	var err error

	if req, ok := FromContext(c); ok && r.engine != nil {
		partialRenderer, ok := r.engine.(template.PartialEngine)
		if partials := req.Request().Header.Get("X-Partial"); partials != "" && ok {
			content, err := partialRenderer.RenderPartials(c, r.Template, r.Data, strings.Split(partials, ","))
			body, err := json.Marshal(content)
			if err != nil {
				return err
			}
			r.Body = bytes.NewBuffer(body)
			r.Header.Set("Content-Type", "application/json; charset=utf-8")
			return r.HTTPResponse.Apply(c, w)
		}
	}

	r.Body, err = r.engine.Render(c, r.Template, r.Data)
	if err != nil {
		return err
	}
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r RenderResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// ServerError creates a 500 error response
func (r *Responder) ServerError(err error) ServerErrorResponse {
	return ServerErrorResponse{
		Error: err,
		HTTPResponse: HTTPResponse{
			Status: http.StatusInternalServerError,
		},
	}
}

// Apply response
func (r ServerErrorResponse) Apply(c context.Context, w http.ResponseWriter) error {
	r.Body = bytes.NewBufferString(r.Error.Error())
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r ServerErrorResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// NotFound creates a 404 error response
func (r *Responder) NotFound(err error) ServerErrorResponse {
	return ServerErrorResponse{
		Error: err,
		HTTPResponse: HTTPResponse{
			Status: http.StatusNotFound,
		},
	}
}

// TODO creates a 501 Not Implemented response
func (r *Responder) TODO() HTTPResponse {
	return HTTPResponse{
		Status: http.StatusNotImplemented,
	}
}
