// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FormationRepository is an autogenerated mock type for the formationRepository type
type FormationRepository struct {
	mock.Mock
}

// ListByFormationNames provides a mock function with given fields: ctx, formationNames, tenantID
func (_m *FormationRepository) ListByFormationNames(ctx context.Context, formationNames []string, tenantID string) ([]*model.Formation, error) {
	ret := _m.Called(ctx, formationNames, tenantID)

	var r0 []*model.Formation
	if rf, ok := ret.Get(0).(func(context.Context, []string, string) []*model.Formation); ok {
		r0 = rf(ctx, formationNames, tenantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Formation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string, string) error); ok {
		r1 = rf(ctx, formationNames, tenantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFormationRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewFormationRepository creates a new instance of FormationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationRepository(t mockConstructorTestingTNewFormationRepository) *FormationRepository {
	mock := &FormationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
