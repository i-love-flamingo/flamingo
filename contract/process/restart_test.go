package process_test

import (
	"syscall"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestart_RedisSessionSurvives(t *testing.T) {
	t.Parallel()

	redis := miniredis.RunT(t)
	env := map[string]string{"CONTRACT_REDIS_ADDR": redis.Addr()}

	first, addrs := serve(t, testappMain, env)

	added := get(t, addrs.HTTP, "/en/basket/add/a")
	require.Equal(t, 200, added.status)
	require.Len(t, added.cookies, 1)
	assert.Equal(t, "flamingo", added.cookies[0].Name)

	first.Signal(t, syscall.SIGTERM)
	require.Equal(t, 0, first.Wait(t))

	// the session is in redis, not in the memory of the first process
	assert.NotEmpty(t, redis.Keys())

	_, addrs = serve(t, testappMain, env)

	assert.JSONEq(t, `["basket.a"]`, get(t, addrs.HTTP, "/en/basket", added.cookies...).body)
}
