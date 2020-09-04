// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// PackageService is an autogenerated mock type for the PackageService type
type PackageService struct {
	mock.Mock
}

// AssociateBundle provides a mock function with given fields: ctx, id, bundleID
func (_m *PackageService) AssociateBundle(ctx context.Context, id string, bundleID string) error {
	ret := _m.Called(ctx, id, bundleID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, id, bundleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, applicationID, in
func (_m *PackageService) Create(ctx context.Context, applicationID string, in model.PackageCreateInput) (string, error) {
	ret := _m.Called(ctx, applicationID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, model.PackageCreateInput) string); ok {
		r0 = rf(ctx, applicationID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.PackageCreateInput) error); ok {
		r1 = rf(ctx, applicationID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *PackageService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *PackageService) Get(ctx context.Context, id string) (*model.Package, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Package
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Package); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
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

// Update provides a mock function with given fields: ctx, id, in
func (_m *PackageService) Update(ctx context.Context, id string, in model.PackageUpdateInput) error {
	ret := _m.Called(ctx, id, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.PackageUpdateInput) error); ok {
		r0 = rf(ctx, id, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
