package role

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/core/security/application/role/mocks"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	ServiceImplTestSuite struct {
		suite.Suite

		service        *ServiceImpl
		firstProvider  *mocks.Provider
		secondProvider *mocks.Provider
		thirdProvider  *mocks.Provider

		context    context.Context
		webSession *web.Session
	}
)

func TestServiceImplTestSuite(t *testing.T) {
	suite.Run(t, &ServiceImplTestSuite{})
}

func (t *ServiceImplTestSuite) SetupSuite() {
	t.context = context.Background()
	t.webSession = web.EmptySession()
}

func (t *ServiceImplTestSuite) SetupTest() {
	t.firstProvider = &mocks.Provider{}
	t.secondProvider = &mocks.Provider{}
	t.thirdProvider = &mocks.Provider{}
	providers := []Provider{
		t.firstProvider,
		t.secondProvider,
		t.thirdProvider,
	}
	t.service = &ServiceImpl{}
	t.service.Inject(providers, &struct {
		PermissionHierarchy config.Map `inject:"config:core.security.roles.permissionHierarchy"`
	}{})
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
		domain.StringRole("SomePermission"),
	}
	t.firstProvider.On("All", t.context, t.webSession).Return(roles).Once()
	t.secondProvider.On("All", t.context, t.webSession).Return(roles).Once()
	t.thirdProvider.On("All", t.context, t.webSession).Return(roles).Once()

	t.Equal([]string{"SomePermission"}, t.service.AllPermissions(t.context, t.webSession))
}

func (t *ServiceImplTestSuite) TestAll_UseHierarchy() {
	firstRoles := []domain.Role{
		domain.StringRole("Permission1"),
	}
	secondRoles := []domain.Role{
		domain.StringRole("Permission2"),
	}
	thirdRoles := []domain.Role{
		domain.StringRole("Permission3"),
	}

	t.service.permissionHierarchy = map[string][]string{
		"Permission1": {"Permission11"},
		"Permission2": {"Permission21", "Permission22"},
		"Permission3": {"Permission31", "Permission32", "Permission33"},
	}

	t.firstProvider.On("All", t.context, t.webSession).Return(firstRoles).Once()
	t.secondProvider.On("All", t.context, t.webSession).Return(secondRoles).Once()
	t.thirdProvider.On("All", t.context, t.webSession).Return(thirdRoles).Once()

	t.ElementsMatch([]string{
		"Permission1",
		"Permission11",
		"Permission2",
		"Permission21",
		"Permission22",
		"Permission3",
		"Permission31",
		"Permission32",
		"Permission33",
	}, t.service.AllPermissions(t.context, t.webSession))
}

func (t *ServiceImplTestSuite) TestAll_Complete() {
	firstRoles := []domain.Role{
		domain.StringRole("Permission1"),
	}
	secondRoles := []domain.Role{
		domain.StringRole("Permission2"),
	}
	thirdRoles := []domain.Role{
		domain.StringRole("Permission3"),
	}

	t.service.permissionHierarchy = map[string][]string{
		"Permission1": {"Permission11", "Permission21", "Permission31"},
		"Permission2": {"Permission21", "Permission22", "Permission32"},
		"Permission3": {"Permission31", "Permission32", "Permission33"},
	}

	t.firstProvider.On("All", t.context, t.webSession).Return(firstRoles).Once()
	t.secondProvider.On("All", t.context, t.webSession).Return(secondRoles).Once()
	t.thirdProvider.On("All", t.context, t.webSession).Return(thirdRoles).Once()

	t.ElementsMatch([]string{
		"Permission1",
		"Permission11",
		"Permission2",
		"Permission21",
		"Permission22",
		"Permission3",
		"Permission31",
		"Permission32",
		"Permission33",
	}, t.service.AllPermissions(t.context, t.webSession))
}
