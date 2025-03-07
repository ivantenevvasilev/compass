// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	webhook "github.com/kyma-incubator/compass/components/director/pkg/webhook"

	webhookclient "github.com/kyma-incubator/compass/components/director/pkg/webhook_client"
)

// NotificationsGenerator is an autogenerated mock type for the notificationsGenerator type
type NotificationsGenerator struct {
	mock.Mock
}

// GenerateFormationLifecycleNotifications provides a mock function with given fields: ctx, formationTemplateWebhooks, tenantID, _a3, formationTemplateName, formationTemplateID, formationOperation, customerTenantContext
func (_m *NotificationsGenerator) GenerateFormationLifecycleNotifications(ctx context.Context, formationTemplateWebhooks []*model.Webhook, tenantID string, _a3 *model.Formation, formationTemplateName string, formationTemplateID string, formationOperation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationNotificationRequest, error) {
	ret := _m.Called(ctx, formationTemplateWebhooks, tenantID, _a3, formationTemplateName, formationTemplateID, formationOperation, customerTenantContext)

	var r0 []*webhookclient.FormationNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Webhook, string, *model.Formation, string, string, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationNotificationRequest, error)); ok {
		return rf(ctx, formationTemplateWebhooks, tenantID, _a3, formationTemplateName, formationTemplateID, formationOperation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []*model.Webhook, string, *model.Formation, string, string, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationNotificationRequest); ok {
		r0 = rf(ctx, formationTemplateWebhooks, tenantID, _a3, formationTemplateName, formationTemplateID, formationOperation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []*model.Webhook, string, *model.Formation, string, string, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, formationTemplateWebhooks, tenantID, _a3, formationTemplateName, formationTemplateID, formationOperation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsAboutApplicationsForTheRuntimeContextThatIsAssigned provides a mock function with given fields: ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsAboutApplicationsForTheRuntimeContextThatIsAssigned(ctx context.Context, tenant string, runtimeCtxID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsAboutApplicationsForTheRuntimeThatIsAssigned provides a mock function with given fields: ctx, tenant, runtimeID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsAboutApplicationsForTheRuntimeThatIsAssigned(ctx context.Context, tenant string, runtimeID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsAboutRuntimeAndRuntimeContextForTheApplicationThatIsAssigned provides a mock function with given fields: ctx, tenant, appID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsAboutRuntimeAndRuntimeContextForTheApplicationThatIsAssigned(ctx context.Context, tenant string, appID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, appID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsForApplicationsAboutTheApplicationThatIsAssigned provides a mock function with given fields: ctx, tenant, appID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsForApplicationsAboutTheApplicationThatIsAssigned(ctx context.Context, tenant string, appID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, appID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsForApplicationsAboutTheRuntimeContextThatIsAssigned provides a mock function with given fields: ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsForApplicationsAboutTheRuntimeContextThatIsAssigned(ctx context.Context, tenant string, runtimeCtxID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, runtimeCtxID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsForApplicationsAboutTheRuntimeThatIsAssigned provides a mock function with given fields: ctx, tenant, runtimeID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsForApplicationsAboutTheRuntimeThatIsAssigned(ctx context.Context, tenant string, runtimeID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, runtimeID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateNotificationsForRuntimeAboutTheApplicationThatIsAssigned provides a mock function with given fields: ctx, tenant, appID, _a3, operation, customerTenantContext
func (_m *NotificationsGenerator) GenerateNotificationsForRuntimeAboutTheApplicationThatIsAssigned(ctx context.Context, tenant string, appID string, _a3 *model.Formation, operation model.FormationOperation, customerTenantContext *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error) {
	ret := _m.Called(ctx, tenant, appID, _a3, operation, customerTenantContext)

	var r0 []*webhookclient.FormationAssignmentNotificationRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) ([]*webhookclient.FormationAssignmentNotificationRequest, error)); ok {
		return rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) []*webhookclient.FormationAssignmentNotificationRequest); ok {
		r0 = rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*webhookclient.FormationAssignmentNotificationRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *model.Formation, model.FormationOperation, *webhook.CustomerTenantContext) error); ok {
		r1 = rf(ctx, tenant, appID, _a3, operation, customerTenantContext)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewNotificationsGenerator creates a new instance of NotificationsGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotificationsGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *NotificationsGenerator {
	mock := &NotificationsGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
