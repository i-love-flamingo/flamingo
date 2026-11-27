package auth_test

import (
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/account"
	"flamingo.me/flamingo/v3/contract/internal/testapp/oidcfake"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/core/security"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

// bootAuth starts an app with the auth, oauth and security modules in the root area and the account routes in
// the area "en", which is at the path "en". Extra modules load in the root area after them.
// The config is the broker part (see oidcfake.Server.BrokerYAML) and a session cookie that is not secure.
func bootAuth(tb testing.TB, brokers string, extra []dingo.Module) *harness.App {
	tb.Helper()

	root := slices.Concat(
		[]dingo.Module{
			new(prefixrouter.Module),
			harness.PassThroughWrapper(),
			new(auth.WebModule),
			new(oauth.Module),
			new(security.Module),
		},
		extra,
	)

	return harness.Boot(tb, root, harness.Config{
		Files: map[string]string{
			"config.yml":    brokers + "flamingo.session.cookie.secure: false\n",
			"en/config.yml": "flamingo.router.path: en\n",
		},
		Areas: []*config.Area{config.NewArea("en", []dingo.Module{new(account.Module)})},
	})
}

// redirectOf sends a GET request with the client of the app to an absolute URL, as a browser does after a redirect.
// It reads and closes the body, checks the status of the response and gives its Location.
func redirectOf(t *testing.T, app *harness.App, target string, status int) *url.URL {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, target, nil)
	require.NoError(t, err)

	resp, err := app.Client.Do(req)
	require.NoError(t, err)

	defer func() { _ = resp.Body.Close() }()

	_, err = io.Copy(io.Discard, resp.Body)
	require.NoError(t, err)

	return location(t, harness.Response{Status: resp.StatusCode, Header: resp.Header}, status)
}

// location gives the Location header of a redirect with the status, parsed.
func location(t *testing.T, resp harness.Response, status int) *url.URL {
	t.Helper()

	require.Equal(t, status, resp.Status)

	target, err := url.Parse(resp.Header.Get("Location"))
	require.NoError(t, err)

	return target
}

// login makes the user "shopper" log in at the fake provider, through the login route of the broker "shop".
// It follows the redirects up to the callback, which sets the identity in the session.
// It returns the URL of the authorize request that the login route gave.
func login(t *testing.T, app *harness.App, server *oidcfake.Server) *url.URL {
	t.Helper()

	server.QueueUser("shopper", "shopper@example.com")

	resp := app.Get(t, "/en/core/auth/login/shop")
	authorize := location(t, resp, http.StatusSeeOther)

	require.True(t, strings.HasPrefix(authorize.String(), server.Issuer+"/authorize?"), "authorize URL: %s", authorize)

	callback := redirectOf(t, app, authorize.String(), http.StatusFound)

	require.True(t, strings.HasPrefix(callback.String(), app.URL+"/en/core/auth/callback/shop?"), "callback URL: %s", callback)

	redirectOf(t, app, callback.String(), http.StatusSeeOther)

	return authorize
}

// logout calls the logout route of the broker "shop" and follows the redirects through the end session endpoint
// of the fake provider and back to the app.
func logout(t *testing.T, app *harness.App, server *oidcfake.Server) {
	t.Helper()

	resp := app.Get(t, "/en/core/auth/logout/shop")

	assert.Equal(t, "/en/core/auth/logoutCallback", location(t, resp, http.StatusSeeOther).Path)

	resp = app.Get(t, "/en/core/auth/logoutCallback")
	endSession := location(t, resp, http.StatusSeeOther)

	assert.True(t, strings.HasPrefix(endSession.String(), server.EndSession+"?"), "end session URL: %s", endSession)
	assert.NotEmpty(t, endSession.Query().Get("id_token_hint"))
	assert.Equal(t, app.URL+"/en/core/auth/logoutCallback", endSession.Query().Get("post_logout_redirect_uri"))

	back := redirectOf(t, app, endSession.String(), http.StatusFound)

	redirectOf(t, app, back.String(), http.StatusSeeOther)
}

// loggedIn gives the answer of the route "account.status".
func loggedIn(t *testing.T, app *harness.App) string {
	t.Helper()

	return app.Get(t, "/en/account/status").Body
}
