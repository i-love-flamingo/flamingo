package canonicalurl_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/canonicalurl"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(canonicalurl.Module)); err != nil {
		t.Error(err)
	}
}
