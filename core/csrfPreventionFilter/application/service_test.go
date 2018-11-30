package application

import (
	"testing"

	"net/http"
	"net/url"

	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/suite"
)

type (
	ServiceTestSuite struct {
		suite.Suite

		service *ServiceImpl

		session    *sessions.Session
		webSession *web.Session
		request    *http.Request
		webRequest *web.Request
	}
)

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}

func (t *ServiceTestSuite) SetupTest() {
	t.service = &ServiceImpl{}
	t.service.Inject(flamingo.NullLogger{}, &struct {
		Secret string  `inject:"config:csrf.secret"`
		Ttl    float64 `inject:"config:csrf.ttl"`
	}{
		Secret: "6368616e676520746869732070617373776f726420746f206120736563726574",
		Ttl:    900,
	})

	t.session = sessions.NewSession(nil, "")
	t.webSession = web.NewSession(t.session)
	t.request = &http.Request{}
	t.webRequest = web.RequestFromRequest(t.request, t.webSession)
}

func (t *ServiceTestSuite) TearDown() {
	t.session = nil
	t.webSession = nil
	t.service = nil
	t.request = nil
	t.webRequest = nil
}

func (t *ServiceTestSuite) TestGenerate_WrongKey() {
	t.service.secret = []byte{}
	t.Empty(t.service.Generate(t.webSession))
}

func (t *ServiceTestSuite) TestGenerate_RightKey() {
	t.NotEmpty(t.service.Generate(t.webSession))

	t.session.ID = "1234567890"
	t.NotEmpty(t.service.Generate(t.webSession))
}

func (t *ServiceTestSuite) TestIsValid_GetRequest() {
	t.request.Method = http.MethodGet
	t.True(t.service.IsValid(t.webRequest))
}

func (t *ServiceTestSuite) TestIsValid_MalformedToken() {
	t.request.Method = http.MethodPost
	t.False(t.service.IsValid(t.webRequest))
}

func (t *ServiceTestSuite) TestIsValid_WrongId() {
	t.session.ID = "first"
	token := t.service.Generate(t.webSession)

	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		TokenName: []string{token},
	}

	t.session.ID = "second"
	t.False(t.service.IsValid(t.webRequest))
}

func (t *ServiceTestSuite) TestIsValid_WrongTime() {
	t.service.ttl = -100000000

	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		TokenName: []string{token},
	}

	t.False(t.service.IsValid(t.webRequest))
}

func (t *ServiceTestSuite) TestIsValid_Success() {
	token := t.service.Generate(t.webSession)
	t.request.Method = http.MethodPost
	t.request.PostForm = url.Values{
		TokenName: []string{token},
	}

	t.True(t.service.IsValid(t.webRequest))
}
