// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	web "flamingo.me/flamingo/v3/framework/web"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

type Service_Expecter struct {
	mock *mock.Mock
}

func (_m *Service) EXPECT() *Service_Expecter {
	return &Service_Expecter{mock: &_m.Mock}
}

// AllPermissions provides a mock function with given fields: _a0, _a1
func (_m *Service) AllPermissions(_a0 context.Context, _a1 *web.Session) []string {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for AllPermissions")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, *web.Session) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Service_AllPermissions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AllPermissions'
type Service_AllPermissions_Call struct {
	*mock.Call
}

// AllPermissions is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *web.Session
func (_e *Service_Expecter) AllPermissions(_a0 interface{}, _a1 interface{}) *Service_AllPermissions_Call {
	return &Service_AllPermissions_Call{Call: _e.mock.On("AllPermissions", _a0, _a1)}
}

func (_c *Service_AllPermissions_Call) Run(run func(_a0 context.Context, _a1 *web.Session)) *Service_AllPermissions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*web.Session))
	})
	return _c
}

func (_c *Service_AllPermissions_Call) Return(_a0 []string) *Service_AllPermissions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Service_AllPermissions_Call) RunAndReturn(run func(context.Context, *web.Session) []string) *Service_AllPermissions_Call {
	_c.Call.Return(run)
	return _c
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
