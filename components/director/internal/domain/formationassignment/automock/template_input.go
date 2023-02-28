// Code generated by mockery. DO NOT EDIT.

package automock

import (
	http "net/http"
	testing "testing"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	webhook "github.com/kyma-incubator/compass/components/director/pkg/webhook"
)

// TemplateInput is an autogenerated mock type for the templateInput type
type TemplateInput struct {
	mock.Mock
}

// Clone provides a mock function with given fields:
func (_m *TemplateInput) Clone() webhook.FormationAssignmentTemplateInput {
	ret := _m.Called()

	var r0 webhook.FormationAssignmentTemplateInput
	if rf, ok := ret.Get(0).(func() webhook.FormationAssignmentTemplateInput); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(webhook.FormationAssignmentTemplateInput)
		}
	}

	return r0
}

// GetAssignment provides a mock function with given fields:
func (_m *TemplateInput) GetAssignment() *model.FormationAssignment {
	ret := _m.Called()

	var r0 *model.FormationAssignment
	if rf, ok := ret.Get(0).(func() *model.FormationAssignment); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationAssignment)
		}
	}

	return r0
}

// GetParticipantsIDs provides a mock function with given fields:
func (_m *TemplateInput) GetParticipantsIDs() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetReverseAssignment provides a mock function with given fields:
func (_m *TemplateInput) GetReverseAssignment() *model.FormationAssignment {
	ret := _m.Called()

	var r0 *model.FormationAssignment
	if rf, ok := ret.Get(0).(func() *model.FormationAssignment); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationAssignment)
		}
	}

	return r0
}

// ParseHeadersTemplate provides a mock function with given fields: tmpl
func (_m *TemplateInput) ParseHeadersTemplate(tmpl *string) (http.Header, error) {
	ret := _m.Called(tmpl)

	var r0 http.Header
	if rf, ok := ret.Get(0).(func(*string) http.Header); ok {
		r0 = rf(tmpl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Header)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(tmpl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseInputTemplate provides a mock function with given fields: tmpl
func (_m *TemplateInput) ParseInputTemplate(tmpl *string) ([]byte, error) {
	ret := _m.Called(tmpl)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(*string) []byte); ok {
		r0 = rf(tmpl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(tmpl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseURLTemplate provides a mock function with given fields: tmpl
func (_m *TemplateInput) ParseURLTemplate(tmpl *string) (*webhook.URL, error) {
	ret := _m.Called(tmpl)

	var r0 *webhook.URL
	if rf, ok := ret.Get(0).(func(*string) *webhook.URL); ok {
		r0 = rf(tmpl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*webhook.URL)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(tmpl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetAssignment provides a mock function with given fields: _a0
func (_m *TemplateInput) SetAssignment(_a0 *model.FormationAssignment) {
	_m.Called(_a0)
}

// SetReverseAssignment provides a mock function with given fields: _a0
func (_m *TemplateInput) SetReverseAssignment(_a0 *model.FormationAssignment) {
	_m.Called(_a0)
}

// NewTemplateInput creates a new instance of TemplateInput. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTemplateInput(t testing.TB) *TemplateInput {
	mock := &TemplateInput{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
