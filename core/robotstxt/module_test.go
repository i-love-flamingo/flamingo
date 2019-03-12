package robotstxt_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/robotstxt"
	"flamingo.me/flamingo/v3/framework/config"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(robotstxt.Module).DefaultConfig(),
	}

	if err := dingo.TryModule(cfgModule, new(robotstxt.Module)); err != nil {
		t.Error(err)
	}
}
