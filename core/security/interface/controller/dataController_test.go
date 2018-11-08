package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"net/http"

	"flamingo.me/flamingo/core/security/application/mocks"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
)

type (
	DataControllerTestSuite struct {
		suite.Suite

		controller      *DataController
		securityService *mocks.SecurityService

		context context.Context
		request *web.Request
		session *sessions.Session
	}
)

func TestDataControllerTestSuite(t *testing.T) {
	suite.Run(t, &DataControllerTestSuite{})
}

func (t *DataControllerTestSuite) SetupSuite() {
	t.context = context.Background()
	t.session = sessions.NewSession(nil, "")
	t.request = web.RequestFromRequest(&http.Request{}, web.NewSession(t.session))
	t.request.LoadParams(map[string]string{
		"permission": "SomePermission",
	})
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
	t.securityService.On("IsLoggedIn", t.context, t.session).Return(true).Once()
	result := t.controller.IsLoggedIn(t.context, t.request)
	isLoggedIn, ok := result.(bool)
	t.True(isLoggedIn)
	t.True(ok)
}

func (t *DataControllerTestSuite) TestIsLoggedOut() {
	t.securityService.On("IsLoggedOut", t.context, t.session).Return(true).Once()
	result := t.controller.IsLoggedOut(t.context, t.request)
	isLoggedOut, ok := result.(bool)
	t.True(isLoggedOut)
	t.True(ok)
}

func (t *DataControllerTestSuite) TestIsGranted() {
	t.securityService.On("IsGranted", t.context, t.session, "SomePermission", nil).Return(true).Once()
	result := t.controller.IsGranted(t.context, t.request)
	isGranted, ok := result.(bool)
	t.True(isGranted)
	t.True(ok)
}
