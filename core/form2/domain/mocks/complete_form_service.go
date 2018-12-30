// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import domain "flamingo.me/flamingo/core/form2/domain"
import mock "github.com/stretchr/testify/mock"
import url "net/url"
import web "flamingo.me/flamingo/framework/web"

// CompleteFormService is an autogenerated mock type for the CompleteFormService type
type CompleteFormService struct {
	mock.Mock
}

// Decode provides a mock function with given fields: ctx, req, values, formData
func (_m *CompleteFormService) Decode(ctx context.Context, req *web.Request, values url.Values, formData interface{}) (interface{}, error) {
	ret := _m.Called(ctx, req, values, formData)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, *web.Request, url.Values, interface{}) interface{}); ok {
		r0 = rf(ctx, req, values, formData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *web.Request, url.Values, interface{}) error); ok {
		r1 = rf(ctx, req, values, formData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFormData provides a mock function with given fields: ctx, req
func (_m *CompleteFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
	ret := _m.Called(ctx, req)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, *web.Request) interface{}); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *web.Request) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Validate provides a mock function with given fields: ctx, req, validatorProvider, formData
func (_m *CompleteFormService) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
	ret := _m.Called(ctx, req, validatorProvider, formData)

	var r0 *domain.ValidationInfo
	if rf, ok := ret.Get(0).(func(context.Context, *web.Request, domain.ValidatorProvider, interface{}) *domain.ValidationInfo); ok {
		r0 = rf(ctx, req, validatorProvider, formData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ValidationInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *web.Request, domain.ValidatorProvider, interface{}) error); ok {
		r1 = rf(ctx, req, validatorProvider, formData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
