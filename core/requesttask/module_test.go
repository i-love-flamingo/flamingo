package requesttask_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/requesttask"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(requesttask.Module)); err != nil {
		t.Error(err)
	}
}
