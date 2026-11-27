package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/config"
)

type tagSliceValues struct {
	Tags config.Slice `inject:"config:brand.tags"`
}

func TestSlice_YAMLListInjectsAsConfigSlice(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, nil, harness.Config{
		Files: map[string]string{"config.yml": "brand:\n  tags: [a, b]\n"},
	})

	assert.Equal(t, tagSliceValues{Tags: config.Slice{"a", "b"}}, rootValues[tagSliceValues](t, app))
}
