package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// Result defines the generic web response
	Result interface {
		// Apply executes the response on the http.ResponseWriter
		Apply(ctx context.Context, rw http.ResponseWriter) error
	}

	// OnResponse hook
	// deprecated: necessary?
	onResponse interface {
		OnResponse(ctx context.Context, r *Request, rw http.ResponseWriter)
	}

	// Responder generates responses
	Responder struct {
		engine flamingo.TemplateEngine
		router *Router

		templateForbidden     string
		templateNotFound      string
		templateUnavailable   string
		templateErrorWithCode string
	}

	// Response contains a status and a body
	Response struct {
		Status uint
		Body   io.Reader
		Header http.Header
	}

	// RouteRedirectResponse redirects to a certain route
	RouteRedirectResponse struct {
		Response
		To     string
		Data   map[string]string
		router *Router
	}

	// URLRedirectResponse redirects to a certain URL
	URLRedirectResponse struct {
		Response
		URL *url.URL
	}

	// DataResponse returns a response containing data, e.g. as JSON
	DataResponse struct {
		Response
		Data interface{}
	}

	// RenderResponse renders data
	RenderResponse struct {
		DataResponse
		Template string
		engine   flamingo.TemplateEngine
	}

	// ServerErrorResponse returns a server error, by default http 500
	ServerErrorResponse struct {
		RenderResponse
		Error error
	}
)

// Inject Responder dependencies
func (r *Responder) Inject(engine flamingo.TemplateEngine, router *Router, cfg *struct {
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

var _ Result = &Response{}

// HTTP Response generator
func (r *Responder) HTTP(status uint, body io.Reader) *Response {
	return &Response{
		Status: status,
		Body:   body,
		Header: make(http.Header),
	}
}

// Apply response
func (r *Response) Apply(c context.Context, w http.ResponseWriter) error {
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

// RouteRedirect generator
func (r *Responder) RouteRedirect(to string, data map[string]string) *RouteRedirectResponse {
	return &RouteRedirectResponse{
		To:     to,
		Data:   data,
		router: r.router,
		Response: Response{
			Status: http.StatusSeeOther,
			Header: make(http.Header),
		},
	}
}

// Apply response
func (r *RouteRedirectResponse) Apply(c context.Context, w http.ResponseWriter) error {
	to, err := r.router.URL(r.To, r.Data)
	if err != nil {
		return err
	}
	w.Header().Set("Location", to.String())
	return r.Response.Apply(c, w)
}

// Permanent marks a redirect as being permanent (http 301)
func (r *RouteRedirectResponse) Permanent() *RouteRedirectResponse {
	r.Status = http.StatusMovedPermanently
	return r
}

// URLRedirect returns a 303 redirect to a given URL
func (r *Responder) URLRedirect(url *url.URL) *URLRedirectResponse {
	return &URLRedirectResponse{
		URL: url,
		Response: Response{
			Status: http.StatusSeeOther,
			Header: make(http.Header),
		},
	}
}

// Apply response
func (r *URLRedirectResponse) Apply(c context.Context, w http.ResponseWriter) error {
	w.Header().Set("Location", r.URL.String())
	return r.Response.Apply(c, w)
}

// Permanent marks a redirect as being permanent (http 301)
func (r *URLRedirectResponse) Permanent() *URLRedirectResponse {
	r.Status = http.StatusMovedPermanently
	return r
}

// Data returns a data response which can be serialized
func (r *Responder) Data(data interface{}) *DataResponse {
	return &DataResponse{
		Data: data,
		Response: Response{
			Status: http.StatusOK,
			Header: make(http.Header),
		},
	}
}

// Apply response
// todo: support more than json
func (r *DataResponse) Apply(c context.Context, w http.ResponseWriter) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(r.Data); err != nil {
		return err
	}
	r.Body = buf
	r.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r.Response.Apply(c, w)
}

// Download returns a download response to handle file downloads
func (r *Responder) Download(data io.ReadCloser, contentType string, fileName string, forceDownload bool) *Response {
	contentDisposition := "inline"
	if forceDownload {
		contentDisposition = "attachement"
	}

	return &Response{
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

	if r.engine == nil {
		return r.DataResponse.Apply(c, w)
	}

	if req := RequestFromContext(c); req != nil && r.engine != nil {
		partialRenderer, ok := r.engine.(flamingo.PartialTemplateEngine)
		if partials := req.Request().Header.Get("X-Partial"); partials != "" && ok {
			content, err := partialRenderer.RenderPartials(c, r.Template, r.Data, strings.Split(partials, ","))
			body, err := json.Marshal(map[string]interface{}{"partials": content, "data": new(GetPartialDataFunc).Func(c).(func() map[string]interface{})()})
			if err != nil {
				return err
			}
			r.Body = bytes.NewBuffer(body)
			r.Header.Set("Content-Type", "application/json; charset=utf-8")
			return r.Response.Apply(c, w)
		}
	}

	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	r.Body, err = r.engine.Render(c, r.Template, r.Data)
	if err != nil {
		return err
	}
	return r.Response.Apply(c, w)
}

// Apply response
func (r *ServerErrorResponse) Apply(c context.Context, w http.ResponseWriter) error {
	return r.RenderResponse.Apply(c, w)
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
				Response: Response{
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
func (r *Responder) TODO() *Response {
	return &Response{
		Status: http.StatusNotImplemented,
		Header: make(http.Header),
	}
}
