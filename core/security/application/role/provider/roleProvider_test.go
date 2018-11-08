package provider

import (
	"context"
	"testing"

	"flamingo.me/flamingo/core/auth/application/fake"
	authDomain "flamingo.me/flamingo/core/auth/domain"
	securityDomain "flamingo.me/flamingo/core/security/domain"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	DefaultRoleProviderTestSuite struct {
		suite.Suite

		provider    *DefaultRoleProvider
		userService *fake.UserService

		context context.Context
	}
)

func TestDefaultRoleProviderTestSuite(t *testing.T) {
	suite.Run(t, &DefaultRoleProviderTestSuite{})
}

func (t *DefaultRoleProviderTestSuite) SetupSuite() {
	t.context = context.Background()
	t.userService = &fake.UserService{}
	t.provider = &DefaultRoleProvider{}
	t.provider.Inject(t.userService)
}

func (t *DefaultRoleProviderTestSuite) TestAll_RoleAnonymous() {
	session := sessions.NewSession(nil, "-")
	t.Equal([]securityDomain.Role{
		securityDomain.RoleAnonymous,
	}, t.provider.All(t.context, session))
}

func (t *DefaultRoleProviderTestSuite) TestAll_RoleUser() {
	session := sessions.NewSession(nil, "-")
	session.Values[fake.UserSessionKey] = authDomain.User{
		Type: authDomain.USER,
	}
	t.Equal([]securityDomain.Role{
		securityDomain.RoleUser,
	}, t.provider.All(t.context, session))
}
