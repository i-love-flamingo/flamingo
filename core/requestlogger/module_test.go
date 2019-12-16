package requestlogger_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(requestlogger.Module)); err != nil {
		t.Error(err)
	}
}
