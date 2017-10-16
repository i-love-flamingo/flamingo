package mocks

import (
	"github.com/gorilla/sessions"
	"go.aoe.com/flamingo/framework/web"
)

// Map builds a parameter map for setting up test request contexts
func Map(p ...string) map[string][]string {
	res := make(map[string][]string)
	for i := 0; i < len(p); i += 2 {
		res[p[i]] = append(res[p[i]], p[i+1])
	}
	return res
}

// RequestContext helper for easy mock creation
func RequestContext(params map[string][]string, forms map[string][]string) web.Context {
	mock := new(Context)

	paramsAll := make(map[string]string)
	for k, v := range params {
		mock.On("MustParam1", k).Return(v[0])
		mock.On("Param1", k).Return(v[0], nil)
		paramsAll[k] = v[0]
	}
	mock.On("ParamAll").Return(paramsAll)

	for k, v := range forms {
		mock.On("MustForm1", k).Return(v[0])
		mock.On("Form1", k).Return(v[0], nil)
		mock.On("MustForm", k).Return(v)
		mock.On("Form", k).Return(v, nil)
	}
	mock.On("FormAll").Return(forms)

	session := sessions.NewSession(nil, "testing")
	mock.On("Session").Return(session)

	return mock
}
