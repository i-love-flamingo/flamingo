package opencensus_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(opencensus.Module).DefaultConfig(),
	}

	if err := dingo.TryModule(cfgModule, new(opencensus.Module)); err != nil {
		t.Error(err)
	}
}
