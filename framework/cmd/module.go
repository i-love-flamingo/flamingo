package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	eventRouterProvider func() flamingo.EventRouter
	flagSetProvider     func() []*pflag.FlagSet

	// Module for DI
	Module struct{}
)

var (
	// ErrCmdRun is returned when an error occurs during CLI command execution
	ErrCmdRun = errors.New("command execution error")

	// ErrGracefulShutdown is returned when graceful shutdown of the application cannot finish cleanly due to timeout or another signal
	ErrGracefulShutdown = errors.New("graceful shutdown error")

	signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(cobra.Command)).
		AnnotatedWith("flamingo").
		ToProvider(RootCommandProvider).
		In(dingo.Singleton)
}

// RootCommandProvider configures root cobra command to be used by the framework
func RootCommandProvider(
	commands []*cobra.Command,
	eventRouterProvider eventRouterProvider,
	logger flamingo.Logger,
	flagSetProvider flagSetProvider,
	config *struct {
		Name string `inject:"config:flamingo.cmd.name"`
	},
) *cobra.Command {
	var (
		// declaring global context.CancelFunc and error to be used from both PersistentPreRun and PersistentPostRunE
		stop             context.CancelFunc
		err              error
		execShutdownOnce = sync.OnceFunc(func() {
			err = shutdown(logger, eventRouterProvider())
		})
	)

	rootCmd := &cobra.Command{
		Use:              config.Name,
		Short:            "Flamingo " + config.Name,
		TraverseChildren: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var ctx context.Context
			ctx, stop = signal.NotifyContext(cmd.Context(), signals...)
			cmd.SetContext(ctx)

			go func() {
				// if in the serve command wait for signal to come (context will be cancelled),
				// then disable listening for signals (calling stop())
				// then execute shutdown func
				<-cmd.Context().Done()

				if stop != nil {
					stop()
				}

				execShutdownOnce()
			}()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// on finished command execution
			// stop listening for signals by calling stop()
			// wait for context to be cancelled (should happen immediately after stop())
			// execute shutdown func
			if stop != nil {
				stop()
			}

			<-cmd.Context().Done()

			execShutdownOnce()

			return err
		},
		Example: `Run with -h or -help to see global debug flags`,
	}

	rootCmd.FParseErrWhitelist.UnknownFlags = true
	for _, set := range flagSetProvider() {
		rootCmd.PersistentFlags().AddFlagSet(set)
	}

	rootCmd.AddCommand(commands...)

	return rootCmd
}

// CueConfig specifies the command name
func (*Module) CueConfig() string {
	return fmt.Sprintf(`flamingo: cmd: name: string | *"%s"`, filepath.Base(os.Args[0]))
}

// FlamingoLegacyConfigAlias maps legacy config to new
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"cmd.name": "flamingo.cmd.name"}
}

// shutdown wait for context ctx to be done and dispatches shutdown event
func shutdown(logger flamingo.Logger, eventRouter flamingo.EventRouter) error {
	logger.Info("start graceful shutdown")

	var (
		group   *errgroup.Group
		sigch   = make(chan os.Signal, 1)
		stopper = make(chan struct{})
		timeout = 30 * time.Second
	)

	signal.Notify(sigch, signals...)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	group, ctx = errgroup.WithContext(ctx)

	group.Go(func() error {
		defer close(stopper)

		eventRouter.Dispatch(ctx, &flamingo.ShutdownEvent{})

		return nil
	})

	group.Go(func() error {
		select {
		case <-sigch:
			logger.Info("second interrupt signal received, hard shutdown")
			return fmt.Errorf("%w: signal received", ErrGracefulShutdown)
		case <-ctx.Done():
			logger.Info("time limit reached, hard shutdown")
			return fmt.Errorf("%w: timed out", ErrGracefulShutdown)
		case <-stopper:
			logger.Info("graceful shutdown complete")
			return nil
		}
	})

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("flamingo shutdown: %w", err)
	}

	return nil
}

// Run the root command
func Run(injector *dingo.Injector) error {
	i, err := injector.GetAnnotatedInstance(new(cobra.Command), "flamingo")
	if err != nil {
		return err
	}

	cmd, ok := i.(*cobra.Command)
	if !ok {
		return fmt.Errorf("%w: resolved instance does not have type *cobra.Command", ErrCmdRun)
	}

	i, err = injector.GetInstance(new(eventRouterProvider))
	if err != nil {
		return err
	}

	erp, ok := i.(eventRouterProvider)
	if !ok {
		return fmt.Errorf("%w: resolved instance does not have type eventRouterProvider", ErrCmdRun)
	}

	erp().Dispatch(cmd.Context(), &flamingo.StartupEvent{})

	return cmd.Execute()
}
