package flamingo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type appmodule struct {
	root   *config.Area
	router *web.Router
	server *http.Server
	logger flamingo.Logger
}

// Inject basic application dependencies
func (a *appmodule) Inject(root *config.Area, router *web.Router, logger flamingo.Logger) {
	a.root = root
	a.router = router
	a.logger = logger
	a.server = &http.Server{
		Addr:    ":3322",
		Handler: &ochttp.Handler{IsPublicEndpoint: true, Handler: a.router, StartOptions: trace.StartOptions{Sampler: opencensus.Sampler}},
	}
}

// Configure dependency injection
func (a *appmodule) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).ToInstance(a)

	injector.BindMulti(new(cobra.Command)).ToProvider(func() *cobra.Command {
		return serveProvider(a, a.logger)
	})
}

func serveProvider(a *appmodule, logger flamingo.Logger) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Default serve command - starts on Port 3322",
		Run: func(cmd *cobra.Command, args []string) {
			a.router.Init(a.root)
			logger.Info(fmt.Sprintf("Starting HTTP Server at %s .....", a.server.Addr))
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
	serveCmd.Flags().StringVarP(&a.server.Addr, "addr", "a", ":3322", "addr on which flamingo runs")

	return serveCmd
}

// Notify upon flamingo Shutdown event
func (a *appmodule) Notify(ctx context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		a.logger.Info("Shutdown server on ", a.server.Addr)

		err := a.server.Shutdown(ctx)
		if err != nil {
			a.logger.Error("unexpected error on server shutdown: ", err)
		}
	}
}

type option func(config *appconfig)

// ConfigDir configuration option
func ConfigDir(configdir string) func(config *appconfig) {
	return func(config *appconfig) {
		config.configDir = configdir
	}
}

// ChildAreas allows to define additional config areas for roots
func ChildAreas(areas ...*config.Area) func(config *appconfig) {
	return func(config *appconfig) {
		config.childAreas = areas
	}
}

type appconfig struct {
	configDir  string
	childAreas []*config.Area
}

type eventRouterProvider func() flamingo.EventRouter

// App is a simple app-runner for flamingo
func App(modules []dingo.Module, options ...option) {
	cfg := &appconfig{
		configDir: "config",
	}
	for _, option := range options {
		option(cfg)
	}

	app := new(appmodule)
	root := config.NewArea("root", modules, cfg.childAreas...)

	root.Modules = append([]dingo.Module{
		new(framework.InitModule),
		new(zap.Module),
		new(cmd.Module),
	}, root.Modules...)

	root.Modules = append(root.Modules, app)
	config.Load(root, cfg.configDir)

	rootCmd := root.Injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command)
	root.Injector.GetInstance(new(eventRouterProvider)).(eventRouterProvider)().Dispatch(context.Background(), new(flamingo.StartupEvent))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
