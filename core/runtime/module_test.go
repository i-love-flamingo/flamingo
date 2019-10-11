package runtime_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/runtime"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(zap.Module).DefaultConfig(),
	}

	cfgModule.Map["area"] = ""

	if err := dingo.TryModule(cfgModule, new(zap.Module), new(runtime.Module)); err != nil {
		t.Error(err)
	}
}
