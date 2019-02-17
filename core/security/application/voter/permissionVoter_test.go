package voter

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/core/security/domain"
	domainMocks "flamingo.me/flamingo/v3/core/security/domain/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	PermissionVoterTestSuite struct {
		suite.Suite

		voter       *PermissionVoter
		object      *domainMocks.PermissionSet

		context    context.Context
		webSession *web.Session
	}
)

func TestPermissionVoterTestSuite(t *testing.T) {
	suite.Run(t, &PermissionVoterTestSuite{})
}

func (t *PermissionVoterTestSuite) SetupSuite() {
	t.context = context.Background()
	t.webSession = web.EmptySession()
}

func (t *PermissionVoterTestSuite) SetupTest() {
	t.object = &domainMocks.PermissionSet{}
	t.voter = &PermissionVoter{}
}

func (t *PermissionVoterTestSuite) TearDownTest() {
	t.object.AssertExpectations(t.T())
	t.object = nil
	t.voter = nil
}

func (t *PermissionVoterTestSuite) TestVote_AccessAbstained() {
	t.Equal(AccessAbstained, t.voter.Vote([]string{}, domain.PermissionAuthorized, nil))
}

func (t *PermissionVoterTestSuite) TestVote_AccessGrantedWithoutObject() {
	t.Equal(AccessGranted, t.voter.Vote([]string{
		domain.PermissionAuthorized,
		"PermissionAdministrator",
	}, "PermissionAdministrator", nil))
}

func (t *PermissionVoterTestSuite) TestVote_AccessGrantedWithObject() {
	t.object.On("Permissions").Return([]string{
		domain.PermissionAuthorized,
		"PermissionAdministrator",
	}).Once()
	t.Equal(AccessGranted, t.voter.Vote([]string{
		domain.PermissionAuthorized,
		"PermissionAdministrator",
	}, "PermissionAdministrator", t.object))
}

func (t *PermissionVoterTestSuite) TestVote_AccessDeniedWithoutObject() {
	t.Equal(AccessDenied, t.voter.Vote([]string{}, "PermissionAdministrator", nil))
}

func (t *PermissionVoterTestSuite) TestVote_AccessDeniedWithObject() {
	t.object.On("Permissions").Return([]string{}).Once()
	t.Equal(AccessDenied, t.voter.Vote([]string{}, "RoleAdministrator", t.object))
}
