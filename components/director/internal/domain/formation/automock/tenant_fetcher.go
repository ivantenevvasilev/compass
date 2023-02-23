// Code generated by mockery. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// TenantFetcher is an autogenerated mock type for the TenantFetcher type
type TenantFetcher struct {
	mock.Mock
}

// FetchOnDemand provides a mock function with given fields: tenant, parentTenant
func (_m *TenantFetcher) FetchOnDemand(tenant string, parentTenant string) error {
	ret := _m.Called(tenant, parentTenant)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(tenant, parentTenant)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewTenantFetcherT interface {
	mock.TestingT
	Cleanup(func())
}

// NewTenantFetcher creates a new instance of TenantFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantFetcher(t NewTenantFetcherT) *TenantFetcher {
	mock := &TenantFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
