package auth

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

type testIdentifier struct{}

func (*testIdentifier) Broker() string {
	return "test"
}

type testIdentity struct{}

func (*testIdentifier) Identify(ctx context.Context, request *web.Request) (Identity, error) {
	return &testIdentity{}, nil
}

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

	identity, err := s.IdentifyAs(context.Background(), nil, new(testIdentityType))
	assert.NoError(t, err)
	testIdentity, ok := identity.(testIdentityType)
	assert.True(t, ok)
	assert.True(t, testIdentity.TestIdentity())

	identity, err = s.IdentifyAs(context.Background(), nil, new(testNotImplementedIdentityType))
	assert.Error(t, err)
	t.Log(err)
	assert.Nil(t, identity)

	identity, err = s.IdentifyAs(context.Background(), nil, 123)
	assert.Error(t, err)
	t.Log(err)
	assert.Nil(t, identity)

	identity, err = s.IdentifyAs(context.Background(), nil, new(string))
	assert.Error(t, err)
	t.Log(err)
	assert.Nil(t, identity)
}
