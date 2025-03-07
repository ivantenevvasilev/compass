// Code generated by mockery. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// Fetcher is an autogenerated mock type for the Fetcher type
type Fetcher struct {
	mock.Mock
}

// FetchOnDemand provides a mock function with given fields: _a0, parentTenant
func (_m *Fetcher) FetchOnDemand(_a0 string, parentTenant string) error {
	ret := _m.Called(_a0, parentTenant)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, parentTenant)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFetcher creates a new instance of Fetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFetcher(t interface {
	mock.TestingT
	Cleanup(func())
}) *Fetcher {
	mock := &Fetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
