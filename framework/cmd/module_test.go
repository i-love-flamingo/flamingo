package cmd_test

import (
	"testing"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/cmd"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(cmd.Module)); err != nil {
		t.Error(err)
	}
}
