// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"
	sync "sync"

	mock "github.com/stretchr/testify/mock"

	systemfetcher "github.com/kyma-incubator/compass/components/director/internal/systemfetcher"
)

// SystemsAPIClient is an autogenerated mock type for the systemsAPIClient type
type SystemsAPIClient struct {
	mock.Mock
}

// FetchSystemsForTenant provides a mock function with given fields: ctx, tenant, mutex
func (_m *SystemsAPIClient) FetchSystemsForTenant(ctx context.Context, tenant string, mutex *sync.Mutex) ([]systemfetcher.System, error) {
	ret := _m.Called(ctx, tenant, mutex)

	var r0 []systemfetcher.System
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *sync.Mutex) ([]systemfetcher.System, error)); ok {
		return rf(ctx, tenant, mutex)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *sync.Mutex) []systemfetcher.System); ok {
		r0 = rf(ctx, tenant, mutex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]systemfetcher.System)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *sync.Mutex) error); ok {
		r1 = rf(ctx, tenant, mutex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSystemsAPIClient creates a new instance of SystemsAPIClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSystemsAPIClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *SystemsAPIClient {
	mock := &SystemsAPIClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
