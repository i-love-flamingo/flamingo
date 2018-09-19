package redirects

import (
	"testing"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
)

func TestModule_Configure(t *testing.T) {
	module := new(Module)
	module.UseInPrefixRouter = true
	module.UseInRouter = true

	cfgModule := &config.Module{
		Map: config.Map{
			"redirects.useInRouter":       true,
			"redirects.useInPrefixRouter": true,
		},
	}

	if err := dingo.TryModule(cfgModule, module); err != nil {
		t.Error(err)
	}
}
