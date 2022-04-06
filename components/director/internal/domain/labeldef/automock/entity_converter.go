// Code generated by mockery v2.10.4. DO NOT EDIT.

package automock

import (
	labeldef "github.com/kyma-incubator/compass/components/director/internal/domain/labeldef"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: in
func (_m *EntityConverter) FromEntity(in labeldef.Entity) (model.LabelDefinition, error) {
	ret := _m.Called(in)

	var r0 model.LabelDefinition
	if rf, ok := ret.Get(0).(func(labeldef.Entity) model.LabelDefinition); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(model.LabelDefinition)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(labeldef.Entity) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in model.LabelDefinition) (labeldef.Entity, error) {
	ret := _m.Called(in)

	var r0 labeldef.Entity
	if rf, ok := ret.Get(0).(func(model.LabelDefinition) labeldef.Entity); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(labeldef.Entity)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.LabelDefinition) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
