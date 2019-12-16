package runtime_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/runtime"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(zap.Module), new(runtime.Module)); err != nil {
		t.Error(err)
	}
}
