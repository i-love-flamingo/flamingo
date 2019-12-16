package security_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/security"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(security.Module)); err != nil {
		t.Error(err)
	}
}
