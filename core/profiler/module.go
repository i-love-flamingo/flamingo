package profiler

import (
	"fmt"
	"net/http"

	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/profiler"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		DebugMode      bool             `inject:"config:debug.mode"`

	}

	roundTripper struct {
		next http.RoundTripper
	}
)

// RoundTrip a request
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if ctx, ok := req.Context().(web.Context); ok {
		req.Header.Add("X-Correlation-ID", ctx.ID())
		defer ctx.Profile("http.request", fmt.Sprintf("%s %s", req.Method, req.URL.String()))()
	} else if ctx, ok = req.Context().Value(web.CONTEXT).(web.Context); ok {
		req.Header.Add("X-Correlation-ID", ctx.ID())
		defer ctx.Profile("http.request", fmt.Sprintf("%s %s", req.Method, req.URL.String()))()
	}

	return rt.next.RoundTrip(req)
}

func init() {
	http.DefaultTransport = &roundTripper{
		next: http.DefaultTransport,
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	if (m.DebugMode) {
		m.RouterRegistry.Route("/_profiler/view/:profile", "_profiler.view")
		m.RouterRegistry.Handle("_profiler.view", new(profileController))

		injector.Override((*profiler.Profiler)(nil), "").To(defaultProfiler{})

		injector.BindMulti((*event.Subscriber)(nil)).To(eventSubscriber{})
	}
}
