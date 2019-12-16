package robotstxt_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/robotstxt"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(robotstxt.Module)); err != nil {
		t.Error(err)
	}
}
