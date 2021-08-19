/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/authentication"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tenantPathParamValue = "tenant"
)

type Tenant struct {
	TenantId   string `json:"tenantId"`
	CustomerId string `json:"customerId"`
}

func TestOnboardingHandler(t *testing.T) {
	t.Run("Success with tenant and customerID", func(t *testing.T) {
		// GIVEN
		providedTenant := Tenant{
			TenantId:   uuid.New().String(),
			CustomerId: uuid.New().String(),
		}

		oldTenantState, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// WHEN
		endpoint := strings.Replace(config.HandlerEndpoint, fmt.Sprintf("{%s}", config.TenantPathParam), tenantPathParamValue, 1)
		url := config.TenantFetcherURL + config.RootAPI + endpoint

		byteTenant, err := json.Marshal(providedTenant)
		require.NoError(t, err)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteTenant))
		require.NoError(t, err)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

		httpClient := http.DefaultClient
		httpClient.Timeout = 15 * time.Second

		response, err := httpClient.Do(request)
		require.NoError(t, err)

		tenants, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// THEN
		require.Equal(t, len(oldTenantState)+2, len(tenants))
		require.Equal(t, http.StatusOK, response.StatusCode)
		containsTenantWithTenantID(providedTenant.TenantId, tenants)
		containsTenantWithTenantID(providedTenant.CustomerId, tenants)
	})

	t.Run("Success with only tenant", func(t *testing.T) {
		// GIVEN
		providedTenant := Tenant{
			TenantId: uuid.New().String(),
		}

		oldTenantState, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// WHEN
		endpoint := strings.Replace(config.HandlerEndpoint, fmt.Sprintf("{%s}", config.TenantPathParam), tenantPathParamValue, 1)
		url := config.TenantFetcherURL + config.RootAPI + endpoint

		byteTenant, err := json.Marshal(providedTenant)
		require.NoError(t, err)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteTenant))
		require.NoError(t, err)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

		httpClient := http.DefaultClient
		httpClient.Timeout = 15 * time.Second

		response, err := httpClient.Do(request)
		require.NoError(t, err)

		tenants, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// THEN
		require.Equal(t, len(oldTenantState)+1, len(tenants))
		require.Equal(t, http.StatusOK, response.StatusCode)
		containsTenantWithTenantID(providedTenant.TenantId, tenants)
	})

	t.Run("Should not fail when tenant already exists", func(t *testing.T) {
		//GIVEN
		providedTenant := Tenant{
			TenantId:   uuid.New().String(),
			CustomerId: uuid.New().String(),
		}

		//WHEN
		endpoint := strings.Replace(config.HandlerEndpoint, fmt.Sprintf("{%s}", config.TenantPathParam), tenantPathParamValue, 1)
		url := config.TenantFetcherURL + config.RootAPI + endpoint

		byteTenant, err := json.Marshal(providedTenant)
		require.NoError(t, err)

		var response *http.Response
		for i := 0; i < 10; i++ {
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteTenant))
			require.NoError(t, err)
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

			httpClient := http.DefaultClient
			httpClient.Timeout = 15 * time.Second

			response, err = httpClient.Do(request)
			require.NoError(t, err)
		}

		// THEN
		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("Should not add already existing tenants", func(t *testing.T) {
		//GIVEN
		providedTenant := Tenant{
			TenantId:   uuid.New().String(),
			CustomerId: uuid.New().String(),
		}

		oldTenantState, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		//WHEN
		endpoint := strings.Replace(config.HandlerEndpoint, fmt.Sprintf("{%s}", config.TenantPathParam), tenantPathParamValue, 1)
		url := config.TenantFetcherURL + config.RootAPI + endpoint

		var response *http.Response
		for i := 0; i < 10; i++ {
			byteTenant, err := json.Marshal(providedTenant)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteTenant))
			require.NoError(t, err)
			request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

			httpClient := http.DefaultClient
			httpClient.Timeout = 15 * time.Second

			response, err = httpClient.Do(request)
			require.NoError(t, err)
		}

		tenants, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState)+2, len(tenants))
		require.Equal(t, http.StatusOK, response.StatusCode)
		containsTenantWithTenantID(providedTenant.TenantId, tenants)
		containsTenantWithTenantID(providedTenant.CustomerId, tenants)
	})

	t.Run("Should fail when no tenantID is provided", func(t *testing.T) {
		providedTenant := Tenant{
			CustomerId: uuid.New().String(),
		}

		oldTenantState, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		endpoint := strings.Replace(config.HandlerEndpoint, fmt.Sprintf("{%s}", config.TenantPathParam), tenantPathParamValue, 1)
		url := config.TenantFetcherURL + config.RootAPI + endpoint

		byteTenant, err := json.Marshal(providedTenant)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteTenant))
		require.NoError(t, err)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

		httpClient := http.DefaultClient
		httpClient.Timeout = 15 * time.Second

		response, err := httpClient.Do(request)
		require.NoError(t, err)

		tenants, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState), len(tenants))
		require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
}

func TestDecommissioningHandler(t *testing.T) {
	t.Run("Success noop", func(t *testing.T) {
		// GIVEN
		providedTenant := Tenant{
			TenantId: uuid.New().String(),
		}

		// WHEN
		endpoint := strings.Replace(config.HandlerEndpoint, fmt.Sprintf("{%s}", config.TenantPathParam), tenantPathParamValue, 1)
		url := config.TenantFetcherURL + config.RootAPI + endpoint

		// Add test tenant
		byteTenant, err := json.Marshal(providedTenant)
		require.NoError(t, err)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(byteTenant))
		require.NoError(t, err)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

		httpClient := http.DefaultClient
		httpClient.Timeout = 15 * time.Second

		response, err := httpClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		// Initial state
		oldTenantState, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		request, err = http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(byteTenant))
		require.NoError(t, err)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authentication.CreateNotSingedToken(t)))

		response, err = httpClient.Do(request)
		require.NoError(t, err)

		newTenantState, err := fixtures.GetTenants(config.DirectorUrl, config.Tenant)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState), len(newTenantState))
		require.Equal(t, http.StatusOK, response.StatusCode)
	})
}

func containsTenantWithTenantID(tenantID string, tenants []*graphql.Tenant) bool {
	for _, t := range tenants {
		if t.ID == tenantID {
			return true
		}
	}
	return false
}
