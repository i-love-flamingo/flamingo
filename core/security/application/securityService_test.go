package application

import (
	"context"
	"testing"

	roleMocks "flamingo.me/flamingo/v3/core/security/application/role/mocks"
	"flamingo.me/flamingo/v3/core/security/application/voter"
	voterMocks "flamingo.me/flamingo/v3/core/security/application/voter/mocks"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	SecurityServiceTestSuite struct {
		suite.Suite

		service     *SecurityServiceImpl
		firstVoter  *voterMocks.SecurityVoter
		secondVoter *voterMocks.SecurityVoter
		thirdVoter  *voterMocks.SecurityVoter
		roleService *roleMocks.Service

		context context.Context
	}

	serviceTestCase struct {
		firstVote         voter.AccessDecision
		secondVote        voter.AccessDecision
		thirdVote         voter.AccessDecision
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
	t.firstVoter = &voterMocks.SecurityVoter{}
	t.secondVoter = &voterMocks.SecurityVoter{}
	t.thirdVoter = &voterMocks.SecurityVoter{}
	voters := []voter.SecurityVoter{
		t.firstVoter,
		t.secondVoter,
		t.thirdVoter,
	}
	t.roleService = &roleMocks.Service{}
	t.service = &SecurityServiceImpl{}
	t.service.Inject(voters, t.roleService, &struct {
		VoterStrategy     string `inject:"config:core.security.roles.voters.strategy"`
		AllowIfAllAbstain bool   `inject:"config:core.security.roles.voters.allowIfAllAbstain"`
	}{})
}

func (t *SecurityServiceTestSuite) TearDownTest() {
	t.firstVoter.AssertExpectations(t.T())
	t.firstVoter = nil
	t.secondVoter.AssertExpectations(t.T())
	t.secondVoter = nil
	t.thirdVoter.AssertExpectations(t.T())
	t.thirdVoter = nil
	t.roleService.AssertExpectations(t.T())
	t.roleService = nil
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
		t.firstVoter.On("Vote", []string{}, domain.PermissionAuthorized, nil).Return(testCase.firstVote).Once()
		t.secondVoter.On("Vote", []string{}, domain.PermissionAuthorized, nil).Return(testCase.secondVote).Once()
		t.thirdVoter.On("Vote", []string{}, domain.PermissionAuthorized, nil).Return(testCase.thirdVote).Once()
		t.roleService.On("AllPermissions", t.context, webSession).Return([]string{}).Once()
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
		t.firstVoter.On("Vote", []string{}, domain.PermissionAuthorized, nil).Return(testCase.firstVote).Once()
		t.secondVoter.On("Vote", []string{}, domain.PermissionAuthorized, nil).Return(testCase.secondVote).Once()
		t.thirdVoter.On("Vote", []string{}, domain.PermissionAuthorized, nil).Return(testCase.thirdVote).Once()
		t.roleService.On("AllPermissions", t.context, webSession).Return([]string{}).Once()
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
		t.firstVoter.On("Vote", []string{}, "SomePermission", nil).Return(testCase.firstVote).Once()
		t.secondVoter.On("Vote", []string{}, "SomePermission", nil).Return(testCase.secondVote).Once()
		t.thirdVoter.On("Vote", []string{}, "SomePermission", nil).Return(testCase.thirdVote).Once()
		t.roleService.On("AllPermissions", t.context, webSession).Return([]string{}).Once()
		t.Equal(testCase.decision, t.service.IsGranted(t.context, webSession, "SomePermission", nil))
	}
}
