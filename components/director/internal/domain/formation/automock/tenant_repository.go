// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// TenantRepository is an autogenerated mock type for the tenantRepository type
type TenantRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, id
func (_m *TenantRepository) Get(ctx context.Context, id string) (*model.BusinessTenantMapping, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.BusinessTenantMapping
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.BusinessTenantMapping); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.BusinessTenantMapping)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCustomerIDParentRecursively provides a mock function with given fields: ctx, tenant
func (_m *TenantRepository) GetCustomerIDParentRecursively(ctx context.Context, tenant string) (string, error) {
	ret := _m.Called(ctx, tenant)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, tenant)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTenantRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTenantRepository creates a new instance of TenantRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantRepository(t mockConstructorTestingTNewTenantRepository) *TenantRepository {
	mock := &TenantRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
