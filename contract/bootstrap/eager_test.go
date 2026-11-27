package bootstrap_test

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/contract/internal/harness"
)

type eagerThing struct{}

// Eager singletons are off by default: AsEagerSingleton alone builds nothing at start.
func TestEagerSingletons_BuiltOnlyWithOption(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	for _, tc := range []struct { //nolint:paralleltest // boots an app, dingo scopes are process-global
		name    string
		options []flamingo.ApplicationOption
		builds  int32
	}{
		{name: "default", builds: 0},
		{name: "enabled", options: []flamingo.ApplicationOption{flamingo.SetEagerSingletons(true)}, builds: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var built atomic.Int32

			harness.Boot(t,
				[]dingo.Module{dingo.ModuleFunc(func(injector *dingo.Injector) {
					injector.Bind(new(eagerThing)).ToProvider(func() *eagerThing {
						built.Add(1)

						return new(eagerThing)
					}).AsEagerSingleton()
				})},
				harness.Config{Options: tc.options},
			)

			assert.Equal(t, tc.builds, built.Load())
		})
	}
}
