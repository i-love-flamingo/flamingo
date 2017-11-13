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
	RootCmd    *cobra.Command  `inject:"flamingo"`
	Root       *config.Area    `inject:""`
	Logger     flamingo.Logger `inject:""`
	defaultmux *http.ServeMux
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.defaultmux = http.NewServeMux()

	var addr string
	var servecmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Run:     Serve(m.Root, m.defaultmux, &addr, m.Logger),
	}

	servecmd.Flags().StringVarP(&addr, "addr", "a", ":3210", "addr on which flamingo runs")

	m.RootCmd.AddCommand(servecmd)

	injector.Bind((*http.ServeMux)(nil)).ToInstance(m.defaultmux)
}

// Serve HTTP Requests
func Serve(root *config.Area, defaultRouter *http.ServeMux, addr *string, logger flamingo.Logger) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.Default(defaultRouter)

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
