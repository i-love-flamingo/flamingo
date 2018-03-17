package prefixrouter

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	// FrontRouter is a http.handler which serves multiple sites based on the host/path prefix
	FrontRouter struct {
		//primaryHandlers a list of handlers used before processing
		primaryHandlers []OptionalHandler
		//router registered to serve the request based on the prefix
		router map[string]http.Handler
		//fallbackHandlers is used if no router is matching
		fallbackHandlers []OptionalHandler
		//finalFallbackHandler is used as final fallback handler - which is called if no other handler can process
		finalFallbackHandler http.Handler
	}

	OptionalHandler interface {
		TryServeHTTP(rw http.ResponseWriter, req *http.Request) error
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

// SetFinalFallbackHandler sets Fallback for undefined Handler
func (fr *FrontRouter) SetFinalFallbackHandler(handler http.Handler) {
	fr.finalFallbackHandler = handler
}

// SetFallbackHandlers sets list of optional fallback Handlers
func (fr *FrontRouter) SetFallbackHandlers(handlers []OptionalHandler) {
	fr.fallbackHandlers = handlers
}

// SetPrimarykHandlers sets list of optional fallback Handlers
func (fr *FrontRouter) SetPrimaryHandlers(handlers []OptionalHandler) {
	fr.primaryHandlers = handlers
}

// ServeHTTP gets Router for Request and lets it handle it
func (fr *FrontRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	//process registered primaryHandlers - and if they are sucessfull exist
	for _, handler := range fr.primaryHandlers {
		err := handler.TryServeHTTP(w, req)
		if err == nil {
			return
		}
	}

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

	//process registered fallbackHandlers - and if they are sucessfull exist
	for _, handler := range fr.fallbackHandlers {
		err := handler.TryServeHTTP(w, req)
		if err == nil {
			return
		}
	}

	//fallback to final handler if given
	if fr.finalFallbackHandler != nil {
		fr.finalFallbackHandler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
	}
}
