package auth_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/contract/internal/testapp/oidcfake"
)

func TestSecurity_HandleIfLoggedIn(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	server := oidcfake.Start(t)
	app := bootAuth(t, server.BrokerYAML("shop"), nil)

	// the middleware redirects to the route "auth.login", which the account module registers
	target := location(t, app.Get(t, "/en/account"), http.StatusSeeOther)

	assert.Equal(t, "/en/account/login", target.Path)
	assert.Equal(t, app.URL+"/en/account", target.Query().Get("redirecturl"))

	login(t, app, server)

	resp := app.Get(t, "/en/account")

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.JSONEq(t, `{"ok":true}`, resp.Body)
}
