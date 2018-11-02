package web

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type (
	ResponseHook func(c context.Context, rw http.ResponseWriter)

	// Response defines the generic web response
	Response interface {
		// Apply executes the response on the http.ResponseWriter
		Apply(context.Context, http.ResponseWriter) error
		GetStatus() int
		GetContentLength() int
		Hook(...ResponseHook) Response
	}

	// BasicResponse defines a response with basic attributes
	BasicResponse struct {
		Status      int
		contentSize int
		hooks       []ResponseHook
	}

	// OnResponse hook
	OnResponse interface {
		OnResponse(context.Context, *Request, http.ResponseWriter)
	}

	// Redirect response with the ability to add data
	Redirect interface {
		Response
		With(key string, data interface{}) Redirect
	}

	// RedirectResponse redirect
	RedirectResponse struct {
		BasicResponse
		Location string
		data     map[string]interface{}
	}

	// ContentResponse contains a response with body
	ContentResponse struct {
		BasicResponse
		Body        io.Reader
		ContentType string
	}

	// JSONResponse returns Data encoded as JSON
	// todo: create a generic data response instead
	JSONResponse struct {
		BasicResponse
		Data interface{}
	}

	// ErrorResponse wraps a response with an error
	ErrorResponse struct {
		Response
		Error error
	}

	// VerboseResponseWriter shadows http.ResponseWriter and tracks written bytes and result Status for logging.
	VerboseResponseWriter struct {
		http.ResponseWriter
		Status int
		Size   int
	}

	// ServeHTTPResponse wraps the original response with a VerboseResponseWriter
	ServeHTTPResponse struct {
		*VerboseResponseWriter
		BasicResponse
	}
)

// Write calls http.ResponseWriter.Write and records the written bytes.
func (response *VerboseResponseWriter) Write(data []byte) (int, error) {
	l, e := response.ResponseWriter.Write(data)
	response.Size += l
	return l, e
}

// WriteHeader calls http.ResponseWriter.WriteHeader and records the Status code.
func (response *VerboseResponseWriter) WriteHeader(h int) {
	response.Status = h
	response.ResponseWriter.WriteHeader(h)
}

// Apply Response (empty, it has already been applied)
func (shr *ServeHTTPResponse) Apply(c context.Context, rw http.ResponseWriter) error {
	return shr.BasicResponse.Apply(c, rw)
}

// Hook appends hooks to the response
func (shr *ServeHTTPResponse) Hook(hooks ...ResponseHook) Response {
	shr.BasicResponse.hooks = append(shr.BasicResponse.hooks, hooks...)
	return shr
}

// GetStatus returns the status of the response
func (br *BasicResponse) GetStatus() int {
	return br.Status
}

// GetContentLength returns the content size of the response
func (br *BasicResponse) GetContentLength() int {
	return br.contentSize
}

// Apply sets status and content size of the response
func (br *BasicResponse) Apply(c context.Context, rw http.ResponseWriter) error {
	if vrb, ok := rw.(*VerboseResponseWriter); ok {
		br.Status = vrb.Status
		br.contentSize = vrb.Size
	}
	return nil
}

// OnResponse callback to apply hooks
func (br *BasicResponse) OnResponse(c context.Context, r *Request, rw http.ResponseWriter) {
	for _, hook := range br.hooks {
		hook(c, rw)
	}
}

// Hook appends hooks to the response
func (br *BasicResponse) Hook(hooks ...ResponseHook) Response {
	br.hooks = append(br.hooks, hooks...)
	return br
}

// Apply Response
func (rr *RedirectResponse) Apply(c context.Context, rw http.ResponseWriter) error {
	if rr.Status == 0 {
		rr.Status = http.StatusTemporaryRedirect
	}

	rw.Header().Set("Location", rr.Location)
	rw.WriteHeader(rr.Status)

	return rr.BasicResponse.Apply(c, rw)
}

// OnResponse Hook
func (rr *RedirectResponse) OnResponse(c context.Context, r *Request, rw http.ResponseWriter) {
	for k, v := range rr.data {
		r.Session().AddFlash(v, k)
	}
}

// With adds data to the web response
func (rr *RedirectResponse) With(key string, data interface{}) Redirect {
	if rr.data == nil {
		rr.data = make(map[string]interface{}, 1)
	}
	rr.data[key] = data

	return rr
}

// Hook appends hooks to the response
func (rr *RedirectResponse) Hook(hooks ...ResponseHook) Response {
	rr.BasicResponse.hooks = append(rr.BasicResponse.hooks, hooks...)
	return rr
}

// Apply ContentResponse
func (cr *ContentResponse) Apply(c context.Context, rw http.ResponseWriter) error {
	if cr.ContentType == "" {
		cr.ContentType = "text/plain; charset=utf-8"
	}

	if cr.Status == 0 {
		cr.Status = http.StatusOK
	}

	rw.Header().Set("Content-Type", cr.ContentType)
	rw.WriteHeader(cr.Status)
	if cr.Body != nil {
		io.Copy(rw, cr.Body)
	}

	return cr.BasicResponse.Apply(c, rw)
}

// Hook appends hooks to the response
func (cr *ContentResponse) Hook(hooks ...ResponseHook) Response {
	cr.BasicResponse.hooks = append(cr.BasicResponse.hooks, hooks...)
	return cr
}

// Apply JSONResponse
func (jr *JSONResponse) Apply(c context.Context, rw http.ResponseWriter) error {
	if jr.Status == 0 {
		jr.Status = http.StatusOK
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(jr.Status)

	p, err := json.Marshal(jr.Data)
	if err != nil {
		return err
	}
	rw.Write(p)

	return jr.BasicResponse.Apply(c, rw)
}

// Hook appends hooks to the response
func (jr *JSONResponse) Hook(hooks ...ResponseHook) Response {
	jr.BasicResponse.hooks = append(jr.BasicResponse.hooks, hooks...)
	return jr
}
