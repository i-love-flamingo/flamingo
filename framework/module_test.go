package framework_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(framework.InitModule)); err != nil {
		t.Error(err)
	}
}
