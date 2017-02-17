package web

import (
	"encoding/json"
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

type JsonResponse struct {
	Data interface{}
}

// Apply Response
func (js JsonResponse) Apply(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	p, err := json.Marshal(js.Data)
	if err != nil {
		panic(err)
	}
	rw.Write(p)
}
