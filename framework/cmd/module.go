package cmd

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
)

// Module for DI
type Module struct{}

var once = sync.Once{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(cobra.Command)).AnnotatedWith("flamingo").ToProvider(
		func(
			commands []*cobra.Command,
			eventRouterProvider web.EventRouterProvider,
			logger flamingo.Logger,
			config *struct {
				Name string `inject:"config:cmd.name"`
			}) *cobra.Command {
			signals := make(chan os.Signal, 1)
			shutdownComplete := make(chan struct{}, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

			once.Do(func() {
				go shutdown(eventRouterProvider(), signals, shutdownComplete, logger)
			})

			rootCmd := &cobra.Command{
				Use:              config.Name,
				Short:            "Flamingo " + config.Name,
				TraverseChildren: true,
				PersistentPostRun: func(cmd *cobra.Command, args []string) {
					signals <- syscall.SIGTERM
					<-shutdownComplete
				},
			}
			rootCmd.FParseErrWhitelist.UnknownFlags = true

			rootCmd.AddCommand(commands...)

			return rootCmd
		},
	)
}

// DefaultConfig specifies the command name
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cmd.name": filepath.Base(os.Args[0]),
	}
}

func shutdown(eventRouter flamingo.EventRouter, signals <-chan os.Signal, complete chan<- struct{}, logger flamingo.Logger) {
	<-signals
	logger.Info("start graceful shutdown")

	stopper := make(chan struct{})

	go func() {
		eventRouter.Dispatch(context.Background(), &flamingo.ShutdownEvent{})
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
		complete <- struct{}{}
		os.Exit(0)
	}
}

// Run the root command
func Run(injector *dingo.Injector) error {
	cmd := injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command)
	injector.GetInstance(new(web.EventRouterProvider)).(web.EventRouterProvider)().Dispatch(context.Background(), &flamingo.StartupEvent{})

	return cmd.Execute()
}
