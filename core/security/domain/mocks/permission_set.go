// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PermissionSet is an autogenerated mock type for the PermissionSet type
type PermissionSet struct {
	mock.Mock
}

type PermissionSet_Expecter struct {
	mock *mock.Mock
}

func (_m *PermissionSet) EXPECT() *PermissionSet_Expecter {
	return &PermissionSet_Expecter{mock: &_m.Mock}
}

// Permissions provides a mock function with given fields:
func (_m *PermissionSet) Permissions() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Permissions")
	}

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

// PermissionSet_Permissions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Permissions'
type PermissionSet_Permissions_Call struct {
	*mock.Call
}

// Permissions is a helper method to define mock.On call
func (_e *PermissionSet_Expecter) Permissions() *PermissionSet_Permissions_Call {
	return &PermissionSet_Permissions_Call{Call: _e.mock.On("Permissions")}
}

func (_c *PermissionSet_Permissions_Call) Run(run func()) *PermissionSet_Permissions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PermissionSet_Permissions_Call) Return(_a0 []string) *PermissionSet_Permissions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PermissionSet_Permissions_Call) RunAndReturn(run func() []string) *PermissionSet_Permissions_Call {
	_c.Call.Return(run)
	return _c
}

// NewPermissionSet creates a new instance of PermissionSet. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPermissionSet(t interface {
	mock.TestingT
	Cleanup(func())
}) *PermissionSet {
	mock := &PermissionSet{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
