package canonicalurl

import (
	"testing"

	"flamingo.me/dingo"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(Module)); err != nil {
		t.Error(err)
	}
}
