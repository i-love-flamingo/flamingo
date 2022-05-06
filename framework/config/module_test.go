package config_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(config.Module)); err != nil {
		t.Error(err)
	}
}
