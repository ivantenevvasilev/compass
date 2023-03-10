// Code generated by mockery. DO NOT EDIT.

package automock

import (
	formationtemplateconstraintreferences "github.com/kyma-incubator/compass/components/director/internal/domain/formationtemplateconstraintreferences"
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// ConstraintReferenceConverter is an autogenerated mock type for the constraintReferenceConverter type
type ConstraintReferenceConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: e
func (_m *ConstraintReferenceConverter) FromEntity(e *formationtemplateconstraintreferences.Entity) *model.FormationTemplateConstraintReference {
	ret := _m.Called(e)

	var r0 *model.FormationTemplateConstraintReference
	if rf, ok := ret.Get(0).(func(*formationtemplateconstraintreferences.Entity) *model.FormationTemplateConstraintReference); ok {
		r0 = rf(e)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationTemplateConstraintReference)
		}
	}

	return r0
}

// ToEntity provides a mock function with given fields: in
func (_m *ConstraintReferenceConverter) ToEntity(in *model.FormationTemplateConstraintReference) *formationtemplateconstraintreferences.Entity {
	ret := _m.Called(in)

	var r0 *formationtemplateconstraintreferences.Entity
	if rf, ok := ret.Get(0).(func(*model.FormationTemplateConstraintReference) *formationtemplateconstraintreferences.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*formationtemplateconstraintreferences.Entity)
		}
	}

	return r0
}

// ToGraphql provides a mock function with given fields: in
func (_m *ConstraintReferenceConverter) ToGraphql(in *model.FormationTemplateConstraintReference) *graphql.ConstraintReference {
	ret := _m.Called(in)

	var r0 *graphql.ConstraintReference
	if rf, ok := ret.Get(0).(func(*model.FormationTemplateConstraintReference) *graphql.ConstraintReference); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.ConstraintReference)
		}
	}

	return r0
}

// ToModel provides a mock function with given fields: in
func (_m *ConstraintReferenceConverter) ToModel(in *graphql.ConstraintReference) *model.FormationTemplateConstraintReference {
	ret := _m.Called(in)

	var r0 *model.FormationTemplateConstraintReference
	if rf, ok := ret.Get(0).(func(*graphql.ConstraintReference) *model.FormationTemplateConstraintReference); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationTemplateConstraintReference)
		}
	}

	return r0
}

type mockConstructorTestingTNewConstraintReferenceConverter interface {
	mock.TestingT
	Cleanup(func())
}

// NewConstraintReferenceConverter creates a new instance of ConstraintReferenceConverter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConstraintReferenceConverter(t mockConstructorTestingTNewConstraintReferenceConverter) *ConstraintReferenceConverter {
	mock := &ConstraintReferenceConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
