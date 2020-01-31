package fake

import (
	"testing"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule(t *testing.T) {
	if err := config.TryModules(config.Map{"flamingo.debug.mode": true}, new(auth.WebModule)); err != nil {
		t.Error(err)
	}
}
