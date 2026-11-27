package session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSession_RedisCookieAttributes(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	for name, secure := range map[string]bool{"secure": true, "not secure": false} { //nolint:paralleltest // boots an app, dingo scopes are process-global
		t.Run(name, func(t *testing.T) {
			app, server := bootRedis(t, nil, shopCookieLines(secure)...)

			cookie := assertShopCookie(t, app, secure)

			assert.Equal(t, []string{"session:" + cookie.Value}, server.Keys())
		})
	}
}
