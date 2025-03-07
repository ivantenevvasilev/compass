// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	formationassignment "github.com/kyma-incubator/compass/components/director/internal/domain/formationassignment"
	formationconstraint "github.com/kyma-incubator/compass/components/director/pkg/formationconstraint"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	webhookclient "github.com/kyma-incubator/compass/components/director/pkg/webhook_client"
)

// FaNotificationService is an autogenerated mock type for the faNotificationService type
type FaNotificationService struct {
	mock.Mock
}

// GenerateFormationAssignmentNotificationExt provides a mock function with given fields: ctx, faRequestMapping, reverseFaRequestMapping, operation
func (_m *FaNotificationService) GenerateFormationAssignmentNotificationExt(ctx context.Context, faRequestMapping *formationassignment.FormationAssignmentRequestMapping, reverseFaRequestMapping *formationassignment.FormationAssignmentRequestMapping, operation model.FormationOperation) (*webhookclient.FormationAssignmentNotificationRequestExt, error) {
	ret := _m.Called(ctx, faRequestMapping, reverseFaRequestMapping, operation)

	var r0 *webhookclient.FormationAssignmentNotificationRequestExt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *formationassignment.FormationAssignmentRequestMapping, *formationassignment.FormationAssignmentRequestMapping, model.FormationOperation) (*webhookclient.FormationAssignmentNotificationRequestExt, error)); ok {
		return rf(ctx, faRequestMapping, reverseFaRequestMapping, operation)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *formationassignment.FormationAssignmentRequestMapping, *formationassignment.FormationAssignmentRequestMapping, model.FormationOperation) *webhookclient.FormationAssignmentNotificationRequestExt); ok {
		r0 = rf(ctx, faRequestMapping, reverseFaRequestMapping, operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhookclient.FormationAssignmentNotificationRequestExt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *formationassignment.FormationAssignmentRequestMapping, *formationassignment.FormationAssignmentRequestMapping, model.FormationOperation) error); ok {
		r1 = rf(ctx, faRequestMapping, reverseFaRequestMapping, operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareDetailsForNotificationStatusReturned provides a mock function with given fields: ctx, tenantID, fa, operation
func (_m *FaNotificationService) PrepareDetailsForNotificationStatusReturned(ctx context.Context, tenantID string, fa *model.FormationAssignment, operation model.FormationOperation) (*formationconstraint.NotificationStatusReturnedOperationDetails, error) {
	ret := _m.Called(ctx, tenantID, fa, operation)

	var r0 *formationconstraint.NotificationStatusReturnedOperationDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.FormationAssignment, model.FormationOperation) (*formationconstraint.NotificationStatusReturnedOperationDetails, error)); ok {
		return rf(ctx, tenantID, fa, operation)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.FormationAssignment, model.FormationOperation) *formationconstraint.NotificationStatusReturnedOperationDetails); ok {
		r0 = rf(ctx, tenantID, fa, operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*formationconstraint.NotificationStatusReturnedOperationDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *model.FormationAssignment, model.FormationOperation) error); ok {
		r1 = rf(ctx, tenantID, fa, operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFaNotificationService creates a new instance of FaNotificationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFaNotificationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *FaNotificationService {
	mock := &FaNotificationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
