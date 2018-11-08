package voter

import (
	"context"
	"testing"

	roleMocks "flamingo.me/flamingo/core/security/application/role/mocks"
	"flamingo.me/flamingo/core/security/domain"
	domainMocks "flamingo.me/flamingo/core/security/domain/mocks"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	RoleVoterTestSuite struct {
		suite.Suite

		voter       *RoleVoter
		roleService *roleMocks.Service
		object      *domainMocks.RoleSet

		context context.Context
		session *sessions.Session
	}
)

func TestRoleVoterTestSuite(t *testing.T) {
	suite.Run(t, &RoleVoterTestSuite{})
}

func (t *RoleVoterTestSuite) SetupSuite() {
	t.context = context.Background()
	t.session = sessions.NewSession(nil, "-")
}

func (t *RoleVoterTestSuite) SetupTest() {
	t.roleService = &roleMocks.Service{}
	t.object = &domainMocks.RoleSet{}
	t.voter = &RoleVoter{}
	t.voter.Inject(t.roleService)
}

func (t *RoleVoterTestSuite) TearDownTest() {
	t.roleService.AssertExpectations(t.T())
	t.roleService = nil
	t.object.AssertExpectations(t.T())
	t.object = nil
	t.voter = nil
}

func (t *RoleVoterTestSuite) TestVote_AccessAbstained() {
	t.Equal(AccessAbstained, t.voter.Vote(t.context, t.session, domain.RoleAnonymous.Permission(), nil))
	t.Equal(AccessAbstained, t.voter.Vote(t.context, t.session, domain.RoleUser.Permission(), nil))
}

func (t *RoleVoterTestSuite) TestVote_AccessGrantedWithoutObject() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleUser,
		domain.DefaultRole("RoleAdministrator"),
	}).Once()
	t.Equal(AccessGranted, t.voter.Vote(t.context, t.session, "RoleAdministrator", nil))
}

func (t *RoleVoterTestSuite) TestVote_AccessGrantedWithObject() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleUser,
		domain.DefaultRole("RoleAdministrator"),
	}).Once()
	t.object.On("Roles").Return([]domain.Role{
		domain.DefaultRole("RoleAdministrator"),
	}).Once()
	t.Equal(AccessGranted, t.voter.Vote(t.context, t.session, "RoleAdministrator", t.object))
}

func (t *RoleVoterTestSuite) TestVote_AccessDeniedWithoutObject() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleAnonymous,
	}).Once()
	t.Equal(AccessDenied, t.voter.Vote(t.context, t.session, "RoleAdministrator", nil))
}

func (t *RoleVoterTestSuite) TestVote_AccessDeniedWithObject() {
	t.object.On("Roles").Return([]domain.Role{}).Once()
	t.Equal(AccessDenied, t.voter.Vote(t.context, t.session, "RoleAdministrator", t.object))
}
