package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"flamingo.me/flamingo/v3/framework/template"
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

		templateForbidden     string
		templateNotFound      string
		templateUnavailable   string
		templateErrorWithCode string
	}

	// HTTPResponse contains a status and a body
	HTTPResponse struct {
		Status uint
		Body   io.Reader
		Header http.Header
		hooks  []ResponseHook
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
		RenderResponse
		Error error
	}
)

// Inject Responder dependencies
func (r *Responder) Inject(engine template.Engine, router ReverseRouter, cfg *struct {
	TemplateForbidden     string `inject:"config:flamingo.template.err403"`
	TemplateNotFound      string `inject:"config:flamingo.template.err404"`
	TemplateUnavailable   string `inject:"config:flamingo.template.err503"`
	TemplateErrorWithCode string `inject:"config:flamingo.template.errWithCode"`
}) *Responder {
	r.engine = engine
	r.router = router
	r.templateForbidden = cfg.TemplateForbidden
	r.templateNotFound = cfg.TemplateNotFound
	r.templateUnavailable = cfg.TemplateUnavailable
	r.templateErrorWithCode = cfg.TemplateErrorWithCode

	return r
}

var _ Response = &HTTPResponse{}

// Apply response
func (r *HTTPResponse) Apply(c context.Context, w http.ResponseWriter) error {
	for _, hook := range r.hooks {
		hook(c, w)
	}
	for name, vals := range r.Header {
		for _, val := range vals {
			w.Header().Add(name, val)
		}
	}

	w.WriteHeader(int(r.Status))
	if r.Body == nil {
		return nil
	}

	_, err := io.Copy(w, r.Body)
	return err
}

// GetStatus returns the HTTP status
func (r *HTTPResponse) GetStatus() int {
	return int(r.Status)
}

// GetContentLength returns the bodies content length
func (HTTPResponse) GetContentLength() int {
	return 0
}

// Hook helper
// deprecated: to be removed
func (r *HTTPResponse) Hook(hooks ...ResponseHook) Response {
	r.hooks = append(r.hooks, hooks...)
	return r
}

// HTTP Response generator
func (r *Responder) HTTP(status uint, body io.Reader) *HTTPResponse {
	return &HTTPResponse{
		Status: status,
		Body:   body,
		Header: make(http.Header),
	}
}

// RouteRedirect generator
func (r *Responder) RouteRedirect(to string, data map[string]string) *RouteRedirectResponse {
	return &RouteRedirectResponse{
		To:     to,
		Data:   data,
		router: r.router,
		HTTPResponse: HTTPResponse{
			Status: http.StatusSeeOther,
			Header: make(http.Header),
		},
	}
}

// Apply response
func (r *RouteRedirectResponse) Apply(c context.Context, w http.ResponseWriter) error {
	to := r.router.URL(r.To, r.Data)
	w.Header().Set("Location", to.String())
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r *RouteRedirectResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// URLRedirect returns a 303 redirect to a given URL
func (r *Responder) URLRedirect(url *url.URL) *URLRedirectResponse {
	return &URLRedirectResponse{
		URL: url,
		HTTPResponse: HTTPResponse{
			Status: http.StatusSeeOther,
			Header: make(http.Header),
		},
	}
}

// Apply response
func (r *URLRedirectResponse) Apply(c context.Context, w http.ResponseWriter) error {
	w.Header().Set("Location", r.URL.String())
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r *URLRedirectResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// Data returns a data response which can be serialized
func (r *Responder) Data(data interface{}) *DataResponse {
	return &DataResponse{
		Data: data,
		HTTPResponse: HTTPResponse{
			Status: http.StatusOK,
			Header: make(http.Header),
		},
	}
}

// Apply response
func (r *DataResponse) Apply(c context.Context, w http.ResponseWriter) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(r.Data); err != nil {
		return err
	}
	r.Body = buf
	r.HTTPResponse.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r *DataResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// Download returns a download response to handle file downloads
func (r *Responder) Download(data io.ReadCloser, contentType string, fileName string, forceDownload bool) *HTTPResponse {
	contentDisposition := "inline"
	if forceDownload {
		contentDisposition = "attachement"
	}

	return &HTTPResponse{
		Status: http.StatusOK,
		Header: http.Header{
			"Content-Type":        []string{contentType},
			"Content-Disposition": []string{contentDisposition + "; filename=" + fileName},
		},
		Body: data,
	}
}

// Render creates a render response, with the supplied template and data
func (r *Responder) Render(tpl string, data interface{}) *RenderResponse {
	return &RenderResponse{
		Template:     tpl,
		engine:       r.engine,
		DataResponse: *r.Data(data),
	}
}

// Apply response
func (r *RenderResponse) Apply(c context.Context, w http.ResponseWriter) error {
	var err error

	if req, ok := FromContext(c); ok && r.engine != nil {
		partialRenderer, ok := r.engine.(template.PartialEngine)
		if partials := req.Request().Header.Get("X-Partial"); partials != "" && ok {
			content, err := partialRenderer.RenderPartials(c, r.Template, r.Data, strings.Split(partials, ","))
			body, err := json.Marshal(map[string]interface{}{"partials": content, "data": new(GetPartialDataFunc).Func(c).(func() map[string]interface{})()})
			if err != nil {
				return err
			}
			r.Body = bytes.NewBuffer(body)
			r.Header.Set("Content-Type", "application/json; charset=utf-8")
			return r.HTTPResponse.Apply(c, w)
		}
	}

	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	r.Body, err = r.engine.Render(c, r.Template, r.Data)
	if err != nil {
		return err
	}
	return r.HTTPResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r *RenderResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// Apply response
func (r *ServerErrorResponse) Apply(c context.Context, w http.ResponseWriter) error {
	return r.RenderResponse.Apply(c, w)
}

// Hook helper
// deprecated: to be removed
func (r *ServerErrorResponse) Hook(hooks ...ResponseHook) Response {
	r.HTTPResponse.Hook(hooks...)
	return r
}

// ServerErrorWithCodeAndTemplate error response with template and http status code
func (r *Responder) ServerErrorWithCodeAndTemplate(err error, tpl string, status uint) *ServerErrorResponse {
	return &ServerErrorResponse{
		Error: err,
		RenderResponse: RenderResponse{
			Template: tpl,
			engine:   r.engine,
			DataResponse: DataResponse{
				Data: map[string]interface{}{
					"code": status,
				},
				HTTPResponse: HTTPResponse{
					Status: status,
					Header: make(http.Header),
				},
			},
		},
	}
}

// ServerError creates a 500 error response
func (r *Responder) ServerError(err error) *ServerErrorResponse {
	return r.ServerErrorWithCodeAndTemplate(err, r.templateErrorWithCode, http.StatusInternalServerError)
}

// Unavailable creates a 503 error response
func (r *Responder) Unavailable(err error) *ServerErrorResponse {
	return r.ServerErrorWithCodeAndTemplate(err, r.templateUnavailable, http.StatusServiceUnavailable)
}

// NotFound creates a 404 error response
func (r *Responder) NotFound(err error) *ServerErrorResponse {
	return r.ServerErrorWithCodeAndTemplate(err, r.templateNotFound, http.StatusNotFound)
}

// Forbidden creates a 403 error response
func (r *Responder) Forbidden(err error) *ServerErrorResponse {
	return r.ServerErrorWithCodeAndTemplate(err, r.templateForbidden, http.StatusForbidden)
}

// TODO creates a 501 Not Implemented response
func (r *Responder) TODO() *HTTPResponse {
	return &HTTPResponse{
		Status: http.StatusNotImplemented,
		Header: make(http.Header),
	}
}
