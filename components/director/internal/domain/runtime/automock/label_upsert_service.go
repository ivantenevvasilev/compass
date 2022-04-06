// Code generated by mockery v2.10.4. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// LabelUpsertService is an autogenerated mock type for the labelUpsertService type
type LabelUpsertService struct {
	mock.Mock
}

// UpsertLabel provides a mock function with given fields: ctx, tenant, labelInput
func (_m *LabelUpsertService) UpsertLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) error {
	ret := _m.Called(ctx, tenant, labelInput)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.LabelInput) error); ok {
		r0 = rf(ctx, tenant, labelInput)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertMultipleLabels provides a mock function with given fields: ctx, tenant, objectType, objectID, labels
func (_m *LabelUpsertService) UpsertMultipleLabels(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, labels map[string]interface{}) error {
	ret := _m.Called(ctx, tenant, objectType, objectID, labels)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string, map[string]interface{}) error); ok {
		r0 = rf(ctx, tenant, objectType, objectID, labels)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
