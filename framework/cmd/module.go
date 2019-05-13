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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// Module for DI
type Module struct{}

var once = sync.Once{}

type (
	eventRouterProvider func() flamingo.EventRouter
	flagSetProvider     func() []*pflag.FlagSet
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(cobra.Command)).AnnotatedWith("flamingo").ToProvider(
		func(
			commands []*cobra.Command,
			eventRouterProvider eventRouterProvider,
			logger flamingo.Logger,
			flagSetProvider flagSetProvider,
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
			for _, set := range flagSetProvider() {
				rootCmd.PersistentFlags().AddFlagSet(set)
			}

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
	injector.GetInstance(new(eventRouterProvider)).(eventRouterProvider)().Dispatch(context.Background(), &flamingo.StartupEvent{})

	return cmd.Execute()
}
