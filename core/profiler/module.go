package profiler

import (
	"fmt"
	"net/http"

	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/profiler"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
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
	m.RouterRegistry.Route("/_profiler/view/:profile", "_profiler.view")
	m.RouterRegistry.Handle("_profiler.view", new(profileController))

	injector.Override((*profiler.Profiler)(nil), "").To(defaultProfiler{})

	injector.BindMulti((*event.Subscriber)(nil)).To(eventSubscriber{})
}
