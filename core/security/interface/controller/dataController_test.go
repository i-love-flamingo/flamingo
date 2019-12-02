package controller

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/core/security/application/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	DataControllerTestSuite struct {
		suite.Suite

		controller      *DataController
		securityService *mocks.SecurityService

		context    context.Context
		request    *web.Request
		params     web.RequestParams
		webSession *web.Session
	}
)

func TestDataControllerTestSuite(t *testing.T) {
	suite.Run(t, &DataControllerTestSuite{})
}

func (t *DataControllerTestSuite) SetupSuite() {
	t.context = context.Background()
	t.webSession = web.EmptySession()
	t.request = web.CreateRequest(nil, t.webSession)
	t.params = web.RequestParams{
		"permission": "SomePermission",
	}

}

func (t *DataControllerTestSuite) SetupTest() {
	t.securityService = &mocks.SecurityService{}
	t.controller = &DataController{}
	t.controller.Inject(t.securityService)
}

func (t *DataControllerTestSuite) TearDownTest() {
	t.securityService.AssertExpectations(t.T())
	t.securityService = nil
	t.controller = nil
}

func (t *DataControllerTestSuite) TestIsLoggedIn() {
	t.securityService.On("IsLoggedIn", t.context, t.webSession).Return(true).Once()
	result := t.controller.IsLoggedIn(t.context, t.request, nil)
	isLoggedIn, ok := result.(bool)
	t.True(isLoggedIn)
	t.True(ok)
}

func (t *DataControllerTestSuite) TestIsLoggedOut() {
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(true).Once()
	result := t.controller.IsLoggedOut(t.context, t.request, nil)
	isLoggedOut, ok := result.(bool)
	t.True(isLoggedOut)
	t.True(ok)
}

func (t *DataControllerTestSuite) TestIsGranted() {
	t.securityService.On("IsGranted", t.context, t.webSession, "SomePermission", nil).Return(true).Once()
	result := t.controller.IsGranted(t.context, t.request, t.params)
	isGranted, ok := result.(bool)
	t.True(isGranted)
	t.True(ok)
}
