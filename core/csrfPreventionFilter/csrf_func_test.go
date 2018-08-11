package csrfPreventionFilter

import (
	"testing"

	"flamingo.me/flamingo/core/csrfPreventionFilter/mocks"
	webmocks "flamingo.me/flamingo/framework/web/mocks"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

type (
	testCsrfData struct {
		name   string
		values interface{}
	}
)

func TestCsrfFuncFunc(t *testing.T) {
	t.Skip("to be updated")

	for _, data := range []testCsrfData{
		{
			name:   "empty session value",
			values: nil,
		},
		{
			name:   `session value "csrfNonces" contains a list of values`,
			values: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
		{
			name:   `session value "csrfNonces" isn't a list of csrfNonces`,
			values: "test",
		},
	} {
		t.Run(data.name, func(t *testing.T) {
			nonce := "17"
			mockNonceGenerator := new(mocks.NonceGenerator)
			mockNonceGenerator.On("GenerateNonce").Once().Return(nonce)

			session := new(sessions.Session)
			session.Values = make(map[interface{}]interface{})
			if data.values != nil {
				session.Values[csrfNonces] = data.values
			}

			ctx := new(webmocks.Context)
			ctx.On("Session").Twice().Return(session)

			csrfFunc := new(CsrfFunc)
			csrfFunc.Generator = mockNonceGenerator
			csrfFunc.TokenLimit = 10
			result := csrfFunc.Func(ctx).(func() interface{})()

			assert.Equal(t, nonce, result)
			assert.True(t, contains(session.Values[csrfNonces].([]string), nonce))

			mockNonceGenerator.AssertNumberOfCalls(t, "GenerateNonce", 1)
			ctx.AssertNumberOfCalls(t, "Session", 2)
		})
	}
}
