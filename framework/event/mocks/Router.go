// Code generated by mockery v1.0.0
package mocks

import context "context"
import event "go.aoe.com/flamingo/framework/event"
import mock "github.com/stretchr/testify/mock"

// Router is an autogenerated mock type for the Router type
type Router struct {
	mock.Mock
}

// Dispatch provides a mock function with given fields: ctx, _a1
func (_m *Router) Dispatch(ctx context.Context, _a1 event.Event) {
	_m.Called(ctx, _a1)
}
