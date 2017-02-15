package web

import (
	"io"
	"net/http"
	"path"
	"strings"
)

type PrefixRouter map[string]http.Handler

func (pr PrefixRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.RequestURI, "/assets/") {
		if r, e := http.Get("http://localhost:1337" + req.RequestURI); e == nil {
			io.Copy(w, r.Body)
			return
		}
	}

	host := req.Host
	if strings.Index(host, ":") > -1 {
		host = strings.Split(host, ":")[0]
	}
	test := path.Join(host, req.RequestURI)
	for prefix, router := range pr {
		if strings.HasPrefix(test, prefix) {
			req.URL.Path = test[len(prefix):]
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
			router.ServeHTTP(w, req)
			return
		}
	}

	test = req.RequestURI
	for prefix, router := range pr {
		if strings.HasPrefix(test, prefix) {
			req.URL.Path = test[len(prefix):]
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
			router.ServeHTTP(w, req)
			return
		}
	}

	w.WriteHeader(404)
	//panic(test + " not routable")
}
