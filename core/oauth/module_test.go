package auth_test

import (
	"testing"

	"flamingo.me/dingo"
	auth "flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(auth.Module).DefaultConfig(),
	}

	cfgModule.Map["session.backend"] = ""

	if err := dingo.TryModule(cfgModule, new(auth.Module)); err != nil {
		t.Error(err)
	}
}
