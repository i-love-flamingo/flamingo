package session_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/web"
)

func TestStore_LoadByIDAndSaveWithoutRequest(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app, _ := bootRedis(t, nil)

	resp := app.Get(t, "/en/basket/add/a")

	cookie := sessionCookie(t, resp, "flamingo")

	instance, err := app.Injector(t).GetInstance(new(web.SessionStore))
	require.NoError(t, err)

	store, ok := instance.(*web.SessionStore)
	require.True(t, ok, "the store is %T", instance)

	ctx := context.Background()

	session, err := store.LoadByID(ctx, cookie.Value)
	require.NoError(t, err)

	assert.Equal(t, "1", session.Try("basket.a"))

	session.Store("basket.a", "2")

	_, err = store.Save(ctx, session)
	require.NoError(t, err)

	reloaded, err := store.LoadByID(ctx, cookie.Value)
	require.NoError(t, err)

	assert.Equal(t, "2", reloaded.Try("basket.a"))
}
