package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/framework/web"
)

type testIdentifier struct{}

func (*testIdentifier) Broker() string {
	return "test"
}

func (*testIdentifier) Identify(context.Context, *web.Request) (Identity, error) {
	return &testIdentity{}, nil
}

type testIdentity struct{}

func (*testIdentity) Subject() string {
	return "test-identity"
}

func (*testIdentity) Broker() string {
	return "test"
}

func (*testIdentity) TestIdentity() bool {
	return true
}

type testIdentityType interface {
	TestIdentity() bool
}

type testNotImplementedIdentityType interface {
	TestIdentityNotImplemented() bool
}

func Test_WebIdentityServiceIdentifyAs(t *testing.T) {
	s := &WebIdentityService{identityProviders: []RequestIdentifier{new(testIdentifier)}}

	t.Run("existing identification", func(t *testing.T) {
		identity, err := s.IdentifyAs(context.Background(), nil, new(testIdentityType))
		assert.NoError(t, err)
		testIdentity, ok := identity.(testIdentityType)
		assert.True(t, ok)
		assert.True(t, testIdentity.TestIdentity())
	})

	t.Run("non-existing indentification", func(t *testing.T) {
		identity, err := s.IdentifyAs(context.Background(), nil, new(testNotImplementedIdentityType))
		assert.Error(t, err)
		t.Log(err)
		assert.Nil(t, identity)
	})

	t.Run("must use pointer", func(t *testing.T) {
		identity, err := s.IdentifyAs(context.Background(), nil, testIdentity{})
		assert.Error(t, err)
		t.Log(err)
		assert.Nil(t, identity)
	})

	t.Run("must use interface", func(t *testing.T) {
		identity, err := s.IdentifyAs(context.Background(), nil, new(testIdentity))
		assert.Error(t, err)
		t.Log(err)
		assert.Nil(t, identity)
	})
}
