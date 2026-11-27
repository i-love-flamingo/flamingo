package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
)

type brandSizeValues struct {
	Size int `inject:"config:brand.size"`
}

// sizeModule gives the default of brand.size with a CUE schema.
type sizeModule struct{}

func (*sizeModule) Configure(*dingo.Injector) {}

func (*sizeModule) CueConfig() string {
	return `brand: { size: int | *3 }`
}

func TestCueConfig_ModuleDefaultsApply(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, []dingo.Module{new(sizeModule)}, harness.Config{})

	assert.Equal(t, 3, rootValues[brandSizeValues](t, app).Size)
}
