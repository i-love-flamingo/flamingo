package flamingo_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"flamingo.me/dingo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/framework/cmd"
	framework "flamingo.me/flamingo/v3/framework/flamingo"
)

type notifyFunc func(ctx context.Context, event framework.Event)

func (nf notifyFunc) Notify(ctx context.Context, event framework.Event) {
	nf(ctx, event)
}

func TestCmdEventsTriggeredProperly(t *testing.T) { //nolint:paralleltest // due to dingo.Singleton
	assertStartupAndShutdownOnce := func(t *testing.T, startupEventCount, shutdownEventCount int32) {
		t.Helper()

		assert.Equal(t, int32(1), startupEventCount, "startupEventCount should be 1")
		assert.Equal(t, int32(1), shutdownEventCount, "shutdownEventCount should be 1")
	}

	tests := []struct {
		name string
		args string
		want func(t *testing.T, startupEventCount, shutdownEventCount int32)
	}{
		{
			name: "command with Run: both startup and shutdown events are triggered",
			args: "test_cmd_run",
			want: assertStartupAndShutdownOnce,
		},
		{
			name: "command with RunE: both startup and shutdown events are triggered",
			args: "test_cmd_run_e",
			want: assertStartupAndShutdownOnce,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // due to dingo.Singleton
		t.Run(tt.name, func(t *testing.T) {
			var (
				startupEventCount  atomic.Int32
				shutdownEventCount atomic.Int32
			)

			modules := []dingo.Module{
				dingo.ModuleFunc(func(injector *dingo.Injector) {
					injector.BindMulti(new(cobra.Command)).ToInstance(&cobra.Command{
						Use:   "test_cmd_run",
						Short: "test_cmd_run",
						Run: func(cmd *cobra.Command, args []string) {
							t.Log("test_cmd_run executed")
							assert.Equal(t, "test_cmd_run", cmd.Short)
							assert.Empty(t, args)
						},
					})
					injector.BindMulti(new(cobra.Command)).ToInstance(&cobra.Command{
						Use:   "test_cmd_run_e",
						Short: "test_cmd_run_e",
						RunE: func(cmd *cobra.Command, args []string) error {
							t.Log("test_cmd_run_e executed")
							assert.Equal(t, "test_cmd_run_e", cmd.Short)
							assert.Empty(t, args)

							return nil
						},
					})
				}),
				dingo.ModuleFunc(func(injector *dingo.Injector) {
					framework.BindEventSubscriber(injector).ToInstance(notifyFunc(func(ctx context.Context, event framework.Event) {
						switch event.(type) {
						case *framework.StartupEvent:
							startupEventCount.Add(1)
						case *framework.ShutdownEvent:
							shutdownEventCount.Add(1)
						}
					}))
				}),
			}

			dingo.Singleton = dingo.NewSingletonScope()
			dingo.ChildSingleton = dingo.NewChildSingletonScope()

			app, err := flamingo.NewApplication(modules,
				flamingo.WithArgs(tt.args),
			)
			require.NoError(t, err)

			err = app.Run()
			require.NoError(t, err)

			tt.want(t, startupEventCount.Load(), shutdownEventCount.Load())
		})
	}
}

func buildSignalSender(t *testing.T) func() {
	t.Helper()

	return func() {
		pid := os.Getpid()
		process, err := os.FindProcess(pid)
		require.NoError(t, err)

		err = process.Signal(os.Interrupt)
		require.NoError(t, err)
	}
}

func TestGracefulShutdown(t *testing.T) { //nolint:paralleltest // due to dingo.Singleton
	assertShutdownOnce := func(t *testing.T, shutdownEventCount int32) {
		t.Helper()

		assert.Equal(t, int32(1), shutdownEventCount, "shutdownEventCount should be 1")
	}

	tests := []struct {
		name                string
		args                string
		onServerStartup     func()
		onShutdown          func()
		insideCommandRun    func(cmd *cobra.Command, args []string)
		insideCommandRunE   func(cmd *cobra.Command, args []string) error
		wantErr             assert.ErrorAssertionFunc
		assertShutdownCount func(t *testing.T, shutdownEventCount int32)
	}{
		{
			name:                "serve command interrupted by SIGINT triggers graceful shutdown",
			args:                "serve",
			onServerStartup:     sync.OnceFunc(buildSignalSender(t)),
			onShutdown:          func() {},
			wantErr:             assert.NoError,
			assertShutdownCount: assertShutdownOnce,
		},
		{
			name: "custom Run command interrupted by SIGINT triggers graceful shutdown",
			args: "test_cmd_run",
			insideCommandRun: func(cmd *cobra.Command, args []string) {
				buildSignalSender(t)()
			},
			onShutdown:          func() {},
			wantErr:             assert.NoError,
			assertShutdownCount: assertShutdownOnce,
		},
		{
			name: "custom RunE command interrupted by SIGINT triggers graceful shutdown",
			args: "test_cmd_run_e",
			insideCommandRunE: func(cmd *cobra.Command, args []string) error {
				buildSignalSender(t)()

				return nil
			},
			onShutdown:          func() {},
			wantErr:             assert.NoError,
			assertShutdownCount: assertShutdownOnce,
		},
		{
			name: "graceful shutdown interrupted by SIGINT forces hard shutdown",
			args: "test_cmd_run",
			insideCommandRun: func(cmd *cobra.Command, args []string) {
				send := buildSignalSender(t)
				send()
				time.Sleep(time.Millisecond)
				send()
			},
			onShutdown: sync.OnceFunc(func() {
				// artificial delay, so that second interrupt could arrive
				time.Sleep(time.Second)
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, cmd.ErrGracefulShutdown)
			},
			assertShutdownCount: assertShutdownOnce,
		},

		// TODO: use testing/synctest for this once it is not experimental
		{
			name: "graceful shutdown timed out forces hard shutdown",
			args: "test_cmd_run",
			insideCommandRun: func(cmd *cobra.Command, args []string) {
				send := buildSignalSender(t)
				send()
			},
			onShutdown: sync.OnceFunc(func() {
				time.Sleep(35 * time.Second)
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, cmd.ErrGracefulShutdown)
			},
			assertShutdownCount: assertShutdownOnce,
		},
	}

	for _, tt := range tests { //nolint:paralleltest // due to dingo.Singleton
		t.Run(tt.name, func(t *testing.T) {
			var shutdownEventCount atomic.Int32

			modules := []dingo.Module{
				dingo.ModuleFunc(func(injector *dingo.Injector) {
					injector.BindMulti(new(cobra.Command)).ToInstance(&cobra.Command{
						Use:   "test_cmd_run",
						Short: "test_cmd_run",
						Run:   tt.insideCommandRun,
					})
					injector.BindMulti(new(cobra.Command)).ToInstance(&cobra.Command{
						Use:   "test_cmd_run_e",
						Short: "test_cmd_run_e",
						RunE:  tt.insideCommandRunE,
					})
				}),
				dingo.ModuleFunc(func(injector *dingo.Injector) {
					framework.
						BindEventSubscriber(injector).
						ToInstance(
							notifyFunc(func(ctx context.Context, event framework.Event) {
								switch ev := event.(type) {
								case *framework.ServerStartEvent:
									if _, err := net.Dial("tcp", fmt.Sprintf(":%s", ev.Port)); err != nil {
										t.Fatalf("failed to connect to server")
									}

									tt.onServerStartup()
								case *framework.ShutdownEvent:
									shutdownEventCount.Add(1)
									tt.onShutdown()
								}
							}))
				}),
			}

			dingo.Singleton = dingo.NewSingletonScope()
			dingo.ChildSingleton = dingo.NewChildSingletonScope()

			app, err := flamingo.NewApplication(modules, flamingo.WithArgs(tt.args))
			require.NoError(t, err)

			err = app.Run()
			tt.wantErr(t, err)
			tt.assertShutdownCount(t, shutdownEventCount.Load())
		})
	}
}
