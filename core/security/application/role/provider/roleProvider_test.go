package provider

import (
	"context"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"

	"flamingo.me/flamingo/core/auth/application/fake"
	authDomain "flamingo.me/flamingo/core/auth/domain"
	securityDomain "flamingo.me/flamingo/core/security/domain"
	"flamingo.me/flamingo/framework/web"
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
	session := sessions.NewSession(nil, "-")
	webSession := web.NewSession(session)
	t.Equal([]securityDomain.Role(nil), t.provider.All(t.context, webSession))
}

func (t *AuthRoleProviderTestSuite) TestAll_RoleUser() {
	session := sessions.NewSession(nil, "-")
	session.Values[fake.UserSessionKey] = authDomain.User{
		Type: authDomain.USER,
	}
	webSession := web.NewSession(session)
	t.Equal([]securityDomain.Role{
		securityDomain.RoleUser,
	}, t.provider.All(t.context, webSession))
}
