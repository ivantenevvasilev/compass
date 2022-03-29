// Code generated by mockery v2.10.0. DO NOT EDIT.

package automock

import (
	runtime "github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	mock "github.com/stretchr/testify/mock"
)

// ExternalSvcCallerProvider is an autogenerated mock type for the ExternalSvcCallerProvider type
type ExternalSvcCallerProvider struct {
	mock.Mock
}

// GetCaller provides a mock function with given fields: _a0, _a1
func (_m *ExternalSvcCallerProvider) GetCaller(_a0 runtime.SelfRegConfig, _a1 string) (runtime.ExternalSvcCaller, error) {
	ret := _m.Called(_a0, _a1)

	var r0 runtime.ExternalSvcCaller
	if rf, ok := ret.Get(0).(func(runtime.SelfRegConfig, string) runtime.ExternalSvcCaller); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(runtime.ExternalSvcCaller)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(runtime.SelfRegConfig, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
