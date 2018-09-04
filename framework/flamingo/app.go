package flamingo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/router"
	"github.com/spf13/cobra"
)

type (
	appmodule struct {
		root   *config.Area
		router *router.Router
		server *http.Server
		logger Logger
	}

	// AppShutdownEvent is dispatched on app shutdown
	AppShutdownEvent struct {
		AppModule dingo.Module
	}
)

func (a *appmodule) Inject(root *config.Area, router *router.Router, logger Logger) {
	a.root = root
	a.router = router
	a.logger = logger
	a.server = &http.Server{
		Addr:    ":3322",
		Handler: a.router,
	}
}

// Configure dependency injection
func (a *appmodule) Configure(injector *dingo.Injector) {
	injector.BindMulti((*event.Subscriber)(nil)).ToInstance(a)
	//pass a function that returns the Command
	injector.BindMulti(new(cobra.Command)).ToProvider(func() *cobra.Command {
		return serveProvider(a, a.logger)
	})
}

func serveProvider(a *appmodule, logger Logger) *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Default serve command - starts on Port 3322",
		Run: func(cmd *cobra.Command, args []string) {
			a.router.Init(a.root)
			err := a.server.ListenAndServe()
			if err != nil {
				logger.Fatal("unexpected error in serving:", err)
			}
			logger.Info(fmt.Sprintf("Starting HTTP Server at %s .....", addr))
		},
	}
	cmd.Flags().StringVarP(&a.server.Addr, "addr", "a", ":3322", "addr on which flamingo runs")

	return cmd
}

func (a *appmodule) Notify(event event.Event) {
	switch event.(type) {
	case *AppShutdownEvent:
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		a.logger.Info("Shutdown server on ", a.server.Addr)

		err := a.server.Shutdown(ctx)
		if err != nil {
			a.logger.Error("unexpected error on server shutdown: ", err)
		}
	}
}

// App is a simple app-runner for flamingo
func App(root *config.Area, configdir string) {
	app := new(appmodule)
	root.Modules = append(root.Modules, app)
	if configdir == "" {
		configdir = "config"
	}
	config.Load(root, configdir)

	if err := root.Injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command).Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
