package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveMode_OnWriteReadOfExistingSessionWritesNothing(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app, server := bootRedis(t, nil, "flamingo.session.saveMode: OnWrite")

	resp := app.Get(t, "/en/basket/count")

	assert.NotEmpty(t, resp.Header.Values("Set-Cookie"), "a request without a session cookie starts a session")
	require.Len(t, server.Keys(), 1, "the new session is stored")

	app.Get(t, "/en/basket/add/a")

	// A new write of the session sets the lifetime of the Redis key again.
	server.FastForward(time.Minute)

	keys := server.Keys()

	require.Len(t, keys, 1)

	lifetime := server.TTL(keys[0])

	resp = app.Get(t, "/en/basket/count")

	require.JSONEq(t, `{"count":1}`, resp.Body, "the request reads the item")
	assert.Empty(t, resp.Header.Values("Set-Cookie"))
	assert.Equal(t, keys, server.Keys())
	assert.Equal(t, lifetime, server.TTL(keys[0]))
}
