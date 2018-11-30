package interfaces

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	applicationMocks "flamingo.me/flamingo/core/csrfPreventionFilter/application/mocks"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	CsrfMiddlewareTestSuite struct {
		suite.Suite

		middleware *CsrfMiddleware
		service    *applicationMocks.Service

		action     router.Action
		context    context.Context
		webRequest *web.Request
	}
)

func TestCsrfMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, &CsrfMiddlewareTestSuite{})
}

func (t *CsrfMiddlewareTestSuite) SetupSuite() {
	t.context = context.Background()
	t.action = func(ctx context.Context, req *web.Request) web.Response {
		return &web.BasicResponse{}
	}
	t.webRequest = web.RequestFromRequest(nil, nil)
}

func (t *CsrfMiddlewareTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.middleware = &CsrfMiddleware{}
	t.middleware.Inject(&web.Responder{}, t.service)
}

func (t *CsrfMiddlewareTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
}

func (t *CsrfMiddlewareTestSuite) TestMiddleware_WrongToken() {
	t.service.On("IsValid", t.webRequest).Return(false).Once()

	handler := t.middleware.Secured(t.action)
	response := handler(t.context, t.webRequest)
	forbidden, ok := response.(*web.ServerErrorResponse)
	t.True(ok)
	t.Equal(uint(http.StatusForbidden), forbidden.Status)
}

func (t *CsrfMiddlewareTestSuite) TestMiddleware_Success() {
	t.service.On("IsValid", t.webRequest).Return(true).Once()

	handler := t.middleware.Secured(t.action)
	response := handler(t.context, t.webRequest)
	t.IsType(&web.BasicResponse{}, response)
}
