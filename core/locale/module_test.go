package locale_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(locale.Module)); err != nil {
		t.Error(err)
	}
}
