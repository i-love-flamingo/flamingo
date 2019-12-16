// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	voter "flamingo.me/flamingo/v3/core/security/application/voter"
	mock "github.com/stretchr/testify/mock"
)

// SecurityVoter is an autogenerated mock type for the SecurityVoter type
type SecurityVoter struct {
	mock.Mock
}

// Vote provides a mock function with given fields: allAssignedPermissions, desiredPermission, forObject
func (_m *SecurityVoter) Vote(allAssignedPermissions []string, desiredPermission string, forObject interface{}) voter.AccessDecision {
	ret := _m.Called(allAssignedPermissions, desiredPermission, forObject)

	var r0 voter.AccessDecision
	if rf, ok := ret.Get(0).(func([]string, string, interface{}) voter.AccessDecision); ok {
		r0 = rf(allAssignedPermissions, desiredPermission, forObject)
	} else {
		r0 = ret.Get(0).(voter.AccessDecision)
	}

	return r0
}
