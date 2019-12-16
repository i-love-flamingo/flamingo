package requesttask_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/requesttask"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(requesttask.Module)); err != nil {
		t.Error(err)
	}
}
