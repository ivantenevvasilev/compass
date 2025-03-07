// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// ApplicationService is an autogenerated mock type for the ApplicationService type
type ApplicationService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *ApplicationService) Create(ctx context.Context, in model.ApplicationRegisterInput) (string, error) {
	ret := _m.Called(ctx, in)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationRegisterInput) (string, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationRegisterInput) string); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ApplicationRegisterInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateFromTemplate provides a mock function with given fields: ctx, in, appTemplateID
func (_m *ApplicationService) CreateFromTemplate(ctx context.Context, in model.ApplicationRegisterInput, appTemplateID *string) (string, error) {
	ret := _m.Called(ctx, in, appTemplateID)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationRegisterInput, *string) (string, error)); ok {
		return rf(ctx, in, appTemplateID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationRegisterInput, *string) string); ok {
		r0 = rf(ctx, in, appTemplateID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.ApplicationRegisterInput, *string) error); ok {
		r1 = rf(ctx, in, appTemplateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *ApplicationService) Get(ctx context.Context, id string) (*model.Application, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Application
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Application, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Application); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Application)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApplicationService creates a new instance of ApplicationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApplicationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApplicationService {
	mock := &ApplicationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
