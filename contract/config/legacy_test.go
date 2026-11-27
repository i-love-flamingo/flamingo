package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/ops"
)

type greetingValues struct {
	Greeting string `inject:"config:ops.greeting"`
}

func TestLegacyAlias_MapsOldKey(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, []dingo.Module{new(ops.Module)}, harness.Config{
		Files: map[string]string{"config.yml": "opsGreeting: legacy\n"},
	})

	assert.Equal(t, "legacy", rootValues[greetingValues](t, app).Greeting)
}

type oldGreetingValues struct {
	Greeting string `inject:"config:opsGreeting"`
}

// Some code still injects the old key.
func TestLegacyAlias_CopiesNewKeyToOldKey(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, []dingo.Module{new(ops.Module)}, harness.Config{
		Files: map[string]string{"config.yml": "ops:\n  greeting: current\n"},
	})

	assert.Equal(t, "current", rootValues[oldGreetingValues](t, app).Greeting)
}
