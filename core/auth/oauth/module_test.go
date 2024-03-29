package oauth_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{"core.auth.web.debugController": false}, new(oauth.Module)); err != nil {
		t.Error(err)
	}
}
