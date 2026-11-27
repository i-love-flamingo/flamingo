package session_test

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/basket"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

// bootSession starts an app with the session module and the area "en", at the path "en", with the basket routes.
// Extra modules load in the root area before the session module.
// The lines go into config.yml. The session cookie is not secure unless a line says otherwise, so the tests do
// not depend on how the cookie jar treats secure cookies.
func bootSession(tb testing.TB, extra []dingo.Module, lines ...string) *harness.App {
	tb.Helper()

	if !slices.ContainsFunc(lines, func(line string) bool { return strings.HasPrefix(line, "flamingo.session.cookie.secure:") }) {
		lines = append(lines, "flamingo.session.cookie.secure: false")
	}

	root := slices.Concat(
		[]dingo.Module{new(prefixrouter.Module), harness.PassThroughWrapper()},
		extra,
		[]dingo.Module{new(flamingo.SessionModule)},
	)

	return harness.Boot(tb, root, harness.Config{
		Files: map[string]string{
			"config.yml":    strings.Join(lines, "\n") + "\n",
			"en/config.yml": "flamingo.router.path: en\n",
		},
		Areas: []*config.Area{
			config.NewArea("en", []dingo.Module{new(basket.Module)}),
		},
	})
}

// bootRedis starts miniredis and an app that keeps its sessions in it.
func bootRedis(tb testing.TB, extra []dingo.Module, lines ...string) (*harness.App, *miniredis.Miniredis) {
	tb.Helper()

	server := miniredis.RunT(tb)

	lines = append(lines,
		"flamingo.session.backend: redis",
		"flamingo.session.redis.host: \""+server.Addr()+"\"",
	)

	return bootSession(tb, extra, lines...), server
}

// sessionCookie gives the cookie of the response with the name.
func sessionCookie(tb testing.TB, resp harness.Response, name string) *http.Cookie {
	tb.Helper()

	for _, cookie := range (&http.Response{Header: resp.Header}).Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}

	tb.Fatalf("response has no cookie %q, Set-Cookie: %q", name, resp.Header.Values("Set-Cookie"))

	return nil
}

// assertShopCookie checks the attributes of the session cookie for the config
// name "shopsession", max age 600, same site strict, path "/en" and the given secure flag.
// The next request sends the cookie in a header, so it does not depend on the cookie jar.
func assertShopCookie(t *testing.T, app *harness.App, secure bool) *http.Cookie {
	t.Helper()

	resp := app.Get(t, "/en/basket/add/a")

	cookie := sessionCookie(t, resp, "shopsession")

	assert.Equal(t, 600, cookie.MaxAge)
	assert.Equal(t, secure, cookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.Equal(t, "/en", cookie.Path)
	assert.True(t, cookie.HttpOnly)

	next := app.Get(t, "/en/basket", "Cookie", cookie.Name+"="+cookie.Value)

	require.JSONEq(t, `["basket.a"]`, next.Body)

	return cookie
}

// shopCookieLines are the config lines that assertShopCookie expects, with the given secure flag.
func shopCookieLines(secure bool) []string {
	return []string{
		"flamingo.session.name: shopsession",
		"flamingo.session.max.age: 600",
		"flamingo.session.cookie.secure: " + strconv.FormatBool(secure),
		"flamingo.session.cookie.sameSite: strict",
		"flamingo.session.cookie.path: /en",
	}
}
