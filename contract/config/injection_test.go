package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/config"
)

type brandMapValues struct {
	Brand config.Map `inject:"config:brand"`
}

type brandSettings struct {
	Name   string
	Size   int
	Nested struct{ Flag bool }
}

// Numbers are float64 and nested maps are config.Map.
func TestConfigMap_InjectsSubTreeWithNestedValues(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, nil, harness.Config{
		Files: map[string]string{
			"config.yml": "brand:\n  name: shop\n  size: 3\n  nested:\n    flag: true\n",
		},
	})

	got := rootValues[brandMapValues](t, app).Brand

	assert.Equal(t, "shop", got["name"])
	assert.Equal(t, float64(3), got["size"])
	assert.Equal(t, config.Map{"flag": true}, got["nested"])

	var settings brandSettings

	require.NoError(t, got.MapInto(&settings))
	assert.Equal(t, "shop", settings.Name)
	assert.Equal(t, 3, settings.Size)
	assert.True(t, settings.Nested.Flag)
}

// injectModule gets T injected, as a module of an app does with its config, and keeps the values.
type injectModule[T any] struct {
	out *T
}

// Configure does nothing, the module only needs the injection.
func (*injectModule[T]) Configure(*dingo.Injector) {}

// Inject the config values.
func (m *injectModule[T]) Inject(values *T) {
	*m.out = *values
}

type optionalValues struct {
	Name    string `inject:"config:brand.name,optional"`
	Missing string `inject:"config:no.such.key,optional"`
	Count   int    `inject:"config:no.such.count,optional"`
}

type requiredValues struct {
	Missing string `inject:"config:no.such.key"`
}

func TestConfigInjection_OptionalKeyIsZeroAndRequiredKeyFailsStart(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	files := map[string]string{"config.yml": "brand:\n  name: shop\n"}

	t.Run("optional", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		var got optionalValues

		harness.Boot(t, []dingo.Module{&injectModule[optionalValues]{out: &got}}, harness.Config{Files: files})

		assert.Equal(t, optionalValues{Name: "shop"}, got)
	})

	t.Run("required", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		var got requiredValues

		app, err := harness.TryBoot(t, []dingo.Module{&injectModule[requiredValues]{out: &got}}, harness.Config{Files: files})

		require.Error(t, err)
		assert.Nil(t, app)
		assert.Contains(t, err.Error(), "no.such.key")
	})
}
