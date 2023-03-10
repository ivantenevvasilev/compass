// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	formationconstraint "github.com/kyma-incubator/compass/components/director/pkg/formationconstraint"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FormationConstraintSvc is an autogenerated mock type for the formationConstraintSvc type
type FormationConstraintSvc struct {
	mock.Mock
}

// ListMatchingConstraints provides a mock function with given fields: ctx, formationTemplateID, location, details
func (_m *FormationConstraintSvc) ListMatchingConstraints(ctx context.Context, formationTemplateID string, location formationconstraint.JoinPointLocation, details formationconstraint.MatchingDetails) ([]*model.FormationConstraint, error) {
	ret := _m.Called(ctx, formationTemplateID, location, details)

	var r0 []*model.FormationConstraint
	if rf, ok := ret.Get(0).(func(context.Context, string, formationconstraint.JoinPointLocation, formationconstraint.MatchingDetails) []*model.FormationConstraint); ok {
		r0 = rf(ctx, formationTemplateID, location, details)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FormationConstraint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, formationconstraint.JoinPointLocation, formationconstraint.MatchingDetails) error); ok {
		r1 = rf(ctx, formationTemplateID, location, details)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFormationConstraintSvc interface {
	mock.TestingT
	Cleanup(func())
}

// NewFormationConstraintSvc creates a new instance of FormationConstraintSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationConstraintSvc(t mockConstructorTestingTNewFormationConstraintSvc) *FormationConstraintSvc {
	mock := &FormationConstraintSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
