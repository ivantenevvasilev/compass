// Code generated by mockery. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromGraphql provides a mock function with given fields: id, in
func (_m *Converter) FromGraphql(id string, in graphql.CertificateSubjectMappingInput) *model.CertSubjectMapping {
	ret := _m.Called(id, in)

	var r0 *model.CertSubjectMapping
	if rf, ok := ret.Get(0).(func(string, graphql.CertificateSubjectMappingInput) *model.CertSubjectMapping); ok {
		r0 = rf(id, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.CertSubjectMapping)
		}
	}

	return r0
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *Converter) MultipleToGraphQL(in []*model.CertSubjectMapping) []*graphql.CertificateSubjectMapping {
	ret := _m.Called(in)

	var r0 []*graphql.CertificateSubjectMapping
	if rf, ok := ret.Get(0).(func([]*model.CertSubjectMapping) []*graphql.CertificateSubjectMapping); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.CertificateSubjectMapping)
		}
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: in
func (_m *Converter) ToGraphQL(in *model.CertSubjectMapping) *graphql.CertificateSubjectMapping {
	ret := _m.Called(in)

	var r0 *graphql.CertificateSubjectMapping
	if rf, ok := ret.Get(0).(func(*model.CertSubjectMapping) *graphql.CertificateSubjectMapping); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.CertificateSubjectMapping)
		}
	}

	return r0
}

// NewConverter creates a new instance of Converter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewConverter(t testing.TB) *Converter {
	mock := &Converter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
