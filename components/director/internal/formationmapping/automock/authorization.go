// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Authorization is an autogenerated mock type for the Authorization type
type Authorization struct {
	mock.Mock
}

// IsAuthorized provides a mock function with given fields: ctx, formationAssignmentID
func (_m *Authorization) IsAuthorized(ctx context.Context, formationAssignmentID string) (bool, error, int) {
	ret := _m.Called(ctx, formationAssignmentID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, formationAssignmentID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, formationAssignmentID)
	} else {
		r1 = ret.Error(1)
	}

	var r2 int
	if rf, ok := ret.Get(2).(func(context.Context, string) int); ok {
		r2 = rf(ctx, formationAssignmentID)
	} else {
		r2 = ret.Get(2).(int)
	}

	return r0, r1, r2
}

// NewAuthorization creates a new instance of Authorization. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthorization(t testing.TB) *Authorization {
	mock := &Authorization{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
