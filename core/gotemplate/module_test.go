package gotemplate_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/gotemplate"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(gotemplate.Module)); err != nil {
		t.Error(err)
	}
}
