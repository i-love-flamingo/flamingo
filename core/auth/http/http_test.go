package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

func TestHTTPBasicAuthIdentifier(t *testing.T) {
	identifier := new(basicAuthIdentifier).Inject(&struct {
		Users config.Map `inject:"config:core.auth.httpbasicusers"`
	}{
		Users: config.Map{
			"alice": "secretpass123",
			"bob":   "donothackmepls",
		},
	})

	t.Run("alice/correct", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("alice", "secretpass123")
		identity := identifier.Identify(context.Background(), req)
		assert.NotNil(t, identity)
		assert.Equal(t, "alice", identity.Subject())
	})

	t.Run("alice/wrong", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("alice", "secretpass12")
		identity := identifier.Identify(context.Background(), req)
		assert.Nil(t, identity)
	})

	t.Run("bob/correct", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("bob", "donothackmepls")
		identity := identifier.Identify(context.Background(), req)
		assert.NotNil(t, identity)
		assert.Equal(t, "bob", identity.Subject())
	})

	t.Run("alice/wrong", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("bob", "---")
		identity := identifier.Identify(context.Background(), req)
		assert.Nil(t, identity)
	})

	t.Run("unknown", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		req.Request().SetBasicAuth("unknown", "---")
		identity := identifier.Identify(context.Background(), req)
		assert.Nil(t, identity)
	})

	t.Run("none", func(t *testing.T) {
		req := web.CreateRequest(nil, nil)
		identity := identifier.Identify(context.Background(), req)
		assert.Nil(t, identity)
	})
}
