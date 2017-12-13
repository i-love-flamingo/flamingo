package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type (
	// Response defines the generic web response
	Response interface {
		// Apply executes the response on the http.ResponseWriter
		Apply(Context, http.ResponseWriter)
		GetStatus() int
		GetContentLength() int
	}

	BasicResponse struct {
		Status      int
		ContentSize int
	}

	// OnResponse hook
	OnResponse interface {
		OnResponse(Context, http.ResponseWriter)
	}

	// Redirect response with the ability to add data
	Redirect interface {
		Response
		With(key string, data interface{}) Redirect
	}

	// RedirectResponse redirect
	RedirectResponse struct {
		BasicResponse
		Status   int
		Location string
		data     map[string]interface{}
	}

	// ContentResponse contains a response with body
	ContentResponse struct {
		BasicResponse
		Status      int
		Body        io.Reader
		ContentType string
	}

	// JSONResponse returns Data encoded as JSON
	JSONResponse struct {
		BasicResponse
		Data   interface{}
		Status int
	}

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
func (shr *ServeHTTPResponse) Apply(c Context, rw http.ResponseWriter) {
	shr.BasicResponse.Apply(c, rw)
}

func (br *BasicResponse) GetStatus() int {
	return br.Status
}

func (br *BasicResponse) GetContentLength() int {
	return br.ContentSize
}

func (br *BasicResponse) Apply(c Context, rw http.ResponseWriter) {
	if vrb, ok := rw.(*VerboseResponseWriter); ok {
		br.Status = vrb.Status
		br.ContentSize = vrb.Size
	}
}

// Apply Response
func (rr *RedirectResponse) Apply(c Context, rw http.ResponseWriter) {
	if rr.Status == 0 {
		rr.Status = http.StatusTemporaryRedirect
	}

	rw.Header().Set("Location", rr.Location)
	rw.WriteHeader(rr.Status)

	rr.BasicResponse.Apply(c, rw)
}

// OnResponse Hook
func (rr *RedirectResponse) OnResponse(c Context, rw http.ResponseWriter) {
	for k, v := range rr.data {
		c.Session().AddFlash(v, k)
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

// Apply ContentResponse
func (cr *ContentResponse) Apply(c Context, rw http.ResponseWriter) {
	if cr.ContentType == "" {
		cr.ContentType = "text/plain; charset=utf-8"
	}

	if cr.Status == 0 {
		cr.Status = http.StatusOK
	}

	rw.Header().Set("Content-Type", cr.ContentType)
	rw.WriteHeader(cr.Status)
	io.Copy(rw, cr.Body)

	cr.BasicResponse.Apply(c, rw)
}

// Apply JSONResponse
func (js *JSONResponse) Apply(c Context, rw http.ResponseWriter) {
	if js.Status == 0 {
		js.Status = http.StatusOK
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(js.Status)

	p, err := json.Marshal(js.Data)
	if err != nil {
		panic(err)
	}
	rw.Write(p)

	js.BasicResponse.Apply(c, rw)
}
