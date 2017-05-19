package prefix_router

import (
	"flamingo/framework/context"
	"flamingo/framework/dingo"
	"log"

	"net/http"
	"os"

	"fmt"

	"flamingo/framework/router"

	"strconv"

	"github.com/spf13/cobra"
)

// Module for core/prefix_router
type Module struct {
	RootCmd    *cobra.Command   `inject:"flamingo"`
	Root       *context.Context `inject:""`
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
func Serve(root *context.Context, defaultRouter *http.ServeMux, port *int) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		frontRouter := NewFrontRouter()
		frontRouter.Default(defaultRouter)

		for _, ctx := range root.GetFlatContexts() {
			ctx.Injector.Bind(new(log.Logger)).ToInstance(log.New(os.Stdout, "["+ctx.Name+"] ", 0))
			log.Println(ctx.Name, "at", ctx.BaseURL)
			frontRouter.Add(ctx.BaseURL, ctx.Injector.GetInstance(router.Router{}).(*router.Router).Init(ctx))
		}

		fmt.Println("Starting HTTP Server at :" + strconv.Itoa(*port) + " .....")
		e := http.ListenAndServe(":"+strconv.Itoa(*port), frontRouter)
		if e != nil {
			fmt.Printf("Unexpected Error: %s", e)
		}
	}
}
