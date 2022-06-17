// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	labelfilter "github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// ApplicationTemplateRepository is an autogenerated mock type for the ApplicationTemplateRepository type
type ApplicationTemplateRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *ApplicationTemplateRepository) Create(ctx context.Context, item model.ApplicationTemplate) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationTemplate) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *ApplicationTemplateRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, id
func (_m *ApplicationTemplateRepository) Exists(ctx context.Context, id string) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *ApplicationTemplateRepository) Get(ctx context.Context, id string) (*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.ApplicationTemplate); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByFilters provides a mock function with given fields: ctx, filter
func (_m *ApplicationTemplateRepository) GetByFilters(ctx context.Context, filter []*labelfilter.LabelFilter) (*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, filter)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter) *model.ApplicationTemplate); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: ctx, id
func (_m *ApplicationTemplateRepository) GetByName(ctx context.Context, id string) (*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.ApplicationTemplate); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, filter, pageSize, cursor
func (_m *ApplicationTemplateRepository) List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (model.ApplicationTemplatePage, error) {
	ret := _m.Called(ctx, filter, pageSize, cursor)

	var r0 model.ApplicationTemplatePage
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter, int, string) model.ApplicationTemplatePage); ok {
		r0 = rf(ctx, filter, pageSize, cursor)
	} else {
		r0 = ret.Get(0).(model.ApplicationTemplatePage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter, int, string) error); ok {
		r1 = rf(ctx, filter, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *ApplicationTemplateRepository) Update(ctx context.Context, _a1 model.ApplicationTemplate) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationTemplate) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewApplicationTemplateRepository creates a new instance of ApplicationTemplateRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplicationTemplateRepository(t testing.TB) *ApplicationTemplateRepository {
	mock := &ApplicationTemplateRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
