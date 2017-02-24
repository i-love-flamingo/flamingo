package prefix_router

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	FrontRouter struct {
		router   map[string]http.Handler
		fallback http.Handler
	}
)

// NewFrontRouter creates new FrontRouter
func NewFrontRouter() *FrontRouter {
	return &FrontRouter{
		router: make(map[string]http.Handler),
	}
}

// Add appends new Handler to Frontrouter
func (fr *FrontRouter) Add(prefix string, handler http.Handler) {
	fr.router[prefix] = handler
}

// Default sets Fallback for undefined Handler
func (fr *FrontRouter) Default(handler http.Handler) {
	fr.fallback = handler
}

// ServeHTTP gets Router for Request and lets it handle it
func (fr *FrontRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	if strings.Index(host, ":") > -1 {
		host = strings.Split(host, ":")[0]
	}

	test := path.Join(host, req.RequestURI)
	for prefix, router := range fr.router {
		if strings.HasPrefix(test, prefix) {
			req.URL, _ = url.Parse(test[len(prefix):])
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
			router.ServeHTTP(w, req)
			return
		}
	}

	test = req.RequestURI
	for prefix, router := range fr.router {
		if strings.HasPrefix(test, prefix) {
			req.URL, _ = url.Parse(test[len(prefix):])
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
			router.ServeHTTP(w, req)
			return
		}
	}

	if fr.fallback != nil {
		fr.fallback.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
	}
}
