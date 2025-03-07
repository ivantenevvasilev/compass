// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	webhookclient "github.com/kyma-incubator/compass/components/director/pkg/webhook_client"
)

// FormationAssignmentNotificationService is an autogenerated mock type for the FormationAssignmentNotificationService type
type FormationAssignmentNotificationService struct {
	mock.Mock
}

// GenerateFormationAssignmentNotification provides a mock function with given fields: ctx, formationAssignment, operation
func (_m *FormationAssignmentNotificationService) GenerateFormationAssignmentNotification(ctx context.Context, formationAssignment *model.FormationAssignment, operation model.FormationOperation) (*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, formationAssignment, operation)

	var r0 *webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FormationAssignment, model.FormationOperation) (*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, formationAssignment, operation)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.FormationAssignment, model.FormationOperation) *webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, formationAssignment, operation)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.FormationAssignment, model.FormationOperation) error); ok {
		r1 = rf(ctx, formationAssignment, operation)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFormationAssignmentNotificationService creates a new instance of FormationAssignmentNotificationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFormationAssignmentNotificationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *FormationAssignmentNotificationService {
	mock := &FormationAssignmentNotificationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
