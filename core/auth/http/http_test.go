package http

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

func TestHTTPBasicAuthIdentifier(t *testing.T) {
	identifier, err := identifierFactory(config.Map{
		"realm":  "test",
		"broker": "test",
		"users": map[string]interface{}{
			"alice": "secretpass123",
			"bob":   "donothackmepls",
		},
	})

	assert.NoError(t, err)

	t.Run("alice/correct", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("alice", "secretpass123")
		identity, err := identifier.Identify(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, identity)
		assert.Equal(t, "alice", identity.Subject())
	})

	t.Run("alice/wrong", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("alice", "secretpass12")
		identity, err := identifier.Identify(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})

	t.Run("bob/correct", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("bob", "donothackmepls")
		identity, err := identifier.Identify(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, identity)
		assert.Equal(t, "bob", identity.Subject())
	})

	t.Run("alice/wrong", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("bob", "---")
		identity, err := identifier.Identify(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})

	t.Run("unknown", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("unknown", "---")
		identity, err := identifier.Identify(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})

	t.Run("none", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		identity, err := identifier.Identify(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, identity)
	})
}
