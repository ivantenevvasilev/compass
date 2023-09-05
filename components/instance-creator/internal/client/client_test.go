package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/kyma-incubator/compass/components/instance-creator/internal/client/resources"

	"github.com/kyma-incubator/compass/components/instance-creator/internal/client/paths"

	"github.com/kyma-incubator/compass/components/instance-creator/internal/client"
	"github.com/kyma-incubator/compass/components/instance-creator/internal/client/automock"
	resmock "github.com/kyma-incubator/compass/components/instance-creator/internal/client/resources/automock"
	"github.com/kyma-incubator/compass/components/instance-creator/internal/client/types"
	"github.com/kyma-incubator/compass/components/instance-creator/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testConfig               = fixConfig()
	testConfigWithInvalidURL = fixConfigWithInvalidURL()
	invalidResponseBody      = []byte(`{"key: value}`)
	emptyResponseBody        = []byte(`{}`)
)

func TestClient_RetrieveResource(t *testing.T) {
	ctx := context.TODO()

	respBody, err := json.Marshal(fixServiceOfferings())
	require.NoError(t, err)

	respBodyWithNoMatchingCatalogName, err := json.Marshal(fixServiceOfferingsWithNoMatchingCatalogName())
	require.NoError(t, err)

	testCases := []struct {
		Name                  string
		CallerProvider        func(t *testing.T, cfg config.Config, region string) *automock.ExternalSvcCallerProvider
		ResourcesFn           func() *resmock.Resources
		ResourceMatchParamsFN func() *resmock.ResourceMatchParameters
		Config                config.Config
		WrongMatchParams      bool
		ExpectedResult        string
		ExpectedErrorMsg      string
	}{
		{
			Name:           "Success",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, respBody),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				resources.On("Match", types.ServiceOfferingMatchParameters{CatalogName: catalogName}).Return(offeringID, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:         testConfig,
			ExpectedResult: offeringID,
		},
		{
			Name:           "Error when building URL",
			CallerProvider: callerThatDoesNotGetCalled,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfigWithInvalidURL,
			ExpectedErrorMsg: "while building service offerings URL",
		},
		{
			Name:           "Error when getting caller fails",
			CallerProvider: callerProviderThatFails,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while getting caller for region:",
		},
		{
			Name:           "Error when caller fails",
			CallerProvider: callerThatDoesNotSucceed,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while executing request for listing service offerings for subaccount with ID:",
		},
		{
			Name:           "Error when response status code is not 200 OK",
			CallerProvider: callerThatReturnsBadStatus,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to get object(s), status:",
		},
		{
			Name:           "Error when unmarshalling response body",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, invalidResponseBody),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal service offerings:",
		},
		{
			Name:           "Error when Match fails due to type assertion error",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, respBodyWithNoMatchingCatalogName),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				resources.On("Match", types.ServicePlanMatchParameters{}).Return("", testError).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			WrongMatchParams: true,
			ExpectedErrorMsg: "while type asserting ResourceMatchParameters to ServiceOfferingMatchParameters",
		},
		{
			Name:           "Error when offering ID is empty",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, respBodyWithNoMatchingCatalogName),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceOfferingType)
				resources.On("Match", types.ServiceOfferingMatchParameters{CatalogName: catalogName}).Return("", testError).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceOfferingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "couldn't find service offering for catalog name:",
		},
	}

	for _, testCase := range testCases {
		// GIVEN
		callerProviderSvc := testCase.CallerProvider(t, testCase.Config, region)
		defer callerProviderSvc.AssertExpectations(t)

		testClient := client.NewClient(testCase.Config, callerProviderSvc)

		var matchParams resources.ResourceMatchParameters
		if testCase.WrongMatchParams {
			matchParams = types.ServicePlanMatchParameters{}
		} else {
			matchParams = types.ServiceOfferingMatchParameters{CatalogName: catalogName}
		}

		// WHEN
		resultResourceID, err := testClient.RetrieveResource(ctx, region, subaccountID, &types.ServiceOfferings{}, matchParams)

		// THEN
		if testCase.ExpectedErrorMsg != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.ExpectedErrorMsg)
		} else {
			assert.Equal(t, testCase.ExpectedResult, resultResourceID)
			assert.NoError(t, err)
		}
	}
}

func TestClient_RetrieveResourceByID(t *testing.T) {
	ctx := context.TODO()

	respBody, err := json.Marshal(fixServiceKey())
	require.NoError(t, err)

	testCases := []struct {
		Name                  string
		CallerProvider        func(t *testing.T, cfg config.Config, region string) *automock.ExternalSvcCallerProvider
		ResourceFn            func() *resmock.Resource
		ResourceMatchParamsFN func() *resmock.ResourceMatchParameters
		Config                config.Config
		ExpectedResult        *types.ServiceKey
		ExpectedErrorMsg      string
	}{
		{
			Name:           "Success",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, respBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceKeyID)
				resource.On("GetResourceType").Return(types.ServiceBindingType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:         testConfig,
			ExpectedResult: fixServiceKey(),
		},
		{
			Name:           "Error when building URL",
			CallerProvider: callerThatDoesNotGetCalled,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceKeyID)
				resource.On("GetResourceType").Return(types.ServiceBindingType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfigWithInvalidURL,
			ExpectedErrorMsg: "while building service binding URL",
		},
		{
			Name:           "Error when getting caller fails",
			CallerProvider: callerProviderThatFails,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceKeyID)
				resource.On("GetResourceType").Return(types.ServiceBindingType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while getting caller for region:",
		},
		{
			Name:           "Error when caller fails",
			CallerProvider: callerThatDoesNotSucceed,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceKeyID)
				resource.On("GetResourceType").Return(types.ServiceBindingType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while executing request for getting service binding for subaccount with ID:",
		},
		{
			Name:           "Error when response status code is not 200 OK",
			CallerProvider: callerThatReturnsBadStatus,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceKeyID)
				resource.On("GetResourceType").Return(types.ServiceBindingType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to get object(s), status:",
		},
		{
			Name:           "Error when unmarshalling response body",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, invalidResponseBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceKeyID)
				resource.On("GetResourceType").Return(types.ServiceBindingType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal service binding:",
		},
	}

	for _, testCase := range testCases {
		// GIVEN
		callerProviderSvc := testCase.CallerProvider(t, testCase.Config, region)
		defer callerProviderSvc.AssertExpectations(t)

		testClient := client.NewClient(testCase.Config, callerProviderSvc)

		// WHEN
		resultServiceKey, err := testClient.RetrieveResourceByID(ctx, region, subaccountID, &types.ServiceKey{}, types.ServiceKeyMatchParameters{})

		// THEN
		if testCase.ExpectedErrorMsg != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.ExpectedErrorMsg)
		} else {
			assert.Equal(t, testCase.ExpectedResult, resultServiceKey.(*types.ServiceKey))
			assert.NoError(t, err)
		}
	}
}

func TestClient_CreateResource(t *testing.T) {
	ctx := context.TODO()

	respBody, err := json.Marshal(fixServiceInstance(serviceInstanceID, serviceInstanceName))
	require.NoError(t, err)

	inProgressResponseBody, err := json.Marshal(fixOperationStatus(types.OperationStateInProgress, ""))
	require.NoError(t, err)

	opRespBody, err := json.Marshal(fixOperationStatus(types.OperationStateSucceeded, serviceInstanceID))
	require.NoError(t, err)

	opRespBodyWithFailedState, err := json.Marshal(fixOperationStatus(types.OperationStateFailed, serviceInstanceID))
	require.NoError(t, err)

	opRespBodyWithEmptyResourceID, err := json.Marshal(fixOperationStatusWithEmptyResourceID(types.OperationStateSucceeded))
	require.NoError(t, err)

	serviceInstanceWithNoID := fixServiceInstance(serviceInstanceID, serviceInstanceName)
	serviceInstanceWithNoID.ID = ""
	respBodyWithEmptyID, err := json.Marshal(serviceInstanceWithNoID)
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		CallerProvider   func(t *testing.T, cfg config.Config, region string) *automock.ExternalSvcCallerProvider
		ResourceFn       func() *resmock.Resource
		Config           config.Config
		InvalidParams    bool
		ExpectedResult   string
		ExpectedErrorMsg string
	}{
		// Sync flow
		{
			Name:           "Success in the synchronous case",
			CallerProvider: callerThatGetsCalledOnce(http.StatusCreated, respBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:         testConfig,
			ExpectedResult: serviceInstanceID,
		},
		{
			Name:           "Error when marshalling request body",
			CallerProvider: callerThatDoesNotGetCalled,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			Config:           testConfig,
			InvalidParams:    true,
			ExpectedErrorMsg: "failed to marshal service instance body",
		},
		{
			Name:           "Error when building URL",
			CallerProvider: callerThatDoesNotGetCalled,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfigWithInvalidURL,
			ExpectedErrorMsg: "while building service instance URL",
		},
		{
			Name:           "Error when getting caller fails",
			CallerProvider: callerProviderThatFails,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while getting caller for region:",
		},
		{
			Name:           "Error when caller fails",
			CallerProvider: callerThatDoesNotSucceed,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: testError.Error(),
		},
		{
			Name:           "Error when response status code is not 201 or 202",
			CallerProvider: callerThatReturnsBadStatus,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to create record of service instance, status:",
		},
		{
			Name:           "Error when unmarshalling response body",
			CallerProvider: callerThatGetsCalledOnce(http.StatusCreated, invalidResponseBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal service instance:",
		},
		{
			Name:           "Error when service instance ID is empty",
			CallerProvider: callerThatGetsCalledOnce(http.StatusCreated, respBodyWithEmptyID),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return("")
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "the service instance ID could not be empty",
		},
		// Async flow
		{
			Name:           "Success in the asynchronous case",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{respBody, opRespBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:         testConfig,
			ExpectedResult: serviceInstanceID,
		},
		{
			Name:           "Success in the asynchronous case when initially operation state is 'In Progress'",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK, http.StatusOK}, [][]byte{respBody, inProgressResponseBody, opRespBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:         testConfig,
			ExpectedResult: serviceInstanceID,
		},
		{
			Name:           "Error when location header key is missing",
			CallerProvider: callerThatGetsCalledOnce(http.StatusAccepted, respBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: fmt.Sprintf("the operation status path from %s header should not be empty", locationHeaderKey),
		},
		{
			Name:             "Error when caller fails the second time it is called",
			CallerProvider:   callerThatDoesNotSucceedTheLastTimeWhenCalled([]int{http.StatusAccepted}, [][]byte{respBody}, 1),
			Config:           testConfig,
			ExpectedErrorMsg: "while handling asynchronous creation of service instance in subaccount with ID:",
		},
		{
			Name:           "Error when operation response status code is not 200 OK",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusBadRequest}, [][]byte{respBody, opRespBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to get asynchronous object operation status. Received status:",
		},
		{
			Name:           "Error when unmarshalling operation status",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{respBody, invalidResponseBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal object operation status:",
		},
		{
			Name:           "Error when operation state is not 'Succeeded' or 'In Progress'",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{respBody, opRespBodyWithFailedState}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: fmt.Sprintf("the asynchronous object operation finished with state: %q", types.OperationStateFailed),
		},
		{
			Name:           "Error when returned service instance ID is empty",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{respBody, opRespBodyWithEmptyResourceID}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				resource.On("GetResourceURLPath", paths.ServiceInstancesPath)
				return resource
			},
			Config:           testConfig,
			ExpectedErrorMsg: "the service instance ID could not be empty",
		},
	}

	for _, testCase := range testCases {
		// GIVEN
		callerProviderSvc := testCase.CallerProvider(t, testCase.Config, region)
		defer callerProviderSvc.AssertExpectations(t)

		testClient := client.NewClient(testCase.Config, callerProviderSvc)

		inputParams := parameters
		if testCase.InvalidParams {
			inputParams = invalidParameters
		}

		// WHEN
		resultServiceInstanceID, err := testClient.CreateResource(ctx, region, subaccountID, types.ServiceInstanceReqBody{Name: serviceInstanceName, ServicePlanID: planID, Parameters: inputParams}, &types.ServiceInstance{})

		// THEN
		if testCase.ExpectedErrorMsg != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.ExpectedErrorMsg)
		} else {
			assert.Equal(t, testCase.ExpectedResult, resultServiceInstanceID)
			assert.NoError(t, err)
		}
	}
}

func TestClient_DeleteResource(t *testing.T) {
	ctx := context.TODO()

	inProgressResponseBody, err := json.Marshal(fixOperationStatus(types.OperationStateInProgress, ""))
	require.NoError(t, err)

	opRespBody, err := json.Marshal(fixOperationStatusWithEmptyResourceID(types.OperationStateSucceeded))
	require.NoError(t, err)

	opRespBodyWithFailedState, err := json.Marshal(fixOperationStatusWithEmptyResourceID(types.OperationStateFailed))
	require.NoError(t, err)

	testCases := []struct {
		Name                  string
		CallerProvider        func(t *testing.T, cfg config.Config, region string) *automock.ExternalSvcCallerProvider
		ResourceFn            func() *resmock.Resource
		ResourceMatchParamsFN func() *resmock.ResourceMatchParameters
		Config                config.Config
		ExpectedErrorMsg      string
	}{
		// Sync flow
		{
			Name:           "Success in the synchronous case",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, emptyResponseBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config: testConfig,
		},
		{
			Name:           "Error when building URL",
			CallerProvider: callerThatDoesNotGetCalled,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfigWithInvalidURL,
			ExpectedErrorMsg: "while building service instance URL",
		},
		{
			Name:           "Error when getting caller fails",
			CallerProvider: callerProviderThatFails,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while getting caller for region:",
		},
		{
			Name:           "Error when caller fails",
			CallerProvider: callerThatDoesNotSucceed,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: testError.Error(),
		},
		{
			Name:           "Error when response status code is not 200 or 202",
			CallerProvider: callerThatReturnsBadStatus,
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to delete service instance, status:",
		},
		// Async flow
		{
			Name:           "Success in the asynchronous case",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{emptyResponseBody, opRespBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config: testConfig,
		},
		{
			Name:           "Success in the asynchronous case when initially operation state is 'In Progress'",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK, http.StatusOK}, [][]byte{emptyResponseBody, inProgressResponseBody, opRespBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config: testConfig,
		},
		{
			Name:           "Error when location header key is missing",
			CallerProvider: callerThatGetsCalledOnce(http.StatusAccepted, emptyResponseBody),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: fmt.Sprintf("the operation status path from %s header should not be empty", locationHeaderKey),
		},
		{
			Name:           "Error when caller fails the second time it is called",
			CallerProvider: callerThatDoesNotSucceedTheLastTimeWhenCalled([]int{http.StatusAccepted}, [][]byte{emptyResponseBody}, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while deleting service instance with ID:",
		},
		{
			Name:           "Error when operation response status code is not 200 OK",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusBadRequest}, [][]byte{emptyResponseBody, opRespBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to get asynchronous object operation status. Received status:",
		},
		{
			Name:           "Error when unmarshalling operation status",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{emptyResponseBody, invalidResponseBody}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal object operation status:",
		},
		{
			Name:           "Error when operation state is not 'Succeeded' or 'In Progress'",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusAccepted, http.StatusOK}, [][]byte{emptyResponseBody, opRespBodyWithFailedState}, false, 1),
			ResourceFn: func() *resmock.Resource {
				resource := &resmock.Resource{}
				resource.On("GetResourceID").Return(serviceInstanceID)
				resource.On("GetResourceType").Return(types.ServiceInstanceType)
				return resource
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceInstancesPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: fmt.Sprintf("the asynchronous object operation finished with state: %q", types.OperationStateFailed),
		},
	}

	for _, testCase := range testCases {
		// GIVEN
		callerProviderSvc := testCase.CallerProvider(t, testCase.Config, region)
		defer callerProviderSvc.AssertExpectations(t)

		testClient := client.NewClient(testCase.Config, callerProviderSvc)

		// WHEN
		err := testClient.DeleteResource(ctx, region, subaccountID, types.ServiceInstance{ID: serviceInstanceID}, types.ServiceInstanceMatchParameters{})

		// THEN
		if testCase.ExpectedErrorMsg != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.ExpectedErrorMsg)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestClient_DeleteMultipleResources(t *testing.T) {
	ctx := context.TODO()

	respBody, err := json.Marshal(fixServiceKeys())
	require.NoError(t, err)

	opRespBody, err := json.Marshal(fixOperationStatusWithEmptyResourceID(types.OperationStateSucceeded))
	require.NoError(t, err)

	inProgressResponseBody, err := json.Marshal(fixOperationStatus(types.OperationStateInProgress, ""))
	require.NoError(t, err)

	opRespBodyWithFailedState, err := json.Marshal(fixOperationStatusWithEmptyResourceID(types.OperationStateFailed))
	require.NoError(t, err)

	testCases := []struct {
		Name                  string
		CallerProvider        func(t *testing.T, cfg config.Config, region string) *automock.ExternalSvcCallerProvider
		ResourcesFn           func() *resmock.Resources
		ResourceMatchParamsFN func() *resmock.ResourceMatchParameters
		Config                config.Config
		ExpectedResult        error
		ExpectedErrorMsg      string
	}{
		// Sync flow
		{
			Name:           "Success in the synchronous case",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusOK}, [][]byte{respBody, emptyResponseBody}, true, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config: testConfig,
		},
		{
			Name:           "Error when building URL",
			CallerProvider: callerThatDoesNotGetCalled,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfigWithInvalidURL,
			ExpectedErrorMsg: "while building service bindings URL",
		},
		{
			Name:             "Error when getting caller fails",
			CallerProvider:   callerProviderThatFails,
			Config:           testConfig,
			ExpectedErrorMsg: "while getting caller for region:",
		},
		{
			Name:           "Error when caller fails",
			CallerProvider: callerThatDoesNotSucceed,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: testError.Error(),
		},
		{
			Name:           "Error when response status code is not 200",
			CallerProvider: callerThatReturnsBadStatus,
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while executing request for listing service bindings for subaccount with ID:",
		},
		{
			Name:           "Error when unmarshalling service keys",
			CallerProvider: callerThatGetsCalledOnce(http.StatusOK, invalidResponseBody),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal service bindings:",
		},
		{
			Name:           "Error when caller fails on the second call",
			CallerProvider: callerThatDoesNotSucceedTheLastTimeWhenCalled([]int{http.StatusOK}, [][]byte{respBody}, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: testError.Error(),
		},
		{
			Name:           "Error when caller does not return 200 or 202 on the second call",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusBadRequest}, [][]byte{respBody, emptyResponseBody}, true, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to delete service binding, status:",
		},
		// Async flow
		{
			Name:           "Success in the asynchronous case",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusAccepted, http.StatusOK}, [][]byte{respBody, emptyResponseBody, opRespBody}, false, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config: testConfig,
		},
		{
			Name:           "Success in the asynchronous case when initially operation state is 'In Progress'",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusAccepted, http.StatusOK, http.StatusOK}, [][]byte{respBody, emptyResponseBody, inProgressResponseBody, opRespBody}, false, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config: testConfig,
		},
		{
			Name:           "Error when location header key is missing",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusAccepted}, [][]byte{respBody, emptyResponseBody}, true, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: fmt.Sprintf("the operation status path from %s header should not be empty", locationHeaderKey),
		},
		{
			Name:           "Error when caller fails the last time it is called",
			CallerProvider: callerThatDoesNotSucceedTheLastTimeWhenCalled([]int{http.StatusOK, http.StatusAccepted}, [][]byte{respBody, emptyResponseBody}, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "while deleting service binding with ID:",
		},
		{
			Name:           "Error when operation response status code is not 200 OK",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusAccepted, http.StatusBadRequest}, [][]byte{respBody, emptyResponseBody, emptyResponseBody}, false, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to get asynchronous object operation status. Received status:",
		},
		{
			Name:           "Error when unmarshalling operation status",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusAccepted, http.StatusOK}, [][]byte{respBody, emptyResponseBody, invalidResponseBody}, false, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: "failed to unmarshal object operation status:",
		},
		{
			Name:           "Error when operation state is not 'Succeeded' or 'In Progress'",
			CallerProvider: callerThatGetsCalledSeveralTimesInAsyncCase([]int{http.StatusOK, http.StatusAccepted, http.StatusOK}, [][]byte{respBody, emptyResponseBody, opRespBodyWithFailedState}, false, 2),
			ResourcesFn: func() *resmock.Resources {
				resources := &resmock.Resources{}
				resources.On("GetType").Return(types.ServiceBindingsType)
				resources.On("MatchMultiple", types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID}).Return([]string{serviceKeyID}, nil).Once()
				return resources
			},
			ResourceMatchParamsFN: func() *resmock.ResourceMatchParameters {
				resourceMatchParams := &resmock.ResourceMatchParameters{}
				resourceMatchParams.On("GetURLPath").Return(paths.ServiceBindingsPath)
				return resourceMatchParams
			},
			Config:           testConfig,
			ExpectedErrorMsg: fmt.Sprintf("the asynchronous object operation finished with state: %q", types.OperationStateFailed),
		},
	}

	for _, testCase := range testCases {
		// GIVEN
		callerProviderSvc := testCase.CallerProvider(t, testCase.Config, region)
		defer callerProviderSvc.AssertExpectations(t)

		testClient := client.NewClient(testCase.Config, callerProviderSvc)

		// WHEN
		err := testClient.DeleteMultipleResources(ctx, region, subaccountID, &types.ServiceKeys{}, types.ServiceKeyMatchParameters{ServiceInstanceID: serviceInstanceID})

		// THEN
		if testCase.ExpectedErrorMsg != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), testCase.ExpectedErrorMsg)
		} else {
			assert.NoError(t, err)
		}
	}
}
