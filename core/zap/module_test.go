package zap_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(zap.Module).DefaultConfig(),
	}

	cfgModule.Map["area"] = ""

	if err := dingo.TryModule(cfgModule, new(zap.Module)); err != nil {
		t.Error(err)
	}
}
