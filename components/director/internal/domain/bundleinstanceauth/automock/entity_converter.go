// Code generated by mockery v2.10.4. DO NOT EDIT.

package automock

import (
	bundleinstanceauth "github.com/kyma-incubator/compass/components/director/internal/domain/bundleinstanceauth"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *bundleinstanceauth.Entity) (*model.BundleInstanceAuth, error) {
	ret := _m.Called(entity)

	var r0 *model.BundleInstanceAuth
	if rf, ok := ret.Get(0).(func(*bundleinstanceauth.Entity) *model.BundleInstanceAuth); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.BundleInstanceAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*bundleinstanceauth.Entity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.BundleInstanceAuth) (*bundleinstanceauth.Entity, error) {
	ret := _m.Called(in)

	var r0 *bundleinstanceauth.Entity
	if rf, ok := ret.Get(0).(func(*model.BundleInstanceAuth) *bundleinstanceauth.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bundleinstanceauth.Entity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.BundleInstanceAuth) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
