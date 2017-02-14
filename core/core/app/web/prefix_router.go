package web

import (
	"io"
	"net"
	"net/http"
	"path"
	"strings"
)

type PrefixRouter map[string]http.Handler

func (pr PrefixRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "/__webpack_hmr" {
		h := w.(http.Hijacker)
		_, buf, _ := h.Hijack()

		c, _ := net.Dial("tcp", "localhost:1337")
		c.Write([]byte("GET /__webpack_hmr HTTP/1.1\r\nHost: localhost\r\n\r\n"))

		go func() {
			io.Copy(buf, c)
		}()
		io.Copy(c, buf)

		return
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
	w.WriteHeader(404)
	//panic(test + " not routable")
}
