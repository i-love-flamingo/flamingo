package web

import (
	"io"
	"net/http"
)

type (
	Response interface {
		Apply(http.ResponseWriter)
	}
)

// RedirectResponse redirect
type RedirectResponse struct {
	Status   int
	Location string
}

// Apply Response
func (rr RedirectResponse) Apply(rw http.ResponseWriter) {
	rw.Header().Set("Location", rr.Location)
	rw.WriteHeader(rr.Status)
}

// ContentResponse contains a response with body
type ContentResponse struct {
	Status      int
	Body        io.Reader
	ContentType string
}

// Apply Response
func (cr ContentResponse) Apply(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", cr.ContentType)
	rw.WriteHeader(cr.Status)
	io.Copy(rw, cr.Body)
}
