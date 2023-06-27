// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// EventDefinitionService is an autogenerated mock type for the EventDefinitionService type
type EventDefinitionService struct {
	mock.Mock
}

// GetForApplication provides a mock function with given fields: ctx, id, appID
func (_m *EventDefinitionService) GetForApplication(ctx context.Context, id string, appID string) (*model.EventDefinition, error) {
	ret := _m.Called(ctx, id, appID)

	var r0 *model.EventDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.EventDefinition); ok {
		r0 = rf(ctx, id, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.EventDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEventDefinitionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewEventDefinitionService creates a new instance of EventDefinitionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEventDefinitionService(t mockConstructorTestingTNewEventDefinitionService) *EventDefinitionService {
	mock := &EventDefinitionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
