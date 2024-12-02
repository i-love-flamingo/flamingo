package framework_test

import (
	"testing"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(framework.InitModule)); err != nil {
		t.Error(err)
	}
}
