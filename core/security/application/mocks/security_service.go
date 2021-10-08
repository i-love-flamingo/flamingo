// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	web "flamingo.me/flamingo/v3/framework/web"
	mock "github.com/stretchr/testify/mock"
)

// SecurityService is an autogenerated mock type for the SecurityService type
type SecurityService struct {
	mock.Mock
}

// IsGranted provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *SecurityService) IsGranted(_a0 context.Context, _a1 *web.Session, _a2 string, _a3 interface{}) bool {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *web.Session, string, interface{}) bool); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsLoggedIn provides a mock function with given fields: _a0, _a1
func (_m *SecurityService) IsLoggedIn(_a0 context.Context, _a1 *web.Session) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *web.Session) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsLoggedOut provides a mock function with given fields: _a0, _a1
func (_m *SecurityService) IsLoggedOut(_a0 context.Context, _a1 *web.Session) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *web.Session) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
