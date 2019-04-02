package oauth_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(oauth.Module).DefaultConfig(),
	}

	cfgModule.Map["session.backend"] = ""

	if err := dingo.TryModule(cfgModule, new(oauth.Module)); err != nil {
		t.Error(err)
	}
}
