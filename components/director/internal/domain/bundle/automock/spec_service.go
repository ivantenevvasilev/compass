// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	resource "github.com/kyma-incubator/compass/components/director/pkg/resource"
)

// SpecService is an autogenerated mock type for the SpecService type
type SpecService struct {
	mock.Mock
}

// GetByReferenceObjectID provides a mock function with given fields: ctx, resourceType, objectType, objectID
func (_m *SpecService) GetByReferenceObjectID(ctx context.Context, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) (*model.Spec, error) {
	ret := _m.Called(ctx, resourceType, objectType, objectID)

	var r0 *model.Spec
	if rf, ok := ret.Get(0).(func(context.Context, resource.Type, model.SpecReferenceObjectType, string) *model.Spec); ok {
		r0 = rf(ctx, resourceType, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Spec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, resource.Type, model.SpecReferenceObjectType, string) error); ok {
		r1 = rf(ctx, resourceType, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByReferenceObjectIDs provides a mock function with given fields: ctx, objectType, objectIDs
func (_m *SpecService) ListByReferenceObjectIDs(ctx context.Context, objectType model.SpecReferenceObjectType, objectIDs []string) ([]*model.Spec, error) {
	ret := _m.Called(ctx, objectType, objectIDs)

	var r0 []*model.Spec
	if rf, ok := ret.Get(0).(func(context.Context, model.SpecReferenceObjectType, []string) []*model.Spec); ok {
		r0 = rf(ctx, objectType, objectIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Spec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.SpecReferenceObjectType, []string) error); ok {
		r1 = rf(ctx, objectType, objectIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSpecService interface {
	mock.TestingT
	Cleanup(func())
}

// NewSpecService creates a new instance of SpecService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSpecService(t mockConstructorTestingTNewSpecService) *SpecService {
	mock := &SpecService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
