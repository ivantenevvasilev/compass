// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	webhook "github.com/kyma-incubator/compass/components/director/pkg/webhook"
)

// DataInputBuilder is an autogenerated mock type for the DataInputBuilder type
type DataInputBuilder struct {
	mock.Mock
}

// PrepareApplicationAndAppTemplateWithLabels provides a mock function with given fields: ctx, tenant, appID
func (_m *DataInputBuilder) PrepareApplicationAndAppTemplateWithLabels(ctx context.Context, tenant string, appID string) (*webhook.ApplicationWithLabels, *webhook.ApplicationTemplateWithLabels, error) {
	ret := _m.Called(ctx, tenant, appID)

	var r0 *webhook.ApplicationWithLabels
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *webhook.ApplicationWithLabels); ok {
		r0 = rf(ctx, tenant, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhook.ApplicationWithLabels)
		}
	}

	var r1 *webhook.ApplicationTemplateWithLabels
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *webhook.ApplicationTemplateWithLabels); ok {
		r1 = rf(ctx, tenant, appID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*webhook.ApplicationTemplateWithLabels)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, tenant, appID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PrepareApplicationMappingsInFormation provides a mock function with given fields: ctx, tenant, scenario
func (_m *DataInputBuilder) PrepareApplicationMappingsInFormation(ctx context.Context, tenant string, scenario string) (map[string]*webhook.ApplicationWithLabels, map[string]*webhook.ApplicationTemplateWithLabels, error) {
	ret := _m.Called(ctx, tenant, scenario)

	var r0 map[string]*webhook.ApplicationWithLabels
	if rf, ok := ret.Get(0).(func(context.Context, string, string) map[string]*webhook.ApplicationWithLabels); ok {
		r0 = rf(ctx, tenant, scenario)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*webhook.ApplicationWithLabels)
		}
	}

	var r1 map[string]*webhook.ApplicationTemplateWithLabels
	if rf, ok := ret.Get(1).(func(context.Context, string, string) map[string]*webhook.ApplicationTemplateWithLabels); ok {
		r1 = rf(ctx, tenant, scenario)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]*webhook.ApplicationTemplateWithLabels)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, tenant, scenario)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PrepareRuntimeAndRuntimeContextWithLabels provides a mock function with given fields: ctx, tenant, runtimeID
func (_m *DataInputBuilder) PrepareRuntimeAndRuntimeContextWithLabels(ctx context.Context, tenant string, runtimeID string) (*webhook.RuntimeWithLabels, *webhook.RuntimeContextWithLabels, error) {
	ret := _m.Called(ctx, tenant, runtimeID)

	var r0 *webhook.RuntimeWithLabels
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *webhook.RuntimeWithLabels); ok {
		r0 = rf(ctx, tenant, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhook.RuntimeWithLabels)
		}
	}

	var r1 *webhook.RuntimeContextWithLabels
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *webhook.RuntimeContextWithLabels); ok {
		r1 = rf(ctx, tenant, runtimeID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*webhook.RuntimeContextWithLabels)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, tenant, runtimeID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PrepareRuntimeContextWithLabels provides a mock function with given fields: ctx, tenant, runtimeCtxID
func (_m *DataInputBuilder) PrepareRuntimeContextWithLabels(ctx context.Context, tenant string, runtimeCtxID string) (*webhook.RuntimeContextWithLabels, error) {
	ret := _m.Called(ctx, tenant, runtimeCtxID)

	var r0 *webhook.RuntimeContextWithLabels
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *webhook.RuntimeContextWithLabels); ok {
		r0 = rf(ctx, tenant, runtimeCtxID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhook.RuntimeContextWithLabels)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, runtimeCtxID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareRuntimeWithLabels provides a mock function with given fields: ctx, tenant, runtimeID
func (_m *DataInputBuilder) PrepareRuntimeWithLabels(ctx context.Context, tenant string, runtimeID string) (*webhook.RuntimeWithLabels, error) {
	ret := _m.Called(ctx, tenant, runtimeID)

	var r0 *webhook.RuntimeWithLabels
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *webhook.RuntimeWithLabels); ok {
		r0 = rf(ctx, tenant, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhook.RuntimeWithLabels)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareRuntimesAndRuntimeContextsMappingsInFormation provides a mock function with given fields: ctx, tenant, scenario
func (_m *DataInputBuilder) PrepareRuntimesAndRuntimeContextsMappingsInFormation(ctx context.Context, tenant string, scenario string) (map[string]*webhook.RuntimeWithLabels, map[string]*webhook.RuntimeContextWithLabels, error) {
	ret := _m.Called(ctx, tenant, scenario)

	var r0 map[string]*webhook.RuntimeWithLabels
	if rf, ok := ret.Get(0).(func(context.Context, string, string) map[string]*webhook.RuntimeWithLabels); ok {
		r0 = rf(ctx, tenant, scenario)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*webhook.RuntimeWithLabels)
		}
	}

	var r1 map[string]*webhook.RuntimeContextWithLabels
	if rf, ok := ret.Get(1).(func(context.Context, string, string) map[string]*webhook.RuntimeContextWithLabels); ok {
		r1 = rf(ctx, tenant, scenario)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]*webhook.RuntimeContextWithLabels)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, tenant, scenario)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewDataInputBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewDataInputBuilder creates a new instance of DataInputBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDataInputBuilder(t mockConstructorTestingTNewDataInputBuilder) *DataInputBuilder {
	mock := &DataInputBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
