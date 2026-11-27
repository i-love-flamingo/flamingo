package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/contract/internal/testapp/oidcfake"
)

func TestOIDC_LoginCallbackLogout(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	server := oidcfake.Start(t)
	app := bootAuth(t, server.BrokerYAML("shop"), nil)

	require.JSONEq(t, `{"loggedIn":false}`, loggedIn(t, app))

	login(t, app, server)

	assert.JSONEq(t, `{"loggedIn":true}`, loggedIn(t, app))
	assert.Zero(t, server.EndSessionCalls())

	logout(t, app, server)

	assert.Equal(t, 1, server.EndSessionCalls())
	assert.JSONEq(t, `{"loggedIn":false}`, loggedIn(t, app))
}

// shopperClaims is the answer of the route "account.claims" for the user of login.
const shopperClaims = `{"sub":"shopper","userID":"shopper","email":"shopper@example.com",` +
	`"tokenSubject":"shopper","hasAccessToken":true}`

// The config maps claims to other names, so a claim read under the wrong key is empty.
func TestOIDC_IdentifyAsExposesClaimsAndToken(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	server := oidcfake.Start(t)
	app := bootAuth(t, server.BrokerYAML("shop"), nil)

	login(t, app, server)

	resp := app.Get(t, "/en/account/claims")

	assert.JSONEq(t, shopperClaims, resp.Body)
}

// RefreshIdentity only marks the token. The next identification gets a new one from the provider and keeps it.
func TestOIDC_RefreshIdentityRefreshesToken(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	server := oidcfake.Start(t)
	app := bootAuth(t, server.BrokerYAML("shop"), nil)

	login(t, app, server)

	require.Equal(t, []string{"authorization_code"}, server.GrantTypes())

	resp := app.Get(t, "/en/account/refresh")

	require.JSONEq(t, `{"refreshed":true}`, resp.Body)
	require.Equal(t, []string{"authorization_code"}, server.GrantTypes())

	resp = app.Get(t, "/en/account/claims")

	assert.JSONEq(t, shopperClaims, resp.Body)
	assert.Equal(t, []string{"authorization_code", "refresh_token"}, server.GrantTypes())

	app.Get(t, "/en/account/claims")

	assert.Equal(t, []string{"authorization_code", "refresh_token"}, server.GrantTypes(), "the new token is kept")
}
