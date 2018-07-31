package web

import (
	"context"
	"net/http"
)

// NoCache is a response hook to enforce no cache headers
func NoCache(c context.Context, r *Request, rw http.ResponseWriter) {
	rw.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate, no-store")
}

// AddHeader adds headers in a response hook
func AddHeader(k, v string) ResponseHook {
	return func(c context.Context, r *Request, rw http.ResponseWriter) {
		rw.Header().Add(k, v)
	}
}
