package config_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type areaTracingGuard struct{ calls int }

type areaTracingCycle struct{}

type areaTracingCycleProvider func() *areaTracingCycle

func (*areaTracingCycle) Inject(next areaTracingCycleProvider, guard *areaTracingGuard) {
	guard.calls++
	if guard.calls > 4 {
		panic("recursion guard")
	}

	next()
}

func assertAreaCircularTracing(t *testing.T, injector *dingo.Injector, guard *areaTracingGuard, enabled bool) {
	t.Helper()

	guard.calls = 0

	want := "recursion guard"
	if enabled {
		want = "detected circular dependency"
	}

	require.PanicsWithValue(t, want, func() {
		_, _ = injector.GetInstance(new(areaTracingCycle))
	})

	if enabled {
		assert.LessOrEqual(t, guard.calls, 4, "the circular diagnostic must precede the recursion guard")
	}
}

func TestAreaCircularTracingBeforeModuleInitialization(t *testing.T) {
	t.Parallel()

	for _, parentTracing := range []bool{false, true} {
		t.Run(fmt.Sprintf("parent=%t", parentTracing), func(t *testing.T) {
			t.Parallel()

			guards := make(map[*config.Area]*areaTracingGuard)
			newArea := func(name string, enabled bool, children ...*config.Area) *config.Area {
				guard := new(areaTracingGuard)
				area := config.NewArea(name, []dingo.Module{dingo.ModuleFunc(func(injector *dingo.Injector) {
					injector.Bind(new(areaTracingGuard)).ToInstance(guard)
					assertAreaCircularTracing(t, injector, guard, enabled)
				})}, children...)
				guards[area] = guard

				return area
			}

			child := newArea("child", true)
			sibling := newArea("sibling", parentTracing)

			root := newArea("root", parentTracing, child, sibling)
			if parentTracing {
				root.SetDingoTracing(true, false)
			} else {
				child.SetDingoTracing(true, false)
			}

			for _, area := range []*config.Area{child, root, sibling} {
				injector, err := area.GetInitializedInjector()
				require.NoError(t, err)
				assertAreaCircularTracing(t, injector, guards[area], area == child || parentTracing)
			}
		})
	}
}

type areaTracingProbe struct{}

//nolint:paralleltest // Captures the process-wide slog logger.
func TestAreaInjectionTracingBeforeModuleInitialization(t *testing.T) {
	var output bytes.Buffer

	previousLogger := slog.Default()

	slog.SetDefault(slog.New(slog.NewTextHandler(&output, nil)))
	t.Cleanup(func() { slog.SetDefault(previousLogger) })

	for _, parentTracing := range []bool{false, true} { //nolint:paralleltest // Shares the process-wide slog logger.
		t.Run(fmt.Sprintf("parent=%t", parentTracing), func(t *testing.T) {
			assertTracing := func(injector *dingo.Injector, enabled bool) {
				output.Reset()

				_, err := injector.GetInstance(new(areaTracingProbe))
				require.NoError(t, err)
				assert.Equal(t, enabled, bytes.Contains(output.Bytes(), []byte("INJECTING: flamingo.me/flamingo/v3/framework/config_test#areaTracingProbe")))
			}
			newArea := func(name string, enabled bool, children ...*config.Area) *config.Area {
				return config.NewArea(name, []dingo.Module{dingo.ModuleFunc(func(injector *dingo.Injector) {
					assertTracing(injector, enabled)
				})}, children...)
			}

			child := newArea("child", true)
			sibling := newArea("sibling", parentTracing)

			root := newArea("root", parentTracing, child, sibling)
			if parentTracing {
				root.SetDingoTracing(false, true)
			} else {
				child.SetDingoTracing(false, true)
			}

			for _, area := range []*config.Area{child, root, sibling} {
				injector, err := area.GetInitializedInjector()
				require.NoError(t, err)
				assertTracing(injector, area == child || parentTracing)
			}
		})
	}
}
