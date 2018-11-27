package role

import (
	"context"
	"testing"

	"flamingo.me/flamingo/core/security/application/role/provider"
	"flamingo.me/flamingo/core/security/application/role/provider/mocks"
	"flamingo.me/flamingo/core/security/domain"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	ServiceImplTestSuite struct {
		suite.Suite

		service        *ServiceImpl
		firstProvider  *mocks.RoleProvider
		secondProvider *mocks.RoleProvider
		thirdProvider  *mocks.RoleProvider

		context    context.Context
		session    *sessions.Session
		webSession *web.Session
	}
)

func TestServiceImplTestSuite(t *testing.T) {
	suite.Run(t, &ServiceImplTestSuite{})
}

func (t *ServiceImplTestSuite) SetupSuite() {
	t.context = context.Background()
	t.session = sessions.NewSession(nil, "-")
	t.webSession = web.NewSession(t.session)
}

func (t *ServiceImplTestSuite) SetupTest() {
	t.firstProvider = &mocks.RoleProvider{}
	t.secondProvider = &mocks.RoleProvider{}
	t.thirdProvider = &mocks.RoleProvider{}
	providers := []provider.RoleProvider{
		t.firstProvider,
		t.secondProvider,
		t.thirdProvider,
	}
	t.service = &ServiceImpl{}
	t.service.Inject(providers, &struct {
		RolesHierarchy config.Map `inject:"config:security.roles.hierarchy"`
	}{
		RolesHierarchy: config.Map{},
	})
}

func (t *ServiceImplTestSuite) TearDownTest() {
	t.firstProvider.AssertExpectations(t.T())
	t.firstProvider = nil
	t.secondProvider.AssertExpectations(t.T())
	t.secondProvider = nil
	t.thirdProvider.AssertExpectations(t.T())
	t.thirdProvider = nil
	t.service = nil
}

func (t *ServiceImplTestSuite) TestAll_RemoveDuplicates() {
	roles := []domain.Role{
		"SomePermission",
	}
	t.firstProvider.On("All", t.context, t.webSession).Return(roles).Once()
	t.secondProvider.On("All", t.context, t.webSession).Return(roles).Once()
	t.thirdProvider.On("All", t.context, t.webSession).Return(roles).Once()

	t.Equal(roles, t.service.All(t.context, t.webSession))
}

func (t *ServiceImplTestSuite) TestAll_UseHierarchy() {
	firstRoles := []domain.Role{
		"Permission1",
	}
	secondRoles := []domain.Role{
		"Permission2",
	}
	thirdRoles := []domain.Role{
		"Permission3",
	}

	t.service.rolesHierarchy = config.Map{
		"Permission1": config.Slice{"Permission11"},
		"Permission2": config.Slice{"Permission21", "Permission22"},
		"Permission3": config.Slice{"Permission31", "Permission32", "Permission33"},
	}

	t.firstProvider.On("All", t.context, t.webSession).Return(firstRoles).Once()
	t.secondProvider.On("All", t.context, t.webSession).Return(secondRoles).Once()
	t.thirdProvider.On("All", t.context, t.webSession).Return(thirdRoles).Once()

	t.ElementsMatch([]domain.Role{
		"Permission1",
		"Permission11",
		"Permission2",
		"Permission21",
		"Permission22",
		"Permission3",
		"Permission31",
		"Permission32",
		"Permission33",
	}, t.service.All(t.context, t.webSession))
}

func (t *ServiceImplTestSuite) TestAll_Complete() {
	firstRoles := []domain.Role{
		"Permission1",
	}
	secondRoles := []domain.Role{
		"Permission2",
	}
	thirdRoles := []domain.Role{
		"Permission3",
	}

	t.service.rolesHierarchy = config.Map{
		"Permission1": config.Slice{"Permission11", "Permission21", "Permission31"},
		"Permission2": config.Slice{"Permission21", "Permission22", "Permission32"},
		"Permission3": config.Slice{"Permission31", "Permission32", "Permission33"},
	}

	t.firstProvider.On("All", t.context, t.webSession).Return(firstRoles).Once()
	t.secondProvider.On("All", t.context, t.webSession).Return(secondRoles).Once()
	t.thirdProvider.On("All", t.context, t.webSession).Return(thirdRoles).Once()

	t.ElementsMatch([]domain.Role{
		"Permission1",
		"Permission11",
		"Permission2",
		"Permission21",
		"Permission22",
		"Permission3",
		"Permission31",
		"Permission32",
		"Permission33",
	}, t.service.All(t.context, t.webSession))
}
