package voter

import (
	"context"
	"testing"

	"flamingo.me/flamingo/core/security/application/role/mocks"
	"flamingo.me/flamingo/core/security/domain"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	IsLoggedInVoterTestSuite struct {
		suite.Suite

		voter       *IsLoggedInVoter
		roleService *mocks.Service

		context context.Context
		session *sessions.Session
	}
)

func TestIsLoggedInVoterTestSuite(t *testing.T) {
	suite.Run(t, &IsLoggedInVoterTestSuite{})
}

func (t *IsLoggedInVoterTestSuite) SetupSuite() {
	t.context = context.Background()
	t.session = sessions.NewSession(nil, "-")
}

func (t *IsLoggedInVoterTestSuite) SetupTest() {
	t.roleService = &mocks.Service{}
	t.voter = &IsLoggedInVoter{}
	t.voter.Inject(t.roleService)
}

func (t *IsLoggedInVoterTestSuite) TearDownTest() {
	t.roleService.AssertExpectations(t.T())
	t.roleService = nil
	t.voter = nil
}

func (t *IsLoggedInVoterTestSuite) TestVote_AccessAbstained() {
	t.Equal(AccessAbstained, t.voter.Vote(t.context, t.session, "SomePermission", nil))
}

func (t *IsLoggedInVoterTestSuite) TestVote_AccessGranted() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleUser,
	}).Once()
	t.Equal(AccessGranted, t.voter.Vote(t.context, t.session, domain.RoleUser.Permission(), nil))
}

func (t *IsLoggedInVoterTestSuite) TestVote_AccessDenied() {
	t.roleService.On("All", t.context, t.session).Return([]domain.Role{
		domain.RoleAnonymous,
	}).Once()
	t.Equal(AccessDenied, t.voter.Vote(t.context, t.session, domain.RoleUser.Permission(), nil))
}
