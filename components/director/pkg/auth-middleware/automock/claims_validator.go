// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	idtokenclaims "github.com/kyma-incubator/compass/components/director/pkg/idtokenclaims"
	mock "github.com/stretchr/testify/mock"
)

// ClaimsValidator is an autogenerated mock type for the ClaimsValidator type
type ClaimsValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: _a0, _a1
func (_m *ClaimsValidator) Validate(_a0 context.Context, _a1 idtokenclaims.Claims) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, idtokenclaims.Claims) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewClaimsValidator interface {
	mock.TestingT
	Cleanup(func())
}

// NewClaimsValidator creates a new instance of ClaimsValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClaimsValidator(t mockConstructorTestingTNewClaimsValidator) *ClaimsValidator {
	mock := &ClaimsValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
