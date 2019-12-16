package internalauth_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/internalauth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(internalauth.InternalAuth)); err != nil {
		t.Error(err)
	}
}
