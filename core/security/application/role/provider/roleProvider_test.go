package provider

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/core/auth/application/fake"
	authDomain "flamingo.me/flamingo/v3/core/auth/domain"
	securityDomain "flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	AuthRoleProviderTestSuite struct {
		suite.Suite

		provider    *AuthRoleProvider
		userService *fake.UserService

		context context.Context
	}
)

func TestAuthRoleProviderTestSuite(t *testing.T) {
	suite.Run(t, &AuthRoleProviderTestSuite{})
}

func (t *AuthRoleProviderTestSuite) SetupSuite() {
	t.context = context.Background()
	t.userService = &fake.UserService{}
	t.provider = &AuthRoleProvider{}
	t.provider.Inject(t.userService)
}

func (t *AuthRoleProviderTestSuite) TestAll_Empty() {
	webSession := web.EmptySession()
	t.Equal([]securityDomain.Role(nil), t.provider.All(t.context, webSession))
}

func (t *AuthRoleProviderTestSuite) TestAll_RoleUser() {
	webSession := web.EmptySession()
	webSession.Store(fake.UserSessionKey, authDomain.User{
		Type: authDomain.USER,
	})
	t.Equal([]securityDomain.Role{
		securityDomain.RoleUser,
	}, t.provider.All(t.context, webSession))
}
