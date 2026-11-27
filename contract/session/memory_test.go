package session_test

import "testing"

// "true" is the config default for the secure flag, so the test also runs with "false".
func TestSession_MemoryCookieAttributes(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	for name, secure := range map[string]bool{"secure": true, "not secure": false} { //nolint:paralleltest // boots an app, dingo scopes are process-global
		t.Run(name, func(t *testing.T) {
			app := bootSession(t, nil, shopCookieLines(secure)...)

			assertShopCookie(t, app, secure)
		})
	}
}
