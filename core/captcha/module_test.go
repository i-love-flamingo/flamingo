package captcha_test

import (
	"testing"

	"flamingo.me/flamingo/core/captcha"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
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
