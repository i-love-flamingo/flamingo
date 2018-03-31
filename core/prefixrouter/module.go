package prefixrouter

import (
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/router"
)

// Module for core/prefix_router
type Module struct {
	RootCmd *cobra.Command  `inject:"flamingo"`
	Root    *config.Area    `inject:""`
	Logger  flamingo.Logger `inject:""`
	//Optional Register a PrimaryHandlersHandlers which is passed to the FrontendRouter
	PrimaryHandlers []OptionalHandler `inject:"primaryHandlers,optional"`
	//Optional Register a FallbackHandlers which is passed to the FrontendRouter
	FallbackHandlers []OptionalHandler `inject:"fallback,optional"`
	//This need to be discussed... currently seems to be used as the default HandlerFunc which can be get via sideeffekt to register finalFallbackHandler routes - maybe provide exactly this - a inject for additional fallbackroutes?
	defaultmux *http.ServeMux
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	//Why this indirect sideeffect way? can this not be an injection filled by someone else from outside (project specific module?)
	m.defaultmux = http.NewServeMux()

	var addr string
	var servecmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Run:     Serve(m.Root, m.defaultmux, &addr, m.PrimaryHandlers, m.FallbackHandlers, m.Logger),
	}

	servecmd.Flags().StringVarP(&addr, "addr", "a", ":3210", "addr on which flamingo runs")

	m.RootCmd.AddCommand(servecmd)

	injector.Bind((*http.ServeMux)(nil)).ToInstance(m.defaultmux)
}

// Serve HTTP Requests
func Serve(root *config.Area, defaultRouter *http.ServeMux, addr *string, primaryHandlers []OptionalHandler, fallbackHandlers []OptionalHandler, logger flamingo.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.SetFinalFallbackHandler(defaultRouter)
		frontRouter.SetFallbackHandlers(fallbackHandlers)
		frontRouter.SetPrimaryHandlers(primaryHandlers)

		for _, area := range root.GetFlatContexts() {
			baseurl, ok := area.Configuration.Get("prefixrouter.baseurl")
			if !ok {
				continue
			}
			logger.Println("Routing", area.Name, "at", baseurl)
			area.Injector.Bind((*flamingo.Logger)(nil)).ToInstance(logger.WithField("area", area.Name))
			areaRouter := area.Injector.GetInstance(router.Router{}).(*router.Router)
			areaRouter.Init(area)
			bu, _ := url.Parse("scheme://" + baseurl.(string))
			areaRouter.SetBase(bu)
			frontRouter.Add(baseurl.(string), areaRouter)
		}

		logger.Printf("Starting HTTP Server at %s .....", *addr)
		e := http.ListenAndServe(*addr, frontRouter)
		if e != nil {
			logger.Error("Unexpected Error", e)
		}
	}
}
