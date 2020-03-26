package prefixrouter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var rt = stats.Int64("flamingo/prefixrouter/requesttimes", "prefixrouter request times", stats.UnitMilliseconds)

func init() {
	if err := view.Register(
		&view.View{
			Name:        "flamingo/prefixrouter/requests",
			Description: "request times",
			Aggregation: view.Distribution(10, 100, 500, 1000, 2500, 5000, 10000),
			Measure:     rt,
			TagKeys:     []tag.Key{opencensus.KeyArea},
		},
	); err != nil {
		panic(err)
	}
}

type (
	// FrontRouter is a http.handler which serves multiple sites based on the host/path prefix
	FrontRouter struct {
		// primaryHandlers a list of handlers used before processing
		primaryHandlers []OptionalHandler
		// router registered to serve the request based on the prefix
		router map[string]routerHandler
		// fallbackHandlers is used if no router is matching
		fallbackHandlers []OptionalHandler
		// finalFallbackHandler is used as final fallback handler - which is called if no other handler can process
		finalFallbackHandler http.Handler
	}

	routerHandler struct {
		area    string
		handler http.Handler
	}

	// OptionalHandler tries to handle a request
	OptionalHandler interface {
		TryServeHTTP(rw http.ResponseWriter, req *http.Request) (proceed bool, err error)
	}
)

// NewFrontRouter creates new FrontRouter
func NewFrontRouter() *FrontRouter {
	return &FrontRouter{
		router: make(map[string]routerHandler),
	}
}

// Add appends new Handler to Frontrouter
func (fr *FrontRouter) Add(prefix string, handler routerHandler) {
	if h, alreadyPresent := fr.router[prefix]; alreadyPresent {
		panic(
			fmt.Sprintf(
				"prefixrouter: duplicate handler registration on prefix %q from areas %q and %q",
				prefix,
				h.area,
				handler.area,
			),
		)
	}
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

// SetPrimaryHandlers sets list of optional fallback Handlers
func (fr *FrontRouter) SetPrimaryHandlers(handlers []OptionalHandler) {
	fr.primaryHandlers = handlers
}

// ServeHTTP gets Router for Request and lets it handle it
func (fr *FrontRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r, _ := tag.New(req.Context(), tag.Insert(opencensus.KeyArea, "-"))
	req = req.WithContext(r)

	start := time.Now()
	defer func() {
		stats.Record(req.Context(), rt.M(time.Since(start).Nanoseconds()/1000000))
	}()

	ctx, span := trace.StartSpan(req.Context(), "prefixrouter/ServeHTTP")
	req = req.WithContext(ctx)
	defer span.End()

	// process registered primaryHandlers - and if they are successful exist
	for _, handler := range fr.primaryHandlers {
		proceed, _ := handler.TryServeHTTP(w, req)
		if !proceed {
			return
		}
	}

	host := req.Host
	if strings.Index(host, ":") > -1 {
		host = strings.Split(host, ":")[0]
	}

	path := req.RequestURI
	path = "/" + strings.TrimLeft(path, "/")

	var matchedPrefixes []string
	for prefix := range fr.router {
		if strings.HasPrefix(host+path, prefix) {
			matchedPrefixes = append(matchedPrefixes, prefix)
		}
	}
	if len(matchedPrefixes) > 0 {
		prefix := longest(matchedPrefixes)
		router := fr.router[prefix]

		r, _ := tag.New(req.Context(), tag.Upsert(opencensus.KeyArea, router.area))
		req = req.WithContext(r)

		req.URL, _ = url.Parse(path[len(prefix)-len(host):])
		req.URL.Path = "/" + strings.TrimLeft(req.URL.Path, "/")

		span.End()
		router.handler.ServeHTTP(w, req)
		return
	}

	matchedPrefixes = nil
	for prefix := range fr.router {
		if strings.HasPrefix(path, prefix) {
			matchedPrefixes = append(matchedPrefixes, prefix)
		}
	}
	if len(matchedPrefixes) > 0 {
		prefix := longest(matchedPrefixes)
		router := fr.router[prefix]
		r, _ := tag.New(req.Context(), tag.Upsert(opencensus.KeyArea, router.area))
		req = req.WithContext(r)

		req.URL, _ = url.Parse(path[len(prefix):])
		req.URL.Path = "/" + strings.TrimLeft(req.URL.Path, "/")

		span.End()
		router.handler.ServeHTTP(w, req)
		return
	}

	// process registered fallbackHandlers - and if they are successful exist
	for _, handler := range fr.fallbackHandlers {
		proceed, _ := handler.TryServeHTTP(w, req)
		if !proceed {
			return
		}
	}

	// fallback to final handler if given
	if fr.finalFallbackHandler != nil {
		span.End()
		fr.finalFallbackHandler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
	}
}

func longest(strings []string) string {
	var best string
	var length int

	for _, s := range strings {
		if len(s) > length {
			best = s
			length = len(s)
		}
	}

	return best
}
