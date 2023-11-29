package prefixrouter

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"flamingo.me/opentelemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
)

// Module for core/prefix_router
type Module struct {
	server                    *http.Server
	eventRouter               flamingo.EventRouter
	logger                    flamingo.Logger
	enableRootRedirectHandler bool
	publicEndpoint            bool
}

// Inject dependencies
func (m *Module) Inject(
	eventRouter flamingo.EventRouter,
	logger flamingo.Logger,
	config *struct {
		EnableRootRedirectHandler bool `inject:"config:flamingo.prefixrouter.rootRedirectHandler.enabled,optional"`
		PublicEndpoint            bool `inject:"config:flamingo.opentelemetry.publicEndpoint,optional"`
	},
) {
	m.eventRouter = eventRouter
	m.logger = logger
	m.enableRootRedirectHandler = config.EnableRootRedirectHandler
	m.publicEndpoint = config.PublicEndpoint
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(http.ServeMux)).ToInstance(http.NewServeMux())
	injector.BindMulti(new(cobra.Command)).ToProvider(m.serveCmd)
	if m.enableRootRedirectHandler {
		injector.BindMulti((*OptionalHandler)(nil)).AnnotatedWith("fallback").To(rootRedirectHandler{})
	}
	flamingo.BindEventSubscriber(injector).ToInstance(m)
}

// CueConfig defines the prefixrouter configuration
func (*Module) CueConfig() string {
	return `
flamingo: prefixrouter: rootRedirectHandler: enabled?: bool
if flamingo.prefixrouter.rootRedirectHandler.enabled {
	flamingo: prefixrouter: rootRedirectHandler: redirectTarget: string
}
`
}

// FlamingoLegacyConfigAlias legacy mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"prefixrouter.rootRedirectHandler.enabled":        "flamingo.prefixrouter.rootRedirectHandler.enabled",
		"prefixrouter.rootRedirectHandler.redirectTarget": "flamingo.prefixrouter.rootRedirectHandler.redirectTarget",
	}
}

type serveCmdConfig struct {
	PrimaryHandlers  []OptionalHandler `inject:"primaryHandlers,optional"` // Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
	FallbackHandlers []OptionalHandler `inject:"fallback,optional"`        // Optional Register a FallbackHandlers which is passed to the FrontendRouter
}

func (m *Module) serveCmd(area *config.Area, defaultmux *http.ServeMux, configuredURLPrefixSampler *opentelemetry.ConfiguredURLPrefixSampler, config *serveCmdConfig) *cobra.Command {
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

// serve HTTP Requests
func (m *Module) serve(
	root *config.Area,
	defaultRouter *http.ServeMux,
	addr *string,
	configuredURLPrefixSampler *opentelemetry.ConfiguredURLPrefixSampler,
	primaryHandlers []OptionalHandler,
	fallbackHandlers []OptionalHandler,
) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.SetFinalFallbackHandler(defaultRouter)
		frontRouter.SetFallbackHandlers(fallbackHandlers)
		frontRouter.SetPrimaryHandlers(primaryHandlers)

		areas, _ := root.GetFlatContexts()
		for _, area := range areas {
			pathValue, pathSet := area.Configuration.Get("flamingo.router.path")
			hostValue, hostSet := area.Configuration.Get("flamingo.router.host")

			if !pathSet && !hostSet {
				m.logger.WithField("category", "prefixrouter").Info("No prefix configured for config area ", area.Name, "!  Area is not routed by prefixrouter!")
				continue
			}

			injector, err := area.GetInitializedInjector()
			if err != nil {
				log.Fatal(err)
			}

			injector.Bind((*flamingo.Logger)(nil)).ToInstance(m.logger.WithField("area", area.Name))
			i, err := injector.GetInstance(web.Router{})
			if err != nil {
				panic(err)
			}
			areaRouter := i.(*web.Router)

			prefix := "/"
			if pathSet {
				prefix = path.Join("/", pathValue.(string), "/")
			}
			if hostSet && hostValue != "" {
				prefix = hostValue.(string) + prefix
			}

			m.logger.WithField("category", "prefixrouter").Info("Routing area '", area.Name, "' at prefix '", prefix, "'")
			frontRouter.Add(prefix, routerHandler{area: area.Name, handler: areaRouter.Handler()})
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

		options := []otelhttp.Option{
			otelhttp.WithFilter(opentelemetry.URLPrefixSampler(whitelist, blacklist,
				configuredURLPrefixSampler.AllowParentTrace)),
		}
		if m.publicEndpoint {
			options = append(options, otelhttp.WithPublicEndpoint())
		}

		m.server = &http.Server{
			Addr:    *addr,
			Handler: otelhttp.NewHandler(frontRouter, "server", options...),
		}

		e := m.listenAndServe()
		if e != nil && e != http.ErrServerClosed {
			m.logger.WithField("category", "prefixrouter").Error("Unexpected Error ", e)
		}
	}
}

func (m *Module) listenAndServe() error {
	listener, err := net.Listen("tcp", m.server.Addr)
	if err != nil {
		return err
	}

	addr := listener.Addr().String()
	m.logger.WithField("category", "prefixrouter").Info(fmt.Sprintf("Starting HTTP Server (Prefixrouter) at %s", addr))

	port := addr[strings.LastIndex(addr, ":")+1:]
	m.eventRouter.Dispatch(context.Background(), &flamingo.ServerStartEvent{Port: port})
	defer m.eventRouter.Dispatch(context.Background(), &flamingo.ServerShutdownEvent{})

	return m.server.Serve(listener)
}

// Notify handles the app shutdown event
func (m *Module) Notify(ctx context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		if m.server == nil {
			m.logger.WithField("category", "prefixrouter").Info("Shutdown: server not started.. ")
			return
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		m.logger.WithField("category", "prefixrouter").Info(fmt.Sprintf("Shutdown server on: %v ", m.server.Addr))

		err := m.server.Shutdown(ctx)
		if err != nil {
			m.logger.WithField("category", "prefixrouter").Error("unexpected error on server shutdown: ", err)
		}
	}
}
