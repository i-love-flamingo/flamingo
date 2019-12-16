package prefixrouter

import (
	"net/http"
)

type (
	rootRedirectHandler struct {
		redirectTarget string
	}
)

// Inject
func (r *rootRedirectHandler) Inject(config *struct {
	RedirectTarget string `inject:"config:flamingo.prefixrouter.rootRedirectHandler.redirectTarget,optional"`
}) {
	r.redirectTarget = config.RedirectTarget
}

// TryServeHTTP - implementation of OptionalHandler
func (r *rootRedirectHandler) TryServeHTTP(rw http.ResponseWriter, req *http.Request) (bool, error) {
	if r.redirectTarget == "" || r.redirectTarget == "/" {
		return true, nil
	}
	if req.RequestURI == "/" {
		rw.Header().Set("Location", r.redirectTarget)
		rw.WriteHeader(http.StatusTemporaryRedirect)
		return false, nil
	}
	return true, nil
}
