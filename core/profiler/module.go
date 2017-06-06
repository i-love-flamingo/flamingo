package profiler

import (
	"flamingo/framework/dingo"
	"flamingo/framework/event"
	"flamingo/framework/profiler"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"fmt"
	"net/http"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}

	RoundTripper struct {
		Next http.RoundTripper
	}
)

// RoundTrip a request
func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if ctx, ok := req.Context().(web.Context); ok {
		req.Header.Add("X-Request-ID", ctx.ID())
		defer ctx.Profile("http.request", fmt.Sprintf("%s %s", req.Method, req.URL.String()))()
	}

	return rt.Next.RoundTrip(req)
}

func init() {
	http.DefaultTransport = &RoundTripper{
		Next: http.DefaultTransport,
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Route("/_profiler/view/:profile", "_profiler.view")
	m.RouterRegistry.Handle("_profiler.view", new(ProfileController))

	injector.Override((*profiler.Profiler)(nil), "").To(DefaultProfiler{})

	injector.BindMulti((*event.Subscriber)(nil)).To(EventSubscriber{})
}
