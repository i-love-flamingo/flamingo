package cmd_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/cmd"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(cmd.Module)); err != nil {
		t.Error(err)
	}
}
