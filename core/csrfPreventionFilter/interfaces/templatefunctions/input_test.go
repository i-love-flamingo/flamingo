package templatefunctions

import (
	"context"
	"testing"

	applicationMocks "flamingo.me/flamingo/core/csrfPreventionFilter/application/mocks"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/session"
	"flamingo.me/flamingo/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	CsrfInputFuncTestSuite struct {
		suite.Suite

		csrfFunc *CsrfInputFunc

		service *applicationMocks.Service

		context context.Context
		session *web.Session
	}
)

func TestCsrfInputFuncTestSuite(t *testing.T) {
	suite.Run(t, &CsrfInputFuncTestSuite{})
}

func (t *CsrfInputFuncTestSuite) SetupSuite() {
	t.session = web.NewSession(nil)
	t.context = session.Context(context.Background(), t.session)
}

func (t *CsrfInputFuncTestSuite) SetupTest() {
	t.service = &applicationMocks.Service{}

	t.csrfFunc = &CsrfInputFunc{}
	t.csrfFunc.Inject(t.service, flamingo.NullLogger{})
}

func (t *CsrfInputFuncTestSuite) TearDown() {
	t.service.AssertExpectations(t.T())
	t.service = nil
	t.csrfFunc = nil
}

func (t *CsrfInputFuncTestSuite) TestFunc() {
	t.service.On("Generate", t.session).Return("token").Once()

	function := t.csrfFunc.Func(t.context)
	csrfFunc, ok := function.(func() interface{})
	t.True(ok)
	content := csrfFunc()
	t.Equal(`<input type="hidden" name="csrftoken" value="token" />`, content)
}
