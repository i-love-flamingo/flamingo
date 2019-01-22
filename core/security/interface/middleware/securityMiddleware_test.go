package middleware

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	applicationMocks "flamingo.me/flamingo/v3/core/security/application/mocks"
	interfaceMocks "flamingo.me/flamingo/v3/core/security/interface/middleware/mocks"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	SecurityMiddlewareTestSuite struct {
		suite.Suite

		middleware       *SecurityMiddleware
		securityService  *applicationMocks.SecurityService
		redirectURLMaker *interfaceMocks.RedirectUrlMaker

		context    context.Context
		request    *web.Request
		response   web.Result
		session    *sessions.Session
		webSession *web.Session
		action     web.Action
	}
)

func TestSecurityMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &SecurityMiddlewareTestSuite{})
}

func (t *SecurityMiddlewareTestSuite) SetupSuite() {
	t.context = context.Background()
	t.webSession = web.EmptySession()
	t.request = web.CreateRequest(&http.Request{
		URL: &url.URL{
			Path: "/referrer",
		},
		Header: http.Header{
			"Referer": []string{"/http-referrer"},
		},
	}, t.webSession)
	t.request.Params["permission"] = "SomePermission"
	t.response = &web.Response{
		Status: http.StatusOK,
	}
	t.action = func(ctx context.Context, req *web.Request) web.Result {
		return t.response
	}
}

func (t *SecurityMiddlewareTestSuite) SetupTest() {
	t.securityService = &applicationMocks.SecurityService{}
	t.redirectURLMaker = &interfaceMocks.RedirectUrlMaker{}
	t.middleware = &SecurityMiddleware{}
	t.middleware.Inject(&web.Responder{}, t.securityService, t.redirectURLMaker, flamingo.NullLogger{}, &struct {
		LoginPathHandler              string `inject:"config:security.loginPath.handler"`
		LoginPathRedirectStrategy     string `inject:"config:security.loginPath.redirectStrategy"`
		LoginPathRedirectPath         string `inject:"config:security.loginPath.redirectPath"`
		AuthenticatedHomepageStrategy string `inject:"config:security.authenticatedHomepage.strategy"`
		AuthenticatedHomepagePath     string `inject:"config:security.authenticatedHomepage.path"`
		EventLogging                  bool   `inject:"config:security.eventLogging"`
	}{
		LoginPathRedirectPath:     "/home",
		AuthenticatedHomepagePath: "/authenticated",
	})
}

func (t *SecurityMiddlewareTestSuite) TearDownTest() {
	t.securityService.AssertExpectations(t.T())
	t.securityService = nil
	t.redirectURLMaker.AssertExpectations(t.T())
	t.redirectURLMaker = nil
	t.middleware = nil
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedIn_ForbiddenWithReferrer() {
	redirectURL, err := url.Parse("/referrer")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedIn(t.action)
	t.middleware.loginPathRedirectStrategy = ReferrerRedirectStrategy
	t.securityService.On("IsLoggedIn", t.context, t.webSession).Return(false).Once()
	t.redirectURLMaker.On("URL", t.context, "/referrer").Return(redirectURL, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.RouteRedirectResponse)

	t.True(ok)
	t.Equal(map[string]string{
		"redirecturl": "/referrer",
	}, response.Data)
	t.Equal("auth.login", response.To)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedIn_ForbiddenWithPath() {
	redirectURL, err := url.Parse("/home")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedIn(t.action)
	t.middleware.loginPathRedirectStrategy = PathRedirectStrategy
	t.securityService.On("IsLoggedIn", t.context, t.webSession).Return(false).Once()
	t.redirectURLMaker.On("URL", t.context, "/home").Return(redirectURL, nil).Once()

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

func (t *SecurityMiddlewareTestSuite) TestRedirectToLoginFallback() {
	redirectURL, err := url.Parse("/home")
	t.NoError(err)

	t.middleware.loginPathRedirectStrategy = PathRedirectStrategy
	t.redirectURLMaker.On("URL", t.context, "/home").Return(redirectURL, nil).Once()

	result := t.middleware.RedirectToLoginFallback(t.context, t.request)
	response, ok := result.(*web.RouteRedirectResponse)

	t.True(ok)
	t.Equal(map[string]string{
		"redirecturl": "/home",
	}, response.Data)
	t.Equal("auth.login", response.To)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedOut_ForbiddenWithReferrer() {
	redirectURL, err := url.Parse("/http-referrer")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedOut(t.action)
	t.middleware.authenticatedHomepageStrategy = ReferrerRedirectStrategy
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(false).Once()
	t.redirectURLMaker.On("URL", t.context, "/http-referrer").Return(redirectURL, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.URLRedirectResponse)

	t.True(ok)
	t.Equal(redirectURL, response.URL)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedOut_ForbiddenWithPath() {
	redirectURL, err := url.Parse("/authenticated")
	t.NoError(err)

	action := t.middleware.HandleIfLoggedOut(t.action)
	t.middleware.authenticatedHomepageStrategy = PathRedirectStrategy
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(false).Once()
	t.redirectURLMaker.On("URL", t.context, "/authenticated").Return(redirectURL, nil).Once()

	result := action(t.context, t.request)
	response, ok := result.(*web.URLRedirectResponse)

	t.True(ok)
	t.Equal(redirectURL, response.URL)
}

func (t *SecurityMiddlewareTestSuite) TestHandleIfLoggedOut_Allowed() {
	action := t.middleware.HandleIfLoggedOut(t.action)
	t.securityService.On("IsLoggedOut", t.context, t.webSession).Return(true).Once()

	result := action(t.context, t.request)
	t.Exactly(t.response, result)
}

func (t *SecurityMiddlewareTestSuite) TestRedirectToLogoutFallback() {
	redirectURL, err := url.Parse("/authenticated")
	t.NoError(err)

	t.middleware.authenticatedHomepageStrategy = PathRedirectStrategy
	t.redirectURLMaker.On("URL", t.context, "/authenticated").Return(redirectURL, nil).Once()

	result := t.middleware.RedirectToLogoutFallback(t.context, t.request)
	response, ok := result.(*web.URLRedirectResponse)

	t.True(ok)
	t.Equal(redirectURL, response.URL)
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
