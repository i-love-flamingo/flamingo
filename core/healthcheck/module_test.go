package healthcheck_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/framework/config"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(healthcheck.Module).DefaultConfig(),
	}

	cfgModule.Map["session.backend"] = ""

	if err := dingo.TryModule(cfgModule, new(healthcheck.Module)); err != nil {
		t.Error(err)
	}
}
