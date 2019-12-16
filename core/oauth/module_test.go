package oauth_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{
		"flamingo.session.backend": "memory",
		"core.oauth.secret":        "secret",
		"core.oauth.server":        "server",
		"core.oauth.clientid":      "clientid",
	}, new(oauth.Module)); err != nil {
		t.Error(err)
	}
}
