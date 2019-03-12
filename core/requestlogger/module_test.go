package requestlogger_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(requestlogger.Module)); err != nil {
		t.Error(err)
	}
}
