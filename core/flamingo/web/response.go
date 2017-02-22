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
		Apply(http.ResponseWriter)
	}

	// RedirectResponse redirect
	RedirectResponse struct {
		Status   int
		Location string
	}

	// ContentResponse contains a response with body
	ContentResponse struct {
		Status      int
		Body        io.Reader
		ContentType string
	}

	// JsonResponse returns Data encoded as Json
	JsonResponse struct {
		Data   interface{}
		Status int
	}
)

// Apply Response
func (rr RedirectResponse) Apply(rw http.ResponseWriter) {
	if rr.Status == 0 {
		rr.Status = http.StatusTemporaryRedirect
	}

	rw.Header().Set("Location", rr.Location)
	rw.WriteHeader(rr.Status)
}

// Apply Response
func (cr ContentResponse) Apply(rw http.ResponseWriter) {
	if cr.ContentType == "" {
		cr.ContentType = "text/plain"
	}
	if cr.Status == 0 {
		cr.Status = http.StatusOK
	}

	rw.Header().Set("Content-Type", cr.ContentType)
	rw.WriteHeader(cr.Status)
	io.Copy(rw, cr.Body)
}

// Apply Response
func (js JsonResponse) Apply(rw http.ResponseWriter) {
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
}
