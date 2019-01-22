package captcha_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/captcha"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
)

func TestModule_Configure(t *testing.T) {
	module := new(captcha.Module)

	cfgModule := &config.Module{
		Map: module.DefaultConfig(),
	}

	if err := dingo.TryModule(
		cfgModule,
		module,
	); err != nil {
		t.Error(err)
	}
}
