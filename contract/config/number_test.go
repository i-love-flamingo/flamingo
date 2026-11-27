package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/contract/internal/harness"
)

type countValues struct {
	AsInt     int     `inject:"config:brand.count"`
	AsInt64   int64   `inject:"config:brand.count"`
	AsFloat64 float64 `inject:"config:brand.count"`
}

type weightValues struct {
	AsInt     int     `inject:"config:brand.weight"`
	AsInt64   int64   `inject:"config:brand.weight"`
	AsFloat64 float64 `inject:"config:brand.weight"`
}

func TestNumber_WholeNumberInjectsIntoIntInt64Float64(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, nil, harness.Config{
		Files: map[string]string{
			"config.yml": "brand:\n  count: 3\n",
			"config.cue": "brand: weight: 3\n",
		},
	})

	assert.Equal(t, countValues{AsInt: 3, AsInt64: 3, AsFloat64: 3}, rootValues[countValues](t, app), "YAML")
	assert.Equal(t, weightValues{AsInt: 3, AsInt64: 3, AsFloat64: 3}, rootValues[weightValues](t, app), "CUE")
}
