package flamingo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"flamingo.me/flamingo/v3/core/cmd"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/event"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/router"
	"github.com/spf13/cobra"
)

type (
	appmodule struct {
		root   *config.Area
		router *router.Router
		server *http.Server
		logger flamingo.Logger
	}

	// AppShutdownEvent is dispatched on app startup
	AppStartupEvent struct {
		AppModule dingo.Module
	}

	// AppShutdownEvent is dispatched on app shutdown
	AppShutdownEvent struct {
		AppModule dingo.Module
	}
)

func (a *appmodule) Inject(root *config.Area, router *router.Router, logger flamingo.Logger) {
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

func serveProvider(a *appmodule, logger flamingo.Logger) *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Default serve command - starts on Port 3322",
		Run: func(cmd *cobra.Command, args []string) {
			a.router.Init(a.root)
			logger.Info(fmt.Sprintf("Starting HTTP Server at %s .....", addr))
			err := a.server.ListenAndServe()
			if err != nil {
				if err == http.ErrServerClosed {
					logger.Error(err)
				} else {
					logger.Fatal("unexpected error in serving:", err)
				}
			}
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

type option func(config *appconfig)

func ConfigDir(configdir string) option {
	return func(config *appconfig) {
		config.configDir = configdir
	}
}

type appconfig struct {
	configDir string
}

// App is a simple app-runner for flamingo
func App(modules []dingo.Module, options ...option) {
	app := new(appmodule)
	root := config.NewArea("root", modules)

	root.Modules = append([]dingo.Module{
		new(framework.InitModule),
		new(framework.Module),
		new(zap.Module),
		new(cmd.Module),
	}, root.Modules...)

	root.Modules = append(root.Modules, app)
	cfg := &appconfig{
		configDir: "config",
	}
	for _, option := range options {
		option(cfg)
	}
	config.Load(root, cfg.configDir)

	cmd := root.Injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command)
	root.Injector.GetInstance(new(router.EventRouterProvider)).(router.EventRouterProvider)().Dispatch(context.Background(), &flamingo.AppStartupEvent{AppModule: nil})

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
