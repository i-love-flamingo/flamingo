package middleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"

	"net/url"

	applicationMocks "flamingo.me/flamingo/core/security/application/mocks"
	interfaceMocks "flamingo.me/flamingo/core/security/interface/middleware/mocks"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	SecurityMiddlewareTestSuite struct {
		suite.Suite

		middleware       *SecurityMiddleware
		securityService  *applicationMocks.SecurityService
		redirectUrlMaker *interfaceMocks.RedirectUrlMaker

		context    context.Context
		request    *web.Request
		response   web.Response
		session    *sessions.Session
		webSession *web.Session
		action     router.Action
	}
)

func TestSecurityMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &SecurityMiddlewareTestSuite{})
}

func (t *SecurityMiddlewareTestSuite) SetupSuite() {
	t.context = context.Background()
	t.session = sessions.NewSession(nil, "")
	t.webSession = web.NewSession(t.session)
	t.request = web.RequestFromRequest(&http.Request{
		URL: &url.URL{
			Path: "/referrer",
		},
		Header: http.Header{
			"Referer": []string{"/http-referrer"},
		},
	}, web.NewSession(t.session))
	t.request.LoadParams(map[string]string{
		"permission": "SomePermission",
	})
	t.response = &web.HTTPResponse{
		Status: http.StatusOK,
	}
	t.action = func(ctx context.Context, req *web.Request) web.Response {
		return t.response
	}
}

func (t *SecurityMiddlewareTestSuite) SetupTest() {
	t.securityService = &applicationMocks.SecurityService{}
	t.redirectUrlMaker = &interfaceMocks.RedirectUrlMaker{}
	t.middleware = &SecurityMiddleware{}
	t.middleware.Inject(&web.Responder{}, t.securityService, t.redirectUrlMaker, flamingo.NullLogger{}, &struct {
		LoginPathHandler              string `inject:"config:security.loginPath.handler"`
		LoginPathRedirectStrategy     string `inject:"config:security.loginPath.redirectStrategy"`
		LoginPathRedirectPath         string `inject:"config:security.loginPath.redirectPath"`
		AuthenticatedHomepageStrategy string `inject:"config:security.authenticatedHomepage.strategy"`
		AuthenticatedHomepagePath     string `inject:"config:security.authenticatedHomepage.path"`
	}{
		LoginPathRedirectPath:     "/home",
		AuthenticatedHomepagePath: "/authenticated",
	})
}

func (t *SecurityMiddlewareTestSuite) TearDownTest() {
	t.securityService.AssertExpectations(t.T())
	t.securityService = nil
	t.redirectUrlMaker.AssertExpectations(t.T())
	t.redirectUrlMaker = nil
	t.middleware = nil
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedIn_ForbiddenWithReferrer() {
	redirectUrl, err := url.Parse("/referrer")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedIn(t.action)
	t.middleware.loginPathRedirectStrategy = ReferrerRedirectStrategy
	t.securityService.On("IsLoggedIn", t.context, t.webSession).Return(false).Once()
	t.redirectUrlMaker.On("URL", t.context, "/referrer").Return(redirectUrl, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.RouteRedirectResponse)

	t.True(ok)
	t.Equal(map[string]string{
		"redirecturl": "/referrer",
	}, response.Data)
	t.Equal("auth.login", response.To)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedIn_ForbiddenWithPath() {
	redirectUrl, err := url.Parse("/home")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedIn(t.action)
	t.middleware.loginPathRedirectStrategy = PathRedirectStrategy
	t.securityService.On("IsLoggedIn", t.context, t.webSession).Return(false).Once()
	t.redirectUrlMaker.On("URL", t.context, "/home").Return(redirectUrl, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.RouteRedirectResponse)

	t.True(ok)
	t.Equal(map[string]string{
		"redirecturl": "/home",
	}, response.Data)
	t.Equal("auth.login", response.To)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedIn_Allowed() {
	action := t.middleware.HandleIfLoggedIn(t.action)
	t.securityService.On("IsLoggedIn", t.context, t.webSession).Return(true).Once()

	result := action(t.context, t.request)
	t.Exactly(t.response, result)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedOut_ForbiddenWithReferrer() {
	redirectUrl, err := url.Parse("/http-referrer")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedOut(t.action)
	t.middleware.authenticatedHomepageStrategy = ReferrerRedirectStrategy
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(false).Once()
	t.redirectUrlMaker.On("URL", t.context, "/http-referrer").Return(redirectUrl, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.URLRedirectResponse)

	t.True(ok)
	t.Equal(redirectUrl, response.URL)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedOut_ForbiddenWithPath() {
	redirectUrl, err := url.Parse("/authenticated")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedOut(t.action)
	t.middleware.authenticatedHomepageStrategy = PathRedirectStrategy
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(false).Once()
	t.redirectUrlMaker.On("URL", t.context, "/authenticated").Return(redirectUrl, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.URLRedirectResponse)

	t.True(ok)
	t.Equal(redirectUrl, response.URL)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedOut_Allowed() {
	action := t.middleware.HandleIfLoggedOut(t.action)
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(true).Once()

	result := action(t.context, t.request)
	t.Exactly(t.response, result)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfGranted_Forbidden() {
	action := t.middleware.HandleIfGranted(t.action, "SomePermission")
	t.securityService.On("IsGranted", t.context, t.webSession, "SomePermission", nil).Return(false).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.ServerErrorResponse)

	t.True(ok)
	t.Equal(uint(http.StatusForbidden), response.Status)
	t.Equal("Permission SomePermission for path /referrer.", response.Error.Error())
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfGranted_Allowed() {
	action := t.middleware.HandleIfGranted(t.action, "SomePermission")
	t.securityService.On("IsGranted", t.context, t.webSession, "SomePermission", nil).Return(true).Once()

	result := action(t.context, t.request)
	t.Exactly(t.response, result)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfNotGranted_Forbidden() {
	action := t.middleware.HandleIfNotGranted(t.action, "SomePermission")
	t.securityService.On("IsGranted", t.context, t.webSession, "SomePermission", nil).Return(true).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.ServerErrorResponse)

	t.True(ok)
	t.Equal(uint(http.StatusForbidden), response.Status)
	t.Equal("Permission SomePermission for path /referrer.", response.Error.Error())
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfNotGranted_Allowed() {
	action := t.middleware.HandleIfNotGranted(t.action, "SomePermission")
	t.securityService.On("IsGranted", t.context, t.webSession, "SomePermission", nil).Return(false).Once()

	result := action(t.context, t.request)
	t.Exactly(t.response, result)
}
