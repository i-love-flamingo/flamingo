package prefixrouter

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"

	"github.com/spf13/cobra"
)

// Module for core/prefix_router
type Module struct {
	RootCmd    *cobra.Command `inject:"flamingo"`
	Root       *config.Area   `inject:""`
	defaultmux *http.ServeMux
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.defaultmux = http.NewServeMux()

	var port int
	var servecmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server"},
		Run:     Serve(m.Root, m.defaultmux, &port),
	}

	servecmd.Flags().IntVarP(&port, "port", "p", 3210, "port on which flamingo runs")

	m.RootCmd.AddCommand(servecmd)

	injector.Bind((*http.ServeMux)(nil)).ToInstance(m.defaultmux)
}

// Serve HTTP Requests
func Serve(root *config.Area, defaultRouter *http.ServeMux, port *int) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.Default(defaultRouter)

		for _, area := range root.GetFlatContexts() {
			area.Injector.Bind(new(log.Logger)).ToInstance(log.New(os.Stdout, "["+area.Name+"] ", 0))
			log.Println(area.Name, "at", area.BaseURL)
			frontRouter.Add(area.BaseURL, area.Injector.GetInstance(router.Router{}).(*router.Router).Init(area))
		}

		fmt.Println("Starting HTTP Server at :" + strconv.Itoa(*port) + " .....")
		e := http.ListenAndServe(":"+strconv.Itoa(*port), frontRouter)
		if e != nil {
			fmt.Printf("Unexpected Error: %s", e)
		}
	}
}
