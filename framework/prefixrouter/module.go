package prefixrouter

import (
	"net/http"
	"net/url"
	"strings"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
)

// Module for core/prefix_router
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*http.ServeMux)(nil)).ToInstance(http.NewServeMux())
	injector.BindMulti(new(cobra.Command)).ToProvider(serveCmd)
}

func serveCmd(area *config.Area, logger flamingo.Logger, defaultmux *http.ServeMux, config *struct {
	PrimaryHandlers  []OptionalHandler `inject:"primaryHandlers,optional"` //Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
	FallbackHandlers []OptionalHandler `inject:"fallback,optional"`        //Optional Register a FallbackHandlers which is passed to the FrontendRouter
}) *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "run the prefix router",
		Aliases: []string{"server"},
		Run:     serve(area, defaultmux, &addr, config.PrimaryHandlers, config.FallbackHandlers, logger),
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":3210", "addr on which flamingo runs")

	return cmd
}

// serve HTTP Requests
func serve(root *config.Area, defaultRouter *http.ServeMux, addr *string, primaryHandlers []OptionalHandler, fallbackHandlers []OptionalHandler, logger flamingo.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.SetFinalFallbackHandler(defaultRouter)
		frontRouter.SetFallbackHandlers(fallbackHandlers)
		frontRouter.SetPrimaryHandlers(primaryHandlers)

		for _, area := range root.GetFlatContexts() {
			baseurlVal, ok := area.Configuration.Get("prefixrouter.baseurl")

			if !ok {
				logger.WithField("category", "prefixrouter").Warn("No baseurl configured for config area %v", area.Name)
				continue
			}

			baseurl := baseurlVal.(string)

			if strings.HasPrefix(baseurl, "/") {
				baseurl = "host" + baseurl
			}

			logger.WithField("category", "prefixrouter").Info("Routing", area.Name, "at", baseurl)
			area.Injector.Bind((*flamingo.Logger)(nil)).ToInstance(logger.WithField("area", area.Name))
			areaRouter := area.Injector.GetInstance(router.Router{}).(*router.Router)
			areaRouter.Init(area)
			bu, _ := url.Parse("scheme://" + baseurl)

			areaRouter.SetBase(bu)
			frontRouter.Add(bu.Path, routerHandler{area: area.Name, handler: areaRouter})
		}

		logger.Info("Starting HTTP Server at %s .....", *addr)

		e := http.ListenAndServe(*addr, &ochttp.Handler{IsPublicEndpoint: true, Handler: frontRouter})
		if e != nil {
			logger.Error("Unexpected Error", e)
		}
	}
}
