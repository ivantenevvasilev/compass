// Code generated by mockery v2.10.4. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromGraphQL provides a mock function with given fields: i
func (_m *Converter) FromGraphQL(i graphql.FormationInput) model.Formation {
	ret := _m.Called(i)

	var r0 model.Formation
	if rf, ok := ret.Get(0).(func(graphql.FormationInput) model.Formation); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Get(0).(model.Formation)
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: i
func (_m *Converter) ToGraphQL(i *model.Formation) *graphql.Formation {
	ret := _m.Called(i)

	var r0 *graphql.Formation
	if rf, ok := ret.Get(0).(func(*model.Formation) *graphql.Formation); ok {
		r0 = rf(i)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.Formation)
		}
	}

	return r0
}
