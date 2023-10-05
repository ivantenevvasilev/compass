package subscription_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/open_resource_discovery/apiclient"

	"github.com/kyma-incubator/compass/components/director/internal/domain/subscription"
	"github.com/kyma-incubator/compass/components/director/internal/domain/subscription/automock"
	persistenceautomock "github.com/kyma-incubator/compass/components/director/pkg/persistence/automock"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence/txtest"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testErr                            = errors.New("new-error")
	appID                              = "testAppID"
	appTemplateID                      = "testAppTemplateID"
	txGen                              = txtest.NewTransactionContextGenerator(testErr)
	payload                            = fmt.Sprintf("{\"subscriptionGUID\":\"%s\",\"dependentServiceInstancesInfo\":[{\"appId\":\"%s\",\"appName\":\"%s\",\"providerSubaccountID\":\"%s\"}]}", subscriptionID, subscriptionProviderID, subscriptionAppName, providerSubaccountID)
	payloadWithoutAppID                = fmt.Sprintf("{\"subscriptionGUID\":\"%s\",\"dependentServiceInstancesInfo\":[{\"appName\":\"%s\",\"providerSubaccountID\":\"%s\"}]}", subscriptionID, subscriptionAppName, providerSubaccountID)
	payloadWithoutProviderSubaccountID = fmt.Sprintf("{\"subscriptionGUID\":\"%s\",\"dependentServiceInstancesInfo\":[{\"appId\":\"%s\",\"appName\":\"%s\"}]}", subscriptionID, subscriptionProviderID, subscriptionAppName)
	payloadWithoutSubscriptionAppName  = fmt.Sprintf("{\"subscriptionGUID\":\"%s\",\"dependentServiceInstancesInfo\":[{\"appId\":\"%s\",\"providerSubaccountID\":\"%s\"}]}", subscriptionID, subscriptionProviderID, providerSubaccountID)
	payloadWithoutSubscription         = fmt.Sprintf("{\"dependentServiceInstancesInfo\":[{\"appId\":\"%s\",\"appName\":\"%s\",\"providerSubaccountID\":\"%s\"}]}", subscriptionProviderID, subscriptionAppName, providerSubaccountID)
)

func TestResolver_SubscribeTenant(t *testing.T) {
	testCases := []struct {
		Name            string
		TransactionerFn func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn       func() *automock.SubscriptionService
		Payload         string
		ExpectedOutput  bool
		ExpectedErr     error
	}{
		{
			Name:            "Success for Application flow",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.ApplicationTemplate, nil).Once()
				svc.On("SubscribeTenantToApplication", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionAppName, subscriptionID, payload).Return(true, appID, appTemplateID, nil).Once()
				svc.AssertNotCalled(t, "SubscribeTenantToRuntime")
				return svc
			},
			Payload:        payload,
			ExpectedOutput: true,
			ExpectedErr:    nil,
		},
		{
			Name:            "Success for Runtime flow",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.Runtime, nil).Once()
				svc.AssertNotCalled(t, "SubscribeTenantToApplication")
				svc.On("SubscribeTenantToRuntime", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionAppName, subscriptionID).Return(true, nil).Once()
				return svc
			},
			Payload:        payload,
			ExpectedOutput: true,
			ExpectedErr:    nil,
		},
		{
			Name:            "Error on transaction begin",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.AssertNotCalled(t, "SubscribeTenantToRuntime")
				svc.AssertNotCalled(t, "SubscribeTenantToApplication")
				svc.AssertNotCalled(t, "SubscribeTenantToRuntime")
				return svc
			},
			Payload:        payload,
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error on flow determination",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.ApplicationTemplate, testErr).Once()
				svc.AssertNotCalled(t, "SubscribeTenantToApplication")
				svc.AssertNotCalled(t, "SubscribeTenantToRuntime")
				return svc
			},
			Payload:        payload,
			ExpectedOutput: false,
			ExpectedErr:    errors.Wrapf(testErr, "while determining subscription flow"),
		},
		{
			Name:            "Error on subscription to applications fails",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.ApplicationTemplate, nil).Once()
				svc.On("SubscribeTenantToApplication", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionAppName, subscriptionID, payload).Return(false, "", "", testErr).Once()
				svc.AssertNotCalled(t, "SubscribeTenantToRuntime")
				return svc
			},
			Payload:        payload,
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error on subscription to runtime fails",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.Runtime, nil).Once()
				svc.On("SubscribeTenantToRuntime", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionAppName, subscriptionID).Return(false, testErr).Once()
				svc.AssertNotCalled(t, "SubscribeTenantToApplication")
				return svc
			},
			Payload:        payload,
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error on commit fail",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.Runtime, nil).Once()
				svc.On("SubscribeTenantToRuntime", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionAppName, subscriptionID).Return(true, nil).Once()
				svc.AssertNotCalled(t, "SubscribeTenantToApplication") //todo test rtm
				return svc
			},
			Payload:        payload,
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error when provider ID is missing",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				return svc
			},
			Payload:        payloadWithoutAppID,
			ExpectedOutput: false,
			ExpectedErr:    errors.New("Provider ID should not be empty"),
		},
		{
			Name:            "Error when subscription app name is missing",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				return svc
			},
			Payload:        payloadWithoutSubscriptionAppName,
			ExpectedOutput: false,
			ExpectedErr:    errors.New("Subscription app name should not be empty"),
		},
		{
			Name:            "Error when provider subaccount ID is missing",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				return svc
			},
			Payload:        payloadWithoutProviderSubaccountID,
			ExpectedOutput: false,
			ExpectedErr:    errors.New("Provider subaccount ID should not be empty"),
		},
		{
			Name:            "Error when subscription ID is missing",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				return svc
			},
			Payload:        payloadWithoutSubscription,
			ExpectedOutput: false,
			ExpectedErr:    errors.New("Subscription ID should not be empty"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			persistTx, transact := testCase.TransactionerFn()
			svc := testCase.ServiceFn()
			resolver := subscription.NewResolver(transact, svc, apiclient.OrdAggregatorClientConfig{})

			defer mock.AssertExpectationsForObjects(t, transact, persistTx, svc)

			// WHEN
			result, err := resolver.SubscribeTenant(context.TODO(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionAppName, testCase.Payload)

			// THEN
			assert.Equal(t, testCase.ExpectedOutput, result)

			if testCase.ExpectedErr != nil {
				assert.Equal(t, testCase.ExpectedErr.Error(), err.Error())
			}
		})
	}
}

func TestResolver_UnsubscribeTenant(t *testing.T) {
	testCases := []struct {
		Name            string
		TransactionerFn func() (*persistenceautomock.PersistenceTx, *persistenceautomock.Transactioner)
		ServiceFn       func() *automock.SubscriptionService
		ExpectedOutput  bool
		ExpectedErr     error
	}{
		{
			Name:            "Success when flow is runtime",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("UnsubscribeTenantFromRuntime", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionID).Return(true, nil).Once()
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.Runtime, nil).Once()
				svc.AssertNotCalled(t, "UnsubscribeTenantFromApplication")
				return svc
			},
			ExpectedOutput: true,
			ExpectedErr:    nil,
		},
		{
			Name:            "Success when flow is application",
			TransactionerFn: txGen.ThatSucceeds,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("UnsubscribeTenantFromApplication", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionID).Return(true, nil).Once()
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.ApplicationTemplate, nil).Once()
				svc.AssertNotCalled(t, "UnsubscribeTenantFromRuntime")
				return svc
			},
			ExpectedOutput: true,
			ExpectedErr:    nil,
		},
		{
			Name:            "Error on transaction begin",
			TransactionerFn: txGen.ThatFailsOnBegin,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.AssertNotCalled(t, "UnsubscribeTenantFromRuntime")
				return svc
			},
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error determining flow type",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.AssertNotCalled(t, "UnsubscribeTenantFromRuntime")
				svc.AssertNotCalled(t, "UnsubscribeTenantFromApplication")
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.ApplicationTemplate, testErr).Once()
				return svc
			},
			ExpectedOutput: false,
			ExpectedErr:    errors.Wrapf(testErr, "while determining subscription flow"),
		},
		{
			Name:            "Error on unsubscription from runtime fail",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("UnsubscribeTenantFromRuntime", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionID).Return(false, testErr).Once()
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.Runtime, nil).Once()

				return svc
			},
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error on unsubscription from application fail",
			TransactionerFn: txGen.ThatDoesntExpectCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("UnsubscribeTenantFromApplication", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionID).Return(false, testErr).Once()
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.ApplicationTemplate, nil).Once()

				return svc
			},
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
		{
			Name:            "Error on commit fail",
			TransactionerFn: txGen.ThatFailsOnCommit,
			ServiceFn: func() *automock.SubscriptionService {
				svc := &automock.SubscriptionService{}
				svc.On("DetermineSubscriptionFlow", txtest.CtxWithDBMatcher(), subscriptionProviderID, tenantRegion).Return(resource.Runtime, nil).Once()
				svc.On("UnsubscribeTenantFromRuntime", txtest.CtxWithDBMatcher(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, subscriptionID).Return(true, nil).Once()

				return svc
			},
			ExpectedOutput: false,
			ExpectedErr:    testErr,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			persistTx, transact := testCase.TransactionerFn()
			svc := testCase.ServiceFn()
			resolver := subscription.NewResolver(transact, svc, apiclient.OrdAggregatorClientConfig{})

			defer mock.AssertExpectationsForObjects(t, transact, persistTx, svc)

			// WHEN
			result, err := resolver.UnsubscribeTenant(context.TODO(), subscriptionProviderID, subaccountTenantExtID, providerSubaccountID, consumerTenantID, tenantRegion, payload)

			// THEN
			assert.Equal(t, testCase.ExpectedOutput, result)

			if testCase.ExpectedErr != nil {
				assert.Equal(t, testCase.ExpectedErr.Error(), err.Error())
			}
		})
	}
}
