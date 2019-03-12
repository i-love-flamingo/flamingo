package internalauth_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/internalauth"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(internalauth.InternalAuth)); err != nil {
		t.Error(err)
	}
}
