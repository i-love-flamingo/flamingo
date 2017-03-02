package router

import "net/http"

type (
	// VerboseResponseWriter shadows http.ResponseWriter and tracks written bytes and result Status for logging.
	VerboseResponseWriter struct {
		http.ResponseWriter
		Status int
		Size   int
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
