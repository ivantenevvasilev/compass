// Code generated by mockery. DO NOT EDIT.

package automock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Resources is an autogenerated mock type for the Resources type
type Resources struct {
	mock.Mock
}

// GetType provides a mock function with given fields:
func (_m *Resources) GetType() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetURLPath provides a mock function with given fields:
func (_m *Resources) GetURLPath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewResources creates a new instance of Resources. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewResources(t testing.TB) *Resources {
	mock := &Resources{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
