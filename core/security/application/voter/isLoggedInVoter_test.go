package voter

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	IsLoggedInVoterTestSuite struct {
		suite.Suite

		voter *IsLoggedInVoter

		context    context.Context
		webSession *web.Session
	}
)

func TestIsLoggedInVoterTestSuite(t *testing.T) {
	suite.Run(t, &IsLoggedInVoterTestSuite{})
}

func (t *IsLoggedInVoterTestSuite) SetupSuite() {
	t.context = context.Background()
	t.webSession = web.EmptySession()
}

func (t *IsLoggedInVoterTestSuite) SetupTest() {
	t.voter = &IsLoggedInVoter{}
}

func (t *IsLoggedInVoterTestSuite) TearDownTest() {
	t.voter = nil
}

func (t *IsLoggedInVoterTestSuite) TestVote_AccessAbstained() {
	t.Equal(AccessAbstained, t.voter.Vote([]string{}, "SomePermission", nil))
}

func (t *IsLoggedInVoterTestSuite) TestVote_AccessGranted() {
	t.Equal(AccessGranted, t.voter.Vote([]string{domain.PermissionAuthorized}, domain.PermissionAuthorized, nil))
}

func (t *IsLoggedInVoterTestSuite) TestVote_AccessDenied() {
	t.Equal(AccessDenied, t.voter.Vote([]string{}, domain.PermissionAuthorized, nil))
}
