package locale_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/locale"
	"testing"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(locale.Module)); err != nil {
		t.Error(err)
	}
}
