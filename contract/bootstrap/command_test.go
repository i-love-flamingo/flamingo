package bootstrap_test

import (
	"bytes"
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/ops"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
)

type notifyFunc func(ctx context.Context, event flamingoFramework.Event)

func (f notifyFunc) Notify(ctx context.Context, event flamingoFramework.Event) {
	f(ctx, event)
}

func TestCommand_BoundCobraCommandRunsAndShutdownFiresOnce(t *testing.T) { //nolint:paralleltest // runs an app, dingo scopes are process-global
	var (
		shutdowns         atomic.Int32
		greetedBeforeStop atomic.Bool
		out               bytes.Buffer
	)

	modules := []dingo.Module{
		new(ops.Module),
		dingo.ModuleFunc(func(injector *dingo.Injector) {
			flamingoFramework.BindEventSubscriber(injector).ToInstance(
				notifyFunc(func(_ context.Context, event flamingoFramework.Event) {
					if _, ok := event.(*flamingoFramework.ShutdownEvent); ok {
						shutdowns.Add(1)
						greetedBeforeStop.Store(bytes.Contains(out.Bytes(), []byte("hello")))
					}
				}),
			)
		}),
	}

	harness.Isolate(t)

	app, err := flamingo.NewApplication(modules, flamingo.WithArgs("ops-greet"))
	require.NoError(t, err)

	injector, err := app.ConfigArea().GetInitializedInjector()
	require.NoError(t, err)

	instance, err := injector.GetAnnotatedInstance(new(cobra.Command), "flamingo")
	require.NoError(t, err)

	rootCmd, ok := instance.(*cobra.Command)
	require.True(t, ok)

	rootCmd.SetOut(&out)

	require.NoError(t, app.Run())

	assert.Contains(t, out.String(), "hello")
	assert.EqualValues(t, 1, shutdowns.Load())
	assert.Never(t, func() bool { return shutdowns.Load() > 1 }, 50*time.Millisecond, 5*time.Millisecond, "a second shutdown after Run")
	assert.True(t, greetedBeforeStop.Load(), "no shutdown event after the command wrote its output")
}
