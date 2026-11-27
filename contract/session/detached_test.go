package session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetached_KeepsRequestAndSession(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	app := bootSession(t, nil)

	resp := app.Get(t, "/en/basket/detached")

	assert.JSONEq(t, `{"request":true,"session":true,"live":true}`, resp.Body)
}
