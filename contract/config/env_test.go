package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/contract/internal/harness"
)

type brandNameValues struct {
	Name string `inject:"config:brand.name"`
}

type brandFlagValues struct {
	Flag string `inject:"config:brand.flag"`
}

func TestEnv_PlaceholderResolves(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := harness.Boot(t, nil, harness.Config{
		Env:   map[string]string{"CONTRACT_BRAND": "shop"},
		Files: map[string]string{"config.yml": "brand:\n  name: \"%%ENV:CONTRACT_BRAND%%\"\n"},
	})

	assert.Equal(t, "shop", rootValues[brandNameValues](t, app).Name)
}

func TestEnv_OSEnvDrivesCueIf(t *testing.T) {
	files := map[string]string{
		"config.cue": "brand: flag: *\"off\" | string\n" +
			"if flamingo.os.env.CONTRACT_FLAG != _|_ {\n\tbrand: flag: flamingo.os.env.CONTRACT_FLAG\n}\n",
	}

	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{name: "with the variable", env: map[string]string{"CONTRACT_FLAG": "on"}, want: "on"},
		{name: "without the variable", want: "off"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env == nil {
				// t.Setenv restores the original value at the end of the test.
				t.Setenv("CONTRACT_FLAG", "")
				require.NoError(t, os.Unsetenv("CONTRACT_FLAG"))
			}

			app := harness.Boot(t, nil, harness.Config{
				Env:   tt.env,
				Files: files,
			})

			assert.Equal(t, tt.want, rootValues[brandFlagValues](t, app).Flag)
		})
	}
}
