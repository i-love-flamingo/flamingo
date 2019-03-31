package prefixrouter

import (
	"context"
	"net/http"
	"strings"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
)

// Module for core/prefix_router
type Module struct {
	server                    *http.Server
	logger                    flamingo.Logger
	enableRootRedirectHandler bool
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(http.ServeMux)).ToInstance(http.NewServeMux())
	injector.BindMulti(new(cobra.Command)).ToProvider(serveCmd(m))
	if m.enableRootRedirectHandler {
		injector.BindMulti((*OptionalHandler)(nil)).AnnotatedWith("fallback").To(rootRedirectHandler{})
	}
	flamingo.BindEventSubscriber(injector).ToInstance(m)
}

// Inject dependencies
func (m *Module) Inject(l flamingo.Logger, config *struct {
	EnableRootRedirectHandler bool `inject:"config:prefixrouter.rootRedirectHandler.enabled,optional"`
}) {
	m.logger = l
	m.enableRootRedirectHandler = config.EnableRootRedirectHandler
}

func serveCmd(m *Module) func(area *config.Area, defaultmux *http.ServeMux, configuredURLPrefixSampler *opencensus.ConfiguredURLPrefixSampler, config *struct {
	PrimaryHandlers  []OptionalHandler `inject:"primaryHandlers,optional"` // Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
	FallbackHandlers []OptionalHandler `inject:"fallback,optional"`        // Optional Register a FallbackHandlers which is passed to the FrontendRouter
}) *cobra.Command {
	return func(area *config.Area, defaultmux *http.ServeMux, configuredURLPrefixSampler *opencensus.ConfiguredURLPrefixSampler, config *struct {
		PrimaryHandlers  []OptionalHandler `inject:"primaryHandlers,optional"` // Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
		FallbackHandlers []OptionalHandler `inject:"fallback,optional"`        // Optional Register a FallbackHandlers which is passed to the FrontendRouter
	}) *cobra.Command {
		var addr string

		cmd := &cobra.Command{
			Use:     "serve",
			Short:   "run the prefix router",
			Aliases: []string{"server"},
			Run:     m.serve(area, defaultmux, &addr, configuredURLPrefixSampler, config.PrimaryHandlers, config.FallbackHandlers),
		}

		cmd.Flags().StringVarP(&addr, "addr", "a", ":3210", "addr on which flamingo runs")

		return cmd
	}
}

// serve HTTP Requests
func (m *Module) serve(
	root *config.Area,
	defaultRouter *http.ServeMux,
	addr *string,
	configuredURLPrefixSampler *opencensus.ConfiguredURLPrefixSampler,
	primaryHandlers []OptionalHandler,
	fallbackHandlers []OptionalHandler) func(cmd *cobra.Command, args []string,
) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.SetFinalFallbackHandler(defaultRouter)
		frontRouter.SetFallbackHandlers(fallbackHandlers)
		frontRouter.SetPrimaryHandlers(primaryHandlers)

		areas, _ := root.GetFlatContexts()
		for _, area := range areas {
			b, _ := area.Configuration.Get("flamingo.router.path")
			baseURL, _ := b.(string)

			if strings.HasPrefix(baseURL, "/") {
				baseURL = "scheme://host" + baseURL
			}
			if baseURL == "" {
				m.logger.WithField("category", "prefixrouter").Warn("No baseurl configured for config area %v", area.Name)
				continue
			}

			m.logger.WithField("category", "prefixrouter").Info("Routing ", area.Name, " at ", baseURL)
			area.Injector.Bind((*flamingo.Logger)(nil)).ToInstance(m.logger.WithField("area", area.Name))
			areaRouter := area.Injector.GetInstance(web.Router{}).(*web.Router)
			frontRouter.Add(areaRouter.Base().Host+areaRouter.Base().Path, routerHandler{area: area.Name, handler: areaRouter.Handler()})
		}

		whitelist := make([]string, 0, len(configuredURLPrefixSampler.Whitelist)*len(frontRouter.router)+1)
		blacklist := make([]string, 0, len(configuredURLPrefixSampler.Blacklist)*len(frontRouter.router)+1)

		// default routes
		for _, p := range configuredURLPrefixSampler.Whitelist {
			whitelist = append(whitelist, p.(string))
		}
		for _, p := range configuredURLPrefixSampler.Blacklist {
			blacklist = append(blacklist, p.(string))
		}

		// prefixed routes
		for k := range frontRouter.router {
			for _, p := range configuredURLPrefixSampler.Whitelist {
				whitelist = append(whitelist, k+p.(string))
			}
			for _, p := range configuredURLPrefixSampler.Blacklist {
				blacklist = append(blacklist, k+p.(string))
			}
		}

		m.logger.WithField("category", "prefixrouter").Info("Starting HTTP Server (Prefixrouter) at ", *addr, ".....")
		m.server = &http.Server{
			Addr: *addr,
			Handler: &ochttp.Handler{
				IsPublicEndpoint: true,
				Handler:          frontRouter,
				GetStartOptions:  opencensus.URLPrefixSampler(whitelist, blacklist, configuredURLPrefixSampler.AllowParentTrace),
			},
		}

		e := m.server.ListenAndServe()
		if e != nil && e != http.ErrServerClosed {
			m.logger.WithField("category", "prefixrouter").Error("Unexpected Error ", e)
		}
	}
}

// Notify handles the app shutdown event
func (m *Module) Notify(ctx context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ServerShutdownEvent); ok {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		m.logger.WithField("category", "prefixrouter").Info("Shutdown server on ", m.server.Addr)

		err := m.server.Shutdown(ctx)
		if err != nil {
			m.logger.WithField("category", "prefixrouter").Error("unexpected error on server shutdown: ", err)
		}
	}
}
