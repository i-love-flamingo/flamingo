package canonicalUrl

import (
	"testing"

	"flamingo.me/flamingo/framework/dingo"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(Module)); err != nil {
		t.Error(err)
	}
}
