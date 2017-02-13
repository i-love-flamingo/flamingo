package web

import (
	"net/http"
	"path"
	"strings"
)

type PrefixRouter map[string]http.Handler

func (pr PrefixRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	test := path.Join(req.Host, req.RequestURI)
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
	panic(test + " not routable")
}
