package cmd

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"flamingo.me/flamingo/core/cmd/interfaces/command"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"github.com/spf13/cobra"
)

// Module for DI
type Module struct{}

var dingoTrace *bool

var once = sync.Once{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	// command.VersionCmd,
	// command.DiCmd,
	// command.RoutingConfCmd,
	// command.RouterCmd,
	// command.DataControllerCmd,
	// command.TplfuncsCmd,

	injector.Bind(new(cobra.Command)).AnnotatedWith("flamingo").ToProvider(
		func(
			commands []*cobra.Command,
			eventRouterProvider router.EventRouterProvider,
			logger flamingo.Logger,
			config *struct {
				Name string `inject:"config:cmd.name"`
			}) *cobra.Command {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

			eventRouterProvider().Dispatch(context.Background(), &flamingo.AppStartupEvent{AppModule: m})
			once.Do(func() {
				go shutdown(eventRouterProvider(), signals, logger, m)
			})

			rootCmd := &cobra.Command{
				Use:              config.Name,
				Short:            "Flamingo " + config.Name,
				TraverseChildren: true,
			}
			rootCmd.FParseErrWhitelist.UnknownFlags = true

			rootCmd.AddCommand(commands...)

			return rootCmd
		},
	)

	injector.BindMulti(new(cobra.Command)).ToProvider(command.ConfigCmd)
}

// DefaultConfig specifies the command name
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cmd.name": filepath.Base(os.Args[0]),
	}
}

func shutdown(eventRouter event.Router, signals <-chan os.Signal, logger flamingo.Logger, m dingo.Module) {
	<-signals
	logger.Info("start graceful shutdown")

	stopper := make(chan struct{})

	go func() {
		eventRouter.Dispatch(context.Background(), &flamingo.AppShutdownEvent{AppModule: m})
		close(stopper)
	}()

	select {
	case <-signals:
		logger.Info("second interrupt signal received, hard shutdown")
		os.Exit(130)
	case <-time.After(30 * time.Second):
		logger.Info("time limit reached, hard shutdown")
		os.Exit(130)
	case <-stopper:
		logger.Info("graceful shutdown complete")
		os.Exit(0)
	}
}

// Run the root command
func Run(injector *dingo.Injector) error {
	return injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command).Execute()
}
