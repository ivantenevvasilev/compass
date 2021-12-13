package subscription

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	oauth2 "github.com/kyma-incubator/compass/components/external-services-mock/internal/oauth"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	testErr = errors.New("test error")
	url     = "https://target-url.com"
	token   = "token-value"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	if resp := f(req); resp == nil {
		return nil, testErr
	}
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestHandler_SubscriptionAndDeprovisioning(t *testing.T) {
	// GIVEN
	consumerTenantID := "94764028-8cf8-11ec-9ffc-acde48001122"
	apiPath := fmt.Sprintf("/saas-manager/v1/application/tenants/%s/subscriptions", consumerTenantID)
	reqBody := "{\"subscriptionParams\": {}}"
	emptyTenantConfig := Config{}
	emptyProviderConfig := ProviderConfig{}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	tenantCfg := Config{
		TenantFetcherURL:             "https://tenant-fetcher.com",
		RootAPI:                      "/tenants",
		RegionalHandlerEndpoint:      "/v1/regional/{region}/callback/{tenantId}",
		TenantPathParam:              "tenantId",
		RegionPathParam:              "region",
		SubscriptionProviderID:       "id-value!t12345",
		TenantFetcherFullRegionalURL: "",
		TestConsumerAccountID:        "consumerAccountID",
		TestConsumerSubaccountID:     "consumberSubaccountID",
		TestConsumerTenantID:         "consumerTenantID",
	}

	providerCfg := ProviderConfig{
		TenantIDProperty:               "tenantProperty",
		SubaccountTenantIDProperty:     "subaccountProperty",
		SubdomainProperty:              "subdomainProperty",
		SubscriptionProviderIDProperty: "subscriptionProviderProperty",
	}

	t.Run("Error when missing authorization header", func(t *testing.T) {
		//GIVEN
		subscribeReq, err := http.NewRequest(http.MethodPost, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)
		h := NewHandler(httpClient, emptyTenantConfig, emptyProviderConfig, "")
		r := httptest.NewRecorder()

		//WHEN
		h.Subscription(r, subscribeReq)
		resp := r.Result()

		//THEN
		expectedBody := "{\"error\":\"while executing subscription request: authorization header is required\"}\n"
		assertExpectedResponse(t, resp, expectedBody, http.StatusUnauthorized)
	})

	t.Run("Error when missing token", func(t *testing.T) {
		//GIVEN
		subscribeReq, err := http.NewRequest(http.MethodPost, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)
		subscribeReq.Header.Add(oauth2.AuthorizationHeader, "Bearer ")
		h := NewHandler(httpClient, emptyTenantConfig, emptyProviderConfig, "")
		r := httptest.NewRecorder()

		//WHEN
		h.Subscription(r, subscribeReq)
		resp := r.Result()

		//THEN
		expectedBody := "{\"error\":\"while executing subscription request: token value is required\"}\n"
		assertExpectedResponse(t, resp, expectedBody, http.StatusUnauthorized)
	})

	t.Run("Error when missing tenant path param", func(t *testing.T) {
		//GIVEN
		subReq, err := http.NewRequest(http.MethodPost, url+fmt.Sprintf("/saas-manager/v1/application/tenants/%s/subscriptions", ""), bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)
		subReq.Header.Add(oauth2.AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
		h := NewHandler(httpClient, emptyTenantConfig, emptyProviderConfig, "")
		r := httptest.NewRecorder()

		//WHEN
		h.Subscription(r, subReq)
		resp := r.Result()

		//THEN
		expectedBody := "{\"error\":\"while executing subscription request: parameter [tenant_id] not provided\"}\n"
		assertExpectedResponse(t, resp, expectedBody, http.StatusBadRequest)
	})

	t.Run("Error when executing subscription request", func(t *testing.T) {
		//GIVEN
		subscribeReq, err := http.NewRequest(http.MethodPost, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)
		subscribeReq.Header.Add(oauth2.AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
		subscribeReq = mux.SetURLVars(subscribeReq, map[string]string{"tenant_id": consumerTenantID})

		testErr = errors.New("while executing subscription request to tenant fetcher")
		testClient := NewTestClient(func(req *http.Request) *http.Response {
			return nil
		})

		h := NewHandler(testClient, tenantCfg, providerCfg, "")
		r := httptest.NewRecorder()

		//WHEN
		h.Subscription(r, subscribeReq)
		resp := r.Result()

		//THEN
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NotEmpty(t, body)
		require.Contains(t, string(body), "while executing subscription request")
	})

	t.Run("Error when response code from subscription request is not the expected one", func(t *testing.T) {
		//GIVEN
		subscribeReq, err := http.NewRequest(http.MethodPost, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)
		subscribeReq.Header.Add(oauth2.AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
		subscribeReq = mux.SetURLVars(subscribeReq, map[string]string{"tenant_id": consumerTenantID})

		testClient := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusAccepted,
			}
		})

		h := NewHandler(testClient, tenantCfg, providerCfg, "")
		r := httptest.NewRecorder()

		//WHEN
		h.Subscription(r, subscribeReq)
		resp := r.Result()

		//THEN
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NotEmpty(t, body)
		require.Contains(t, string(body), "while executing subscription request: wrong status code while executing subscription request")
	})

	t.Run("Successfully executed subscription and unsubscription requests", func(t *testing.T) {
		subscribeReq, err := http.NewRequest(http.MethodPost, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)

		unsubscribeReq, err := http.NewRequest(http.MethodDelete, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)

		testCases := []struct {
			Name           string
			Request        *http.Request
			IsSubscription bool
		}{
			{
				Name:           "Successfully executed subscription request",
				Request:        subscribeReq,
				IsSubscription: true,
			},
			{
				Name:           "Successfully executed unsubscription request",
				Request:        unsubscribeReq,
				IsSubscription: false,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.Name, func(t *testing.T) {
				//GIVEN
				req := testCase.Request
				req.Header.Add(oauth2.AuthorizationHeader, fmt.Sprintf("Bearer %s", token))
				req = mux.SetURLVars(req, map[string]string{"tenant_id": consumerTenantID})

				testClient := NewTestClient(func(req *http.Request) *http.Response {
					return &http.Response{
						StatusCode: http.StatusOK,
					}
				})

				h := NewHandler(testClient, tenantCfg, providerCfg, "jobID")
				r := httptest.NewRecorder()

				//WHEN
				if testCase.IsSubscription {
					h.Subscription(r, req)
				} else {
					h.Deprovisioning(r, req)
				}
				resp := r.Result()

				//THEN
				require.Equal(t, http.StatusAccepted, resp.StatusCode)
				require.Equal(t, "/api/v1/jobs/jobID", resp.Header.Get("Location"))
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				require.Empty(t, body)
			})
		}
	})

	t.Run("Error when executing unsubscription request", func(t *testing.T) {
		//GIVEN
		subscribeReq, err := http.NewRequest(http.MethodPost, url+apiPath, bytes.NewBuffer([]byte(reqBody)))
		require.NoError(t, err)
		h := NewHandler(httpClient, emptyTenantConfig, emptyProviderConfig, "")
		r := httptest.NewRecorder()

		//WHEN
		h.Deprovisioning(r, subscribeReq)
		resp := r.Result()

		//THEN
		expectedBody := "{\"error\":\"while executing unsubscription request: authorization header is required\"}\n"
		assertExpectedResponse(t, resp, expectedBody, http.StatusUnauthorized)
	})
}

func TestHandler_JobStatus(t *testing.T) {
	jobID := "d1a21d4a-be03-4da5-a0ce-a006fbc851a6"
	apiPath := fmt.Sprintf("/api/v1/jobs/%s", jobID)

	t.Run("Error when missing authorization header", func(t *testing.T) {
		//GIVEN
		getJobReq, err := http.NewRequest(http.MethodGet, url+apiPath, bytes.NewBuffer([]byte{}))
		require.NoError(t, err)
		h := NewHandler(nil, Config{}, ProviderConfig{}, jobID)
		r := httptest.NewRecorder()

		//WHEN
		h.JobStatus(r, getJobReq)
		resp := r.Result()

		//THEN
		expectedBody := "{\"error\":\"authorization header is required\"}\n"
		assertExpectedResponse(t, resp, expectedBody, http.StatusUnauthorized)
	})

	testCases := []struct {
		Name                 string
		RequestMethod        string
		RequestBody          string
		ExpectedResponseCode int
		ExpectedBody         string
		Token                string
	}{
		{
			Name:                 "Error when missing token",
			RequestMethod:        http.MethodGet,
			ExpectedBody:         "{\"error\":\"token value is required\"}\n",
			ExpectedResponseCode: http.StatusUnauthorized,
			Token:                "",
		},
		{
			Name:                 "Error when request method is not the expected one",
			RequestMethod:        http.MethodPost,
			ExpectedResponseCode: http.StatusMethodNotAllowed,
			Token:                token,
		},
		{
			Name:                 "Successful job status response",
			RequestMethod:        http.MethodGet,
			ExpectedResponseCode: http.StatusOK,
			ExpectedBody:         fmt.Sprintf("{\"id\":\"%s\",\"state\":\"SUCCEEDED\"}", jobID),
			Token:                token,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			getJobReq, err := http.NewRequest(testCase.RequestMethod, url+apiPath, bytes.NewBuffer([]byte(testCase.RequestBody)))
			require.NoError(t, err)
			getJobReq.Header.Add(oauth2.AuthorizationHeader, fmt.Sprintf("Bearer %s", testCase.Token))
			h := NewHandler(nil, Config{}, ProviderConfig{}, jobID)
			r := httptest.NewRecorder()

			//WHEN
			h.JobStatus(r, getJobReq)
			resp := r.Result()

			//THEN
			if len(testCase.ExpectedBody) > 0 {
				assertExpectedResponse(t, resp, testCase.ExpectedBody, testCase.ExpectedResponseCode)
			} else {
				require.Equal(t, testCase.ExpectedResponseCode, resp.StatusCode)
			}
		})
	}
}

func TestHandler_OnSubscription(t *testing.T) {
	apiPath := fmt.Sprintf("/tenants/v1/regional/%s/callback/%s", "region", "tenantID")

	testCases := []struct {
		Name                 string
		RequestMethod        string
		RequestBody          string
		ExpectedResponseCode int
		ExpectedBody         string
	}{
		{
			Name:                 "with PUT request",
			RequestMethod:        http.MethodPut,
			ExpectedBody:         "https://github.com/kyma-incubator/compass",
			ExpectedResponseCode: http.StatusOK,
		},
		{
			Name:                 "with DELETE request",
			RequestMethod:        http.MethodDelete,
			ExpectedBody:         "https://github.com/kyma-incubator/compass",
			ExpectedResponseCode: http.StatusOK,
		},
		{
			Name:                 "with invalid request method",
			RequestMethod:        http.MethodPost,
			ExpectedResponseCode: http.StatusMethodNotAllowed,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			onSubReq, err := http.NewRequest(testCase.RequestMethod, url+apiPath, bytes.NewBuffer([]byte(testCase.RequestBody)))
			require.NoError(t, err)
			h := NewHandler(nil, Config{}, ProviderConfig{}, "")
			r := httptest.NewRecorder()

			//WHEN
			h.OnSubscription(r, onSubReq)
			resp := r.Result()

			//THEN
			if len(testCase.ExpectedBody) > 0 {
				assertExpectedResponse(t, resp, testCase.ExpectedBody, testCase.ExpectedResponseCode)
			} else {
				require.Equal(t, testCase.ExpectedResponseCode, resp.StatusCode)
			}
		})
	}
}

func TestHandler_DependenciesConfigure(t *testing.T) {
	apiPath := "/v1/dependencies/configure"
	dependency := "dependency-name"

	testCases := []struct {
		Name                 string
		RequestMethod        string
		RequestBody          string
		ExpectedResponseCode int
		ExpectedBody         string
	}{
		{
			Name:                 "Error when request method is not the expected one",
			RequestMethod:        http.MethodPut,
			ExpectedResponseCode: http.StatusMethodNotAllowed,
		},
		{
			Name:                 "Error when the request body is empty",
			RequestMethod:        http.MethodPost,
			RequestBody:          "",
			ExpectedResponseCode: http.StatusInternalServerError,
			ExpectedBody:         "{\"error\":\"The request body is empty\"}\n",
		},
		{
			Name:                 "Successfully handled dependency configure request",
			RequestMethod:        http.MethodPost,
			RequestBody:          dependency,
			ExpectedResponseCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			//GIVEN
			depConfigureReq, err := http.NewRequest(testCase.RequestMethod, url+apiPath, bytes.NewBuffer([]byte(testCase.RequestBody)))
			require.NoError(t, err)
			h := NewHandler(nil, Config{}, ProviderConfig{}, "")
			r := httptest.NewRecorder()

			//WHEN
			h.DependenciesConfigure(r, depConfigureReq)
			resp := r.Result()

			//THEN
			if len(testCase.ExpectedBody) > 0 {
				assertExpectedResponse(t, resp, testCase.ExpectedBody, testCase.ExpectedResponseCode)
			} else {
				require.Equal(t, testCase.ExpectedResponseCode, resp.StatusCode)
			}
		})
	}
}

func TestHandler_Dependencies(t *testing.T) {
	depConfApiPath := "/v1/dependencies/configure"
	depApiPath := "/v1/dependencies"
	dependency := "dependency-name"

	t.Run("Successfully handled get dependency request", func(t *testing.T) {
		//GIVEN
		depConfigureReq, err := http.NewRequest(http.MethodPost, url+depConfApiPath, bytes.NewBuffer([]byte(dependency)))
		require.NoError(t, err)
		depReq, err := http.NewRequest(http.MethodGet, url+depApiPath, bytes.NewBuffer([]byte{}))
		require.NoError(t, err)
		h := NewHandler(nil, Config{}, ProviderConfig{}, "")

		//WHEN
		r := httptest.NewRecorder()
		h.DependenciesConfigure(r, depConfigureReq)
		resp := r.Result()

		//THEN
		require.Equal(t, http.StatusOK, resp.StatusCode)

		// WHEN
		r = httptest.NewRecorder()
		h.Dependencies(r, depReq)
		resp = r.Result()

		//THEN
		expectedBody := fmt.Sprintf("[{\"xsappname\":\"%s\"}]", dependency)
		assertExpectedResponse(t, resp, expectedBody, http.StatusOK)
	})
}

func assertExpectedResponse(t *testing.T, response *http.Response, expectedBody string, expectedStatusCode int) {
	require.Equal(t, expectedStatusCode, response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Equal(t, expectedBody, string(body))
}
