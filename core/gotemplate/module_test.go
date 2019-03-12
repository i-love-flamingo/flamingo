package gotemplate_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/gotemplate"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(gotemplate.Module)); err != nil {
		t.Error(err)
	}
}
