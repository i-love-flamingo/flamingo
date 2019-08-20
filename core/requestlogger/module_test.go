package requestlogger_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/framework/config"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(requestlogger.Module).DefaultConfig(),
	}
	if err := dingo.TryModule(cfgModule, new(requestlogger.Module)); err != nil {
		t.Error(err)
	}
}
