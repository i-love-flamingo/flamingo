package session_test

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/testapp/oidcfake"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/web"
)

const storedSessionFile = "../testdata/session/stored.json"

var update = flag.Bool("update", false, "write the stored session fixture again")

// storedSession is the Redis entry of one session.
type storedSession struct {
	Cookie string `json:"cookie"` // the value of the session cookie, which is the session ID
	Key    string `json:"key"`    // the Redis key
	Value  string `json:"value"`  // the Redis value, base64
	State  string `json:"state"`  // the state of the pending login
}

// Two versions of an app share one Redis during a deploy, so a new release has to read old sessions.
// The fixture holds a logged-in user, a pending login and a basket item, and flamingo v3.17.6 wrote it.
// -update writes the file with the code of the checkout. After a deliberate format change, run it on the
// release that has the new format and name that release here:
//
//	go test ./contract/session -run TestStoredFormat -update
func TestStoredFormat_ReleasedSessionReads(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	// a session as this code writes it
	var current storedSession

	require.True(t, t.Run("current", func(t *testing.T) { current = newStoredSession(t) }))

	if *update {
		writeStoredSession(t, current)
	}

	raw, err := os.ReadFile(storedSessionFile)
	require.NoError(t, err)

	var stored storedSession
	require.NoError(t, json.Unmarshal(raw, &stored))

	app, server := bootRedis(t, nil)

	instance, err := app.Injector(t).GetInstance(new(web.SessionStore))
	require.NoError(t, err)

	store, ok := instance.(*web.SessionStore)
	require.True(t, ok, "the store is %T", instance)

	load := func(entry storedSession) *web.Session {
		value, err := base64.StdEncoding.DecodeString(entry.Value)
		require.NoError(t, err)
		require.NoError(t, server.Set(entry.Key, string(value)))

		session, err := store.LoadByID(t.Context(), entry.Cookie)
		require.NoError(t, err)

		return session
	}

	session := load(stored)

	// a renamed session key would leave the old data unread
	assert.Subset(t, load(current).Keys(), session.Keys())

	_, hasItem := session.Load("basket.a")
	assert.True(t, hasItem, "item of the basket")

	login, hasLogin := session.Load("core.auth.oidc.shop.sessiondata")
	require.True(t, hasLogin, "data of the login")
	assertLogin(t, login)

	// importing the oauth package registers its gob types
	states, ok := session.Try("core.auth.oidc.shop.states").([]oauth.StateEntry)
	require.True(t, ok, "pending logins")
	require.Len(t, states, 1)
	assert.Equal(t, stored.State, states[0].State)
	assert.False(t, states[0].TS.IsZero())

	resp := app.Get(t, "/en/basket", "Cookie", "flamingo="+stored.Cookie)

	assert.JSONEq(t, `["basket.a"]`, resp.Body)
}

// assertLogin checks the data that the oidc broker keeps. Its type is not exported, so the check uses reflect.
// gob drops a stored field that the type no longer has without an error, so the check reads every field.
func assertLogin(t *testing.T, data any) {
	t.Helper()

	value := reflect.Indirect(reflect.ValueOf(data))
	require.Equal(t, reflect.Struct, value.Kind(), "data of the login is a %T", data)

	field := func(name string) reflect.Value {
		found := value.FieldByName(name)
		require.True(t, found.IsValid(), "data of the login has no field %s", name)

		return found
	}

	assert.Equal(t, "shopper", field("Subject").String())
	assert.NotEmpty(t, field("RawIDToken").String())
	assert.JSONEq(t, `{"userID":"shopper","mail":"shopper@example.com"}`, string(field("IDTokenClaims").Bytes()))
	assert.JSONEq(t, `{"tokenSubject":"shopper"}`, string(field("AccessTokenClaims").Bytes()))

	token, ok := field("Token").Interface().(*oauth2.Token)
	require.True(t, ok, "field Token is a %T", field("Token").Interface())
	require.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
}

// newStoredSession logs in, adds a basket item, starts a second login and gives the Redis entry of the session.
// Call it in a subtest, so that its app stops before the test starts the next one.
func newStoredSession(t *testing.T) storedSession {
	t.Helper()

	provider := oidcfake.Start(t)
	app, server := bootRedis(t, []dingo.Module{new(auth.WebModule), new(oauth.Module)}, provider.BrokerYAML("shop"))

	provider.QueueUser("shopper", "shopper@example.com")

	// a client that follows the redirects through the provider and back
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, app.URL+"/en/core/auth/login/shop", nil)
	require.NoError(t, err)

	resp, err := (&http.Client{Jar: app.Client.Jar}).Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	app.Get(t, "/en/basket/add/a")

	pending := app.Get(t, "/en/core/auth/login/shop")
	require.Equal(t, http.StatusSeeOther, pending.Status)

	authorize, err := url.Parse(pending.Header.Get("Location"))
	require.NoError(t, err)

	keys := server.Keys()
	require.Len(t, keys, 1)

	value, err := server.Get(keys[0])
	require.NoError(t, err)

	return storedSession{
		Cookie: strings.TrimPrefix(keys[0], "session:"),
		Key:    keys[0],
		Value:  base64.StdEncoding.EncodeToString([]byte(value)),
		State:  authorize.Query().Get("state"),
	}
}

func writeStoredSession(t *testing.T, session storedSession) {
	t.Helper()

	raw, err := json.MarshalIndent(session, "", "  ")
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(storedSessionFile, append(raw, '\n'), 0o600))
}
