package handlers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/kyma-incubator/compass/components/ias-adapter/internal/errors"

	"github.com/kyma-incubator/compass/components/ias-adapter/internal/api/internal/handlers/automock"
	"github.com/kyma-incubator/compass/components/ias-adapter/internal/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Tenant Mapping Handler", func() {
	var tenantMapping = types.TenantMapping{
		FormationID: "2d933ae2-10c4-4d6f-b4d4-5e1553e4ff05",
		ReceiverTenant: types.ReceiverTenant{
			ApplicationURL: "localhost",
		},
		AssignedTenants: []types.AssignedTenant{
			{
				UCLApplicationID: "2d933ae2-10c4-4d6f-b4d4-5e1553e4ff05",
				LocalTenantID:    "2d933ae2-10c4-4d6f-b4d4-5e1553e4ff05",
				Operation:        types.OperationAssign,
				Parameters: types.AssignedTenantParameters{
					ClientID: "clientID",
				},
				ReverseAssignmentState: "",
			},
		},
	}
	BeforeEach(func() {
		tenantMapping.AssignedTenants[0].ReverseAssignmentState = ""
	})
	When("Tenant mapping cannot be decoded", func() {
		It("Should fail with 400", func() {
			handler := TenantMappingsHandler{Service: &automock.TenantMappingsService{}}

			w, ctx := createTestRequest("unprocessable body")

			handler.Patch(ctx)
			responseBody, err := io.ReadAll(w.Body)
			Expect(err).Error().ToNot(HaveOccurred())
			Expect(responseBody).To(ContainSubstring(url.QueryEscape("failed to decode tenant mapping body")))
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
	When("Tenant mapping is invalid", func() {
		It("Should fail with 400", func() {
			handler := TenantMappingsHandler{Service: &automock.TenantMappingsService{}}

			w, ctx := createTestRequest(`{"assignedTenants":[{"configuration": ""}]}`)

			handler.Patch(ctx)
			responseBody, err := io.ReadAll(w.Body)
			Expect(err).Error().ToNot(HaveOccurred())
			Expect(responseBody).To(ContainSubstring(url.QueryEscape("tenant mapping body is invalid")))
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
	When("Reverse assignment state is neither INITIAL nor READY", func() {
		BeforeEach(func() {
			tenantMapping.AssignedTenants[0].ReverseAssignmentState = "CREATE_ERROR"
		})
		It("Should fail with 422 CONFIG_PENDING", func() {
			service := &automock.TenantMappingsService{}
			service.On("ProcessTenantMapping", mock.Anything, mock.Anything).Return(nil)
			handler := TenantMappingsHandler{Service: service}

			w, ctx := createTestRequest(tenantMapping)

			handler.Patch(ctx)
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})
	When("Consumed APIs cannot be updated", func() {
		BeforeEach(func() {
			tenantMapping.AssignedTenants[0].ReverseAssignmentState = types.StateInitial
		})
		It("Should fail with 500", func() {
			service := &automock.TenantMappingsService{}
			service.On("ProcessTenantMapping", mock.Anything, mock.Anything).Return(errors.New("error"))
			handler := TenantMappingsHandler{Service: service}

			w, ctx := createTestRequest(tenantMapping)

			handler.Patch(ctx)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})
	})
	When("Consumed APIs cannot be updated due to not found IAS application", func() {
		When("Operation is Assign", func() {
			BeforeEach(func() {
				tenantMapping.AssignedTenants[0].ReverseAssignmentState = types.StateInitial
				tenantMapping.AssignedTenants[0].Operation = types.OperationAssign
			})
			It("Should return 404", func() {
				service := &automock.TenantMappingsService{}
				service.On("ProcessTenantMapping", mock.Anything, mock.Anything).Return(errors.Newf("could not process tenant mapping: %w", errors.IASApplicationNotFound))
				handler := TenantMappingsHandler{Service: service}

				w, ctx := createTestRequest(tenantMapping)

				handler.Patch(ctx)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
	When("Consumed APIs are successfully updated", func() {
		BeforeEach(func() {
			tenantMapping.AssignedTenants[0].ReverseAssignmentState = types.StateReady
		})
		It("Should return 200", func() {
			service := &automock.TenantMappingsService{}
			service.On("ProcessTenantMapping", mock.Anything, mock.Anything).Return(nil)
			handler := TenantMappingsHandler{Service: service}

			w, ctx := createTestRequest(tenantMapping)

			handler.Patch(ctx)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})
