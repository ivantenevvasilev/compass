// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	formationassignment "github.com/kyma-incubator/compass/components/director/internal/domain/formationassignment"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FormationAssignmentUpdater is an autogenerated mock type for the formationAssignmentUpdater type
type FormationAssignmentUpdater struct {
	mock.Mock
}

// SetAssignmentToErrorState provides a mock function with given fields: ctx, assignment, errorMessage, errorCode, state, operation
func (_m *FormationAssignmentUpdater) SetAssignmentToErrorState(ctx context.Context, assignment *model.FormationAssignment, errorMessage string, errorCode formationassignment.AssignmentErrorCode, state model.FormationAssignmentState, operation model.FormationOperation) error {
	ret := _m.Called(ctx, assignment, errorMessage, errorCode, state, operation)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FormationAssignment, string, formationassignment.AssignmentErrorCode, model.FormationAssignmentState, model.FormationOperation) error); ok {
		r0 = rf(ctx, assignment, errorMessage, errorCode, state, operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, fa, operation
func (_m *FormationAssignmentUpdater) Update(ctx context.Context, fa *model.FormationAssignment, operation model.FormationOperation) error {
	ret := _m.Called(ctx, fa, operation)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FormationAssignment, model.FormationOperation) error); ok {
		r0 = rf(ctx, fa, operation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewFormationAssignmentUpdater interface {
	mock.TestingT
	Cleanup(func())
}

// NewFormationAssignmentUpdater creates a new instance of FormationAssignmentUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationAssignmentUpdater(t mockConstructorTestingTNewFormationAssignmentUpdater) *FormationAssignmentUpdater {
	mock := &FormationAssignmentUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
