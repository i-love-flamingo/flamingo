package prefixrouter

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// Module for core/prefix_router
type Module struct {
	server *http.Server
	logger flamingo.Logger
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(http.ServeMux)).ToInstance(http.NewServeMux())
	injector.BindMulti(new(cobra.Command)).ToProvider(serveCmd(m))
	flamingo.BindEventSubscriber(injector).ToInstance(m)
}

// Inject dependencies
func (m *Module) Inject(l flamingo.Logger) {
	m.logger = l
}

func serveCmd(m *Module) func(area *config.Area, defaultmux *http.ServeMux, config *struct {
	PrimaryHandlers  []OptionalHandler `inject:"primaryHandlers,optional"` // Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
	FallbackHandlers []OptionalHandler `inject:"fallback,optional"`        // Optional Register a FallbackHandlers which is passed to the FrontendRouter
}) *cobra.Command {
	return func(area *config.Area, defaultmux *http.ServeMux, config *struct {
		PrimaryHandlers  []OptionalHandler `inject:"primaryHandlers,optional"` // Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
		FallbackHandlers []OptionalHandler `inject:"fallback,optional"`        // Optional Register a FallbackHandlers which is passed to the FrontendRouter
	}) *cobra.Command {
		var addr string

		cmd := &cobra.Command{
			Use:     "serve",
			Short:   "run the prefix router",
			Aliases: []string{"server"},
			Run:     m.serve(area, defaultmux, &addr, config.PrimaryHandlers, config.FallbackHandlers),
		}

		cmd.Flags().StringVarP(&addr, "addr", "a", ":3210", "addr on which flamingo runs")

		return cmd
	}
}

// serve HTTP Requests
func (m *Module) serve(root *config.Area, defaultRouter *http.ServeMux, addr *string, primaryHandlers []OptionalHandler, fallbackHandlers []OptionalHandler) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.SetFinalFallbackHandler(defaultRouter)
		frontRouter.SetFallbackHandlers(fallbackHandlers)
		frontRouter.SetPrimaryHandlers(primaryHandlers)

		areas, _ := root.GetFlatContexts()
		for _, area := range areas {
			baseurlVal, ok := area.Configuration.Get("prefixrouter.baseurl")

			if !ok {
				m.logger.WithField("category", "prefixrouter").Warn("No baseurl configured for config area %v", area.Name)
				continue
			}

			baseurl := baseurlVal.(string)

			if strings.HasPrefix(baseurl, "/") {
				baseurl = "host" + baseurl
			}

			m.logger.WithField("category", "prefixrouter").Info("Routing ", area.Name, " at ", baseurl)
			area.Injector.Bind(new(flamingo.Logger)).ToInstance(m.logger.WithField("area", area.Name))
			areaRouter := area.Injector.GetInstance(web.Router{}).(*web.Router)
			areaRouter.Init(area)
			bu, _ := url.Parse("scheme://" + baseurl)

			areaRouter.SetBase(bu)
			frontRouter.Add(bu.Path, routerHandler{area: area.Name, handler: areaRouter})
		}

		m.logger.WithField("category", "prefixrouter").Info("Starting HTTP Server (Prefixrouter) at ", *addr, ".....")
		m.server = &http.Server{
			Addr:    *addr,
			Handler: &ochttp.Handler{IsPublicEndpoint: true, Handler: frontRouter, StartOptions: trace.StartOptions{Sampler: opencensus.Sampler}},
		}

		e := m.server.ListenAndServe()
		if e != nil && e != http.ErrServerClosed {
			m.logger.WithField("category", "prefixrouter").Error("Unexpected Error ", e)
		}
	}
}

// Notify handles the app shutdown event
func (m *Module) Notify(event flamingo.Event) {
	switch event.(type) {
	case *flamingo.ShutdownEvent:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		m.logger.WithField("category", "prefixrouter").Info("Shutdown server on ", m.server.Addr)

		err := m.server.Shutdown(ctx)
		if err != nil {
			m.logger.WithField("category", "prefixrouter").Error("unexpected error on server shutdown: ", err)
		}
	}
}
