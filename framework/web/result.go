package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// Result defines the generic web response
	Result interface {
		// Apply executes the response on the http.ResponseWriter
		Apply(ctx context.Context, rw http.ResponseWriter) error
	}

	// Responder generates responses
	Responder struct {
		engine flamingo.TemplateEngine
		router *Router
		logger flamingo.Logger
		debug  bool

		templateForbidden     string
		templateNotFound      string
		templateUnavailable   string
		templateErrorWithCode string
	}

	// Response contains a status and a body
	Response struct {
		Status         uint
		Body           io.Reader
		Header         http.Header
		CacheDirective *CacheDirective
	}

	// RouteRedirectResponse redirects to a certain route
	RouteRedirectResponse struct {
		Response
		To       string
		fragment string
		Data     map[string]string
		router   *Router
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
		Error     error
		ErrString string
	}

	// CacheDirectiveBuilder constructs a CacheDirective with the most commonly used options
	CacheDirectiveBuilder struct {
		IsReusable              bool
		RevalidateEachTime      bool
		AllowIntermediateCaches bool
		MaxCacheLifetime        int
		ETag                    string
	}

	// CacheDirective holds the possible directives for Cache Control Headers and other Http Caching
	CacheDirective struct {
		// Visibility: private or public
		// a response marked “private” can be cached (by the browser) but such responses are typically intended for single users hence they aren’t cacheable by intermediate caches
		// A response that is marked “public” can be cached even in cases where it is associated with a HTTP authentication or the HTTP response status code is not cacheable normally. In most cases, a response marked “public” isn’t necessary, since explicit caching information (i.e. “max-age”) shows that a response is cacheable anyway.
		Visibility string
		// NoCache directive to tell this response can’t be used for subsequent requests to the same URL (browser might revalidate or not cache at all)
		NoCache bool
		// NoStore directive: disallows browsers and all intermediate caches from storing any versions of returned responses i.e. responses containing private/personal information or banking data. Every time users request this asset, requests are sent to the server. The assets are downloaded every time.
		NoStore bool
		// tells intermediate caches to not modify headers - especialle The Content-Encoding, Content-Range, and Content-Type headers must remain unchanged
		NoTransform bool
		// MustRevalidate directive is used to tell a cache that it must first revalidate an asset with the origin after it becomes stale
		MustRevalidate bool
		// ProxyRevalidate is the same as MustRevalidate for shared caches
		ProxyRevalidate bool
		// MaxAge defines the max-age directive states the maximum amount of time in seconds that fetched responses are allowed to be used again
		MaxAge int
		// SMaxAge defines the maxAge for shared caches. Supposed to override max-age for CDN for example
		SMaxAge int
		// ETag the key for the Response
		ETag string
		// LastModifiedSince indicates the time a document last changed
		LastModifiedSince *time.Time
	}
)

const (
	// CacheVisibilityPrivate is used as visibility in CacheDirective to indiate no store in intermediate caches
	CacheVisibilityPrivate = "private"
	// CacheVisibilityPublic is used as visibility in CacheDirective to indicate that response can be stored also in intermediate caches
	CacheVisibilityPublic = "public"
)

// Inject Responder dependencies
func (r *Responder) Inject(router *Router, logger flamingo.Logger, cfg *struct {
	Engine                flamingo.TemplateEngine `inject:",optional"`
	Debug                 bool                    `inject:"config:flamingo.debug.mode"`
	TemplateForbidden     string                  `inject:"config:flamingo.template.err403"`
	TemplateNotFound      string                  `inject:"config:flamingo.template.err404"`
	TemplateUnavailable   string                  `inject:"config:flamingo.template.err503"`
	TemplateErrorWithCode string                  `inject:"config:flamingo.template.errWithCode"`
}) *Responder {
	r.engine = cfg.Engine
	r.router = router
	r.templateForbidden = cfg.TemplateForbidden
	r.templateNotFound = cfg.TemplateNotFound
	r.templateUnavailable = cfg.TemplateUnavailable
	r.templateErrorWithCode = cfg.TemplateErrorWithCode
	r.logger = logger.WithField("module", "framework.web").WithField("category", "responder")
	r.debug = cfg.Debug
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
func (r *Response) Apply(_ context.Context, responseWriter http.ResponseWriter) error {
	if r.CacheDirective != nil {
		r.CacheDirective.ApplyHeaders(r.Header)
	}
	for name, vals := range r.Header {
		for _, val := range vals {
			responseWriter.Header().Add(name, val)
		}
	}
	if r.Status == 0 {
		r.Status = http.StatusOK
	}

	responseWriter.WriteHeader(int(r.Status))

	if r.Body == nil {
		return nil
	}

	_, err := io.Copy(responseWriter, r.Body)

	return err
}

// SetNoCache helper
// deprecated: use CacheControlHeader instead
func (r *Response) SetNoCache() *Response {
	r.CacheDirective = CacheDirectiveBuilder{IsReusable: false}.Build()
	return r
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
	if r.router == nil {
		return errors.New("no reverserouter available")
	}

	to, err := r.router.Relative(r.To, r.Data)
	if err != nil {
		return err
	}
	to.Fragment = r.fragment
	w.Header().Set("Location", to.String())
	return r.Response.Apply(c, w)
}

// Fragment adds a fragment to the resulting URL, argument must be given without '#'
func (r *RouteRedirectResponse) Fragment(fragment string) *RouteRedirectResponse {
	r.fragment = fragment
	return r
}

// Permanent marks a redirect as being permanent (http 301)
func (r *RouteRedirectResponse) Permanent() *RouteRedirectResponse {
	r.Status = http.StatusMovedPermanently
	return r
}

// SetNoCache helper
func (r *RouteRedirectResponse) SetNoCache() *RouteRedirectResponse {
	r.Response.SetNoCache()
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
	if r.URL == nil {
		return errors.New("URL is nil")
	}

	w.Header().Set("Location", r.URL.String())
	return r.Response.Apply(c, w)
}

// Permanent marks a redirect as being permanent (http 301)
func (r *URLRedirectResponse) Permanent() *URLRedirectResponse {
	r.Status = http.StatusMovedPermanently
	return r
}

// SetNoCache helper
func (r *URLRedirectResponse) SetNoCache() *URLRedirectResponse {
	r.Response.SetNoCache()
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
	if r.Response.Header == nil {
		r.Response.Header = make(http.Header)
	}
	r.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
	return r.Response.Apply(c, w)
}

// Status changes response status code
func (r *DataResponse) Status(status uint) *DataResponse {
	r.Response.Status = status
	return r
}

// SetNoCache helper
func (r *DataResponse) SetNoCache() *DataResponse {
	r.Response.SetNoCache()
	return r
}

// Download returns a download response to handle file downloads
func (r *Responder) Download(data io.Reader, contentType string, fileName string, forceDownload bool) *Response {
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
			if err != nil {
				return err
			}

			result := make(map[string]string, len(content))
			for k, v := range content {
				buf, err := io.ReadAll(v)
				if err != nil {
					return err
				}
				result[k] = string(buf)
			}

			body, err := json.Marshal(map[string]interface{}{"partials": result, "data": new(GetPartialDataFunc).Func(c).(func() map[string]interface{})()})
			if err != nil {
				return err
			}
			r.Body = bytes.NewBuffer(body)
			r.Header.Set("Content-Type", "application/json; charset=utf-8")
			return r.Response.Apply(c, w)
		}
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set("Content-Type", "text/html; charset=utf-8")
	r.Body, err = r.engine.Render(c, r.Template, r.Data)
	if err != nil {
		return err
	}
	return r.Response.Apply(c, w)
}

// SetNoCache helper
func (r *RenderResponse) SetNoCache() *RenderResponse {
	r.Response.SetNoCache()
	return r
}

// Apply response
func (r *ServerErrorResponse) Apply(c context.Context, w http.ResponseWriter) error {
	if r.RenderResponse.DataResponse.Response.Status == 0 {
		r.RenderResponse.DataResponse.Response.Status = http.StatusInternalServerError
	}

	if err := r.RenderResponse.Apply(c, w); err != nil {
		http.Error(w, r.ErrString, int(r.RenderResponse.DataResponse.Response.Status))
	}

	return nil
}

// ServerErrorWithCodeAndTemplate error response with template and http status code
func (r *Responder) ServerErrorWithCodeAndTemplate(err error, tpl string, status uint) *ServerErrorResponse {
	var errstr string

	if err == nil {
		err = errors.New("")
	}

	if r.debug {
		errstr = fmt.Sprintf("%+v", err)
	} else {
		errstr = err.Error()
	}

	return &ServerErrorResponse{
		Error:     err,
		ErrString: errstr,
		RenderResponse: RenderResponse{
			Template: tpl,
			engine:   r.engine,
			DataResponse: DataResponse{
				Data: map[string]interface{}{
					"code":  status,
					"error": errstr,
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
	if errors.Is(err, context.Canceled) {
		r.getLogger().Debug(fmt.Sprintf("%+v\n", err))
	} else {
		r.getLogger().Error(fmt.Sprintf("%+v\n", err))
	}

	return r.ServerErrorWithCodeAndTemplate(err, r.templateErrorWithCode, http.StatusInternalServerError)
}

// Unavailable creates a 503 error response
func (r *Responder) Unavailable(err error) *ServerErrorResponse {
	r.getLogger().Error(fmt.Sprintf("%+v\n", err))

	return r.ServerErrorWithCodeAndTemplate(err, r.templateUnavailable, http.StatusServiceUnavailable)
}

// NotFound creates a 404 error response
func (r *Responder) NotFound(err error) *ServerErrorResponse {
	r.getLogger().Warn(err)

	return r.ServerErrorWithCodeAndTemplate(err, r.templateNotFound, http.StatusNotFound)
}

// Forbidden creates a 403 error response
func (r *Responder) Forbidden(err error) *ServerErrorResponse {
	r.getLogger().Warn(err)

	return r.ServerErrorWithCodeAndTemplate(err, r.templateForbidden, http.StatusForbidden)
}

// SetNoCache helper
func (r *ServerErrorResponse) SetNoCache() *ServerErrorResponse {
	r.Response.SetNoCache()
	return r
}

// TODO creates a 501 Not Implemented response
func (r *Responder) TODO() *Response {
	return &Response{
		Status: http.StatusNotImplemented,
		Header: make(http.Header),
	}
}

func (r *Responder) getLogger() flamingo.Logger {
	if r.logger != nil {
		return r.logger
	}
	return &flamingo.StdLogger{Logger: *log.New(os.Stdout, "flamingo", log.LstdFlags)}
}

func (r *Responder) completeResult(result Result) Result {
	switch result := result.(type) {
	case *RenderResponse:
		if result.engine == nil {
			result.engine = r.engine
		}
	case *RouteRedirectResponse:
		if result.router == nil {
			result.router = r.router
		}
	case *ServerErrorResponse:
		if result.engine == nil {
			result.engine = r.engine
		}
	}
	return result
}

// ApplyHeaders sets the correct cache control headers
func (c *CacheDirective) ApplyHeaders(header http.Header) {
	var cacheControlValues []string

	if c.NoStore {
		// No store makes all other headers obsolete
		header.Set("Cache-Control", "no-store")
		return
	}

	// Revalidation header:
	if c.MustRevalidate {
		cacheControlValues = append(cacheControlValues, "must-revalidate")
	}
	if c.ProxyRevalidate {
		cacheControlValues = append(cacheControlValues, "proxy-revalidate")
	}
	if c.NoCache {
		cacheControlValues = append(cacheControlValues, "no-cache")
	} else {
		if c.MaxAge > 0 {
			header.Set("Expires", time.Now().Add(time.Duration(int64(c.MaxAge))*time.Second).UTC().Format(time.RFC1123))
			cacheControlValues = append(cacheControlValues, fmt.Sprintf("max-age=%d", c.MaxAge))
		}
		if c.SMaxAge > 0 {
			cacheControlValues = append(cacheControlValues, fmt.Sprintf("s-maxage=%d", c.SMaxAge))
		}
	}

	// Add Validation Headers
	if c.ETag != "" {
		header.Set("ETag", c.ETag)
	}
	if c.LastModifiedSince != nil {
		header.Set("Last-Modified", c.LastModifiedSince.UTC().Format(time.RFC1123))
	}

	// Other directives for caches
	if c.NoTransform {
		cacheControlValues = append(cacheControlValues, "no-transform")
	}
	if c.Visibility == "public" {
		cacheControlValues = append(cacheControlValues, "public")
	}
	if c.Visibility == "private" {
		cacheControlValues = append(cacheControlValues, "private")
	}

	if len(cacheControlValues) > 0 {
		header.Set("Cache-Control", strings.Join(cacheControlValues, ", "))
	}
}

// Build returns the CacheDirective based on the settings
func (c CacheDirectiveBuilder) Build() *CacheDirective {
	if !c.IsReusable {
		return &CacheDirective{
			NoStore: true,
		}
	}
	cd := &CacheDirective{}
	if c.RevalidateEachTime {
		cd.NoCache = true
	}
	if c.AllowIntermediateCaches {
		cd.Visibility = CacheVisibilityPublic
	} else {
		cd.Visibility = CacheVisibilityPrivate
	}
	cd.MaxAge = c.MaxCacheLifetime
	cd.ETag = c.ETag
	return cd
}
