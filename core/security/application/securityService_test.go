package application

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/core/security/application/voter"
	"flamingo.me/flamingo/v3/core/security/application/voter/mocks"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	SecurityServiceTestSuite struct {
		suite.Suite

		service     *SecurityServiceImpl
		firstVoter  *mocks.SecurityVoter
		secondVoter *mocks.SecurityVoter
		thirdVoter  *mocks.SecurityVoter

		context context.Context
	}

	serviceTestCase struct {
		firstVote         int
		secondVote        int
		thirdVote         int
		voterStrategy     string
		allowIfAllAbstain bool
		decision          bool
	}
)

func TestSecurityServiceTestSuite(t *testing.T) {
	suite.Run(t, &SecurityServiceTestSuite{})
}

func (t *SecurityServiceTestSuite) SetupSuite() {
	t.context = context.Background()
}

func (t *SecurityServiceTestSuite) SetupTest() {
	t.firstVoter = &mocks.SecurityVoter{}
	t.secondVoter = &mocks.SecurityVoter{}
	t.thirdVoter = &mocks.SecurityVoter{}
	voters := []voter.SecurityVoter{
		t.firstVoter,
		t.secondVoter,
		t.thirdVoter,
	}
	t.service = &SecurityServiceImpl{}
	t.service.Inject(voters, &struct {
		VoterStrategy     string `inject:"config:security.roles.voters.strategy"`
		AllowIfAllAbstain bool   `inject:"config:security.roles.voters.allowIfAllAbstain"`
	}{})
}

func (t *SecurityServiceTestSuite) TearDownTest() {
	t.firstVoter.AssertExpectations(t.T())
	t.firstVoter = nil
	t.secondVoter.AssertExpectations(t.T())
	t.secondVoter = nil
	t.thirdVoter.AssertExpectations(t.T())
	t.thirdVoter = nil
	t.service = nil
}

func (t *SecurityServiceTestSuite) TestIsLoggedIn() {
	testCases := []serviceTestCase{
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessDenied,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: true,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessGranted,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessDenied,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: true,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessGranted,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessGranted,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: true,
			decision:          true,
		},
	}

	for _, testCase := range testCases {
		webSession := web.EmptySession()
		t.service.voterStrategy = testCase.voterStrategy
		t.service.allowIfAllAbstain = testCase.allowIfAllAbstain
		t.firstVoter.On("Vote", t.context, webSession, domain.RoleUser.Permission(), nil).Once().Return(testCase.firstVote)
		t.secondVoter.On("Vote", t.context, webSession, domain.RoleUser.Permission(), nil).Once().Return(testCase.secondVote)
		t.thirdVoter.On("Vote", t.context, webSession, domain.RoleUser.Permission(), nil).Once().Return(testCase.thirdVote)
		t.Equal(testCase.decision, t.service.IsLoggedIn(t.context, webSession))
	}
}

func (t *SecurityServiceTestSuite) TestIsLoggedOut() {
	testCases := []serviceTestCase{
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessDenied,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: true,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessGranted,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessDenied,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: true,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessGranted,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessGranted,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: true,
			decision:          false,
		},
	}

	for _, testCase := range testCases {
		webSession := web.EmptySession()
		t.service.voterStrategy = testCase.voterStrategy
		t.service.allowIfAllAbstain = testCase.allowIfAllAbstain
		t.firstVoter.On("Vote", t.context, webSession, domain.RoleUser.Permission(), nil).Once().Return(testCase.firstVote)
		t.secondVoter.On("Vote", t.context, webSession, domain.RoleUser.Permission(), nil).Once().Return(testCase.secondVote)
		t.thirdVoter.On("Vote", t.context, webSession, domain.RoleUser.Permission(), nil).Once().Return(testCase.thirdVote)
		t.Equal(testCase.decision, t.service.IsLoggedOut(t.context, webSession))
	}
}

func (t *SecurityServiceTestSuite) TestIsGranted() {
	testCases := []serviceTestCase{
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessDenied,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyAffirmative,
			allowIfAllAbstain: true,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessGranted,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessDenied,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyConsensus,
			allowIfAllAbstain: true,
			decision:          true,
		},
		{
			firstVote:         voter.AccessGranted,
			secondVote:        voter.AccessGranted,
			thirdVote:         voter.AccessDenied,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessGranted,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          true,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: false,
			decision:          false,
		},
		{
			firstVote:         voter.AccessAbstained,
			secondVote:        voter.AccessAbstained,
			thirdVote:         voter.AccessAbstained,
			voterStrategy:     VoterStrategyUnanimous,
			allowIfAllAbstain: true,
			decision:          true,
		},
	}

	for _, testCase := range testCases {
		webSession := web.EmptySession()
		t.service.voterStrategy = testCase.voterStrategy
		t.service.allowIfAllAbstain = testCase.allowIfAllAbstain
		t.firstVoter.On("Vote", t.context, webSession, "SomePermission", nil).Once().Return(testCase.firstVote)
		t.secondVoter.On("Vote", t.context, webSession, "SomePermission", nil).Once().Return(testCase.secondVote)
		t.thirdVoter.On("Vote", t.context, webSession, "SomePermission", nil).Once().Return(testCase.thirdVote)
		t.Equal(testCase.decision, t.service.IsGranted(t.context, webSession, "SomePermission", nil))
	}
}
