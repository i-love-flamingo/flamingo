// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PermissionSet is an autogenerated mock type for the PermissionSet type
type PermissionSet struct {
	mock.Mock
}

// Permissions provides a mock function with given fields:
func (_m *PermissionSet) Permissions() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
