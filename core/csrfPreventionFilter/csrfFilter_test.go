package csrfPreventionFilter

import (
	"net/http"
	"testing"

	"flamingo.me/flamingo/core/csrfPreventionFilter/mocks"
	"flamingo.me/flamingo/framework/router"
	webmocks "flamingo.me/flamingo/framework/web/mocks"
	respondermocks "flamingo.me/flamingo/framework/web/responder/mocks"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type (
	testCsrfFilterData struct {
		name              string
		request           *http.Request
		ignore            bool
		sessionCalls      int
		values            interface{}
		requestNonce      string
		requestError      error
		errorAwareError   bool
		errorMsg          string
		checkDeletedNonce bool
	}
)

func TestCsrfFilter_Filter(t *testing.T) {
	for _, data := range []testCsrfFilterData{
		{
			name:              "positive example",
			request:           &http.Request{Method: "POST"},
			ignore:            false,
			sessionCalls:      3,
			values:            []string{"17"},
			requestNonce:      "17",
			requestError:      nil,
			checkDeletedNonce: true,
		},
		{
			name:         "ignoring CSRF",
			request:      &http.Request{Method: "POST"},
			ignore:       true,
			sessionCalls: 0,
		},
		{
			name:            "no sessionNonce",
			request:         &http.Request{Method: "POST"},
			sessionCalls:    1,
			errorAwareError: true,
			values:          nil,
			errorMsg:        `session hasn't the key "csrfNonces"`,
		},
		{
			name:            "sessionNonce isn't a list",
			request:         &http.Request{Method: "POST"},
			sessionCalls:    1,
			errorAwareError: true,
			values:          "test",
			errorMsg:        `the session key "csrfNonces" isn't a list'"`,
		},
		{
			name:            "no nonce in request",
			request:         &http.Request{Method: "POST"},
			sessionCalls:    1,
			values:          []string{"17"},
			requestError:    errors.New("form value not found"),
			errorAwareError: true,
			errorMsg:        "form value not found",
		},
		{
			name:            "not same nonce",
			request:         &http.Request{Method: "POST"},
			sessionCalls:    1,
			values:          []string{"17"},
			requestNonce:    "42",
			errorAwareError: true,
			errorMsg:        "session doesn't contain the csrf-nonce of the request",
		},
	} {
		t.Run(data.name, func(t *testing.T) {
			session := new(sessions.Session)
			session.Values = make(map[interface{}]interface{})
			if data.values != nil {
				session.Values[csrfNonces] = data.values
			}

			ctxMock := new(webmocks.Context)
			ctxMock.On("Request").Once().Return(data.request)
			ctxMock.On("Session").Return(session)
			ctxMock.On("Form1", "csrf_token").Return(data.requestNonce, data.requestError)

			lastFilter := new(mocks.Filter)
			filterChain := &router.FilterChain{Filters: []router.Filter{lastFilter}}
			if data.ignore {
				controllerMock := &mocks.ControllerOptionAware{}
				controllerMock.On("CheckOption", router.ControllerOption("csrf.ignore")).Once().Return(true)
				filterChain.Controller = controllerMock
			}
			lastFilter.On("Filter", ctxMock, nil, filterChain).Return(nil)

			csrfFilter := new(csrfFilter)
			if data.errorAwareError {
				errorAwareMock := new(respondermocks.ErrorAware)
				errorAwareMock.On("Error", ctxMock, mock.MatchedBy(errorMsg(data.errorMsg))).Return(nil)
				csrfFilter.ErrorAware = errorAwareMock
			}

			//csrfFilter.Filter(ctxMock, web.RequestFromRequest(data.request, session), nil, filterChain)

			//ctxMock.AssertNumberOfCalls(t, "Request", 1)
			//ctxMock.AssertNumberOfCalls(t, "Session", data.sessionCalls)

			//if data.checkDeletedNonce {
			//	assert.False(t, contains(ctxMock.Session().Values[csrfNonces].([]string), data.requestNonce))
			//}
		})
	}
}

func errorMsg(msg string) func(error) bool {
	return func(err error) bool { return err.Error() == msg }
}
