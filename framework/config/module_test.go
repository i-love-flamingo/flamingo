package config_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(config.Module)); err != nil {
		t.Error(err)
	}
}
