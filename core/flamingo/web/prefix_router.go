package web

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

func NewFrontRouter() *FrontRouter {
	return &FrontRouter{
		router: make(map[string]http.Handler),
	}
}

func (fr *FrontRouter) Add(prefix string, handler http.Handler) {
	fr.router[prefix] = handler
}

func (fr *FrontRouter) Default(handler http.Handler) {
	fr.fallback = handler
}

func (fr *FrontRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	/*
		if strings.HasPrefix(req.RequestURI, "/assets/") {
			if r, e := http.Get("http://localhost:1337" + req.RequestURI); e == nil {
				io.Copy(w, r.Body)
				return
			}
		}
	*/

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
