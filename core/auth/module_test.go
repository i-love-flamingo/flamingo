package auth

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule(t *testing.T) {
	if err := config.TryModules(nil, new(WebModule)); err != nil {
		t.Error(err)
	}
}
