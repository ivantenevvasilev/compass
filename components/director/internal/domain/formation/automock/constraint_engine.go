// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	formationconstraint "github.com/kyma-incubator/compass/components/director/pkg/formationconstraint"

	mock "github.com/stretchr/testify/mock"
)

// ConstraintEngine is an autogenerated mock type for the constraintEngine type
type ConstraintEngine struct {
	mock.Mock
}

// EnforceConstraints provides a mock function with given fields: ctx, location, details, formationTemplateID
func (_m *ConstraintEngine) EnforceConstraints(ctx context.Context, location formationconstraint.JoinPointLocation, details formationconstraint.JoinPointDetails, formationTemplateID string) error {
	ret := _m.Called(ctx, location, details, formationTemplateID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, formationconstraint.JoinPointLocation, formationconstraint.JoinPointDetails, string) error); ok {
		r0 = rf(ctx, location, details, formationTemplateID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewConstraintEngineT interface {
	mock.TestingT
	Cleanup(func())
}

// NewConstraintEngine creates a new instance of ConstraintEngine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConstraintEngine(t NewConstraintEngineT) *ConstraintEngine {
	mock := &ConstraintEngine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
