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

var testTypeChecker = func(identity Identity) bool {
	_, ok := identity.(testIdentityType)

	return ok
}

func Test_WebIdentityServiceIdentifyAs(t *testing.T) {
	s := &WebIdentityService{identityProviders: []RequestIdentifier{new(testIdentifier)}}

	t.Run("existing identification", func(t *testing.T) {
		identity, err := s.IdentifyAs(context.Background(), nil, testTypeChecker)
		assert.NoError(t, err)
		testIdentity, ok := identity.(testIdentityType)
		assert.True(t, ok)
		assert.True(t, testIdentity.TestIdentity())
	})

	t.Run("non-existing indentification", func(t *testing.T) {
		identity, err := s.IdentifyAs(context.Background(), nil, func(identity Identity) bool {
			_, ok := identity.(testNotImplementedIdentityType)

			return ok
		})
		assert.Error(t, err)
		t.Log(err)
		assert.Nil(t, identity)
	})
}
