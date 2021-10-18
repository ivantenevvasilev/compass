package tests

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/tests/pkg/assertions"

	directorSchema "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/request"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	tenantHeader            = "Tenant"
	applicationTypeLabelKey = "applicationType"

	descriptionField      = "description"
	shortDescriptionField = "shortDescription"
	apisField             = "apis"
	eventsField           = "events"

	expectedSpecType                        = "openapi-v3"
	expectedSpecFormat                      = "application/json"
	expectedSystemInstanceName              = "test-app"
	expectedSecondSystemInstanceName        = "second-test-app"
	expectedThirdSystemInstanceName         = "third-test-app"
	expectedFourthSystemInstanceName        = "fourth-test-app"
	expectedFifthSystemInstanceName         = "fifth-test-app"
	expectedSystemInstanceDescription       = "test-app1-description"
	expectedSecondSystemInstanceDescription = "test-app2-description"
	expectedThirdSystemInstanceDescription  = "test-app3-description"
	expectedFourthSystemInstanceDescription = "test-app4-description"
	expectedFifthSystemInstanceDescription  = "test-app5-description"
	expectedBundleTitle                     = "BUNDLE TITLE"
	secondExpectedBundleTitle               = "BUNDLE TITLE 2"
	expectedBundleDescription               = "lorem ipsum dolor nsq sme"
	secondExpectedBundleDescription         = "foo bar"
	expectedPackageTitle                    = "PACKAGE 1 TITLE"
	expectedPackageDescription              = "lorem ipsum dolor set"
	expectedProductTitle                    = "PRODUCT TITLE"
	expectedProductShortDescription         = "lorem ipsum"
	firstAPIExpectedTitle                   = "API TITLE"
	firstAPIExpectedDescription             = "lorem ipsum dolor sit amet"
	firstEventTitle                         = "EVENT TITLE"
	firstEventDescription                   = "lorem ipsum dolor sit amet"
	secondEventTitle                        = "EVENT TITLE 2"
	secondEventDescription                  = "lorem ipsum dolor sit amet"
	expectedTombstoneOrdID                  = "ns:apiResource:API_ID2:v1"
	expectedVendorTitle                     = "SAP"

	expectedNumberOfSystemInstances           = 5
	expectedNumberOfPackages                  = 5
	expectedNumberOfBundles                   = 10
	expectedNumberOfProducts                  = 5
	expectedNumberOfAPIs                      = 5
	expectedNumberOfResourceDefinitionsPerAPI = 3
	expectedNumberOfEvents                    = 10
	expectedNumberOfTombstones                = 5
	expectedNumberOfVendors                   = 10

	expectedNumberOfAPIsInFirstBundle    = 1
	expectedNumberOfAPIsInSecondBundle   = 1
	expectedNumberOfEventsInFirstBundle  = 2
	expectedNumberOfEventsInSecondBundle = 2

	testTimeoutAdditionalBuffer = 5 * time.Minute
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randString(n int) string {
	buffer := make([]rune, n)
	for i := range buffer {
		buffer[i] = letters[rand.Intn(len(letters))]
	}
	return string(buffer)
}

func TestORDAggregator(t *testing.T) {
	basicORDConfigSecurity := &fixtures.ORDConfigSecurity{
		Enabled:  true,
		Username: randString(8),
		Password: randString(8),
	}

	oauthORDConfigSecurity := &fixtures.ORDConfigSecurity{
		Enabled:  true,
		Username: randString(8),
		Password: randString(8),
		TokenURL: testConfig.ExternalServicesMockBaseURL + "/oauth/token",
	}

	toggleORDConfigSecurity(t, testConfig.ExternalServicesMockBaseURL+"/.well-known/open-resource-discovery/basic/configure", basicORDConfigSecurity)
	toggleORDConfigSecurity(t, testConfig.ExternalServicesMockBaseURL+"/.well-known/open-resource-discovery/oauth/configure", oauthORDConfigSecurity)

	var appInput, secondAppInput, thirdAppInput, fourthAppInput, fifthAppInput directorSchema.ApplicationRegisterInput
	t.Run("Verifying ORD Document to be valid", func(t *testing.T) {
		appInput = fixtures.FixSampleApplicationRegisterInputWithORDWebhooks(expectedSystemInstanceName, expectedSystemInstanceDescription, testConfig.ExternalServicesMockAbsoluteURL, &fixtures.ORDConfigSecurity{})
		secondAppInput = fixtures.FixSampleApplicationRegisterInputWithORDWebhooks(expectedSecondSystemInstanceName, expectedSecondSystemInstanceDescription, testConfig.ExternalServicesMockBaseURL, &fixtures.ORDConfigSecurity{})
		thirdAppInput = fixtures.FixSampleApplicationRegisterInputWithORDWebhooks(expectedThirdSystemInstanceName, expectedThirdSystemInstanceDescription, testConfig.ExternalServicesMockBaseURL+"/basic/.well-known/open-resource-discovery", basicORDConfigSecurity)
		fourthAppInput = fixtures.FixSampleApplicationRegisterInputWithORDWebhooks(expectedFourthSystemInstanceName, expectedFourthSystemInstanceDescription, testConfig.ExternalServicesMockBaseURL+"/oauth/.well-known/open-resource-discovery", oauthORDConfigSecurity)
		fifthAppInput = fixtures.FixSampleApplicationRegisterInputWithORDWebhooks(expectedFifthSystemInstanceName, expectedFifthSystemInstanceDescription, testConfig.ExternalServicesMockBaseURL+"/cert/.well-known/open-resource-discovery", &fixtures.ORDConfigSecurity{})

		systemInstancesMap := make(map[string]string)
		systemInstancesMap[expectedSystemInstanceName] = expectedSystemInstanceDescription
		systemInstancesMap[expectedSecondSystemInstanceName] = expectedSecondSystemInstanceDescription
		systemInstancesMap[expectedThirdSystemInstanceName] = expectedThirdSystemInstanceDescription
		systemInstancesMap[expectedFourthSystemInstanceName] = expectedFourthSystemInstanceDescription
		systemInstancesMap[expectedFifthSystemInstanceName] = expectedFifthSystemInstanceDescription

		eventsMap := make(map[string]string)
		eventsMap[firstEventTitle] = firstEventDescription
		eventsMap[secondEventTitle] = secondEventDescription

		bundlesMap := make(map[string]string)
		bundlesMap[expectedBundleTitle] = expectedBundleDescription
		bundlesMap[secondExpectedBundleTitle] = secondExpectedBundleDescription

		bundlesAPIsNumberMap := make(map[string]int)
		bundlesAPIsNumberMap[expectedBundleTitle] = expectedNumberOfAPIsInFirstBundle
		bundlesAPIsNumberMap[secondExpectedBundleTitle] = expectedNumberOfAPIsInSecondBundle

		bundlesAPIsData := make(map[string][]string)
		bundlesAPIsData[expectedBundleTitle] = []string{firstAPIExpectedTitle}
		bundlesAPIsData[secondExpectedBundleTitle] = []string{firstAPIExpectedTitle}

		bundlesEventsNumberMap := make(map[string]int)
		bundlesEventsNumberMap[expectedBundleTitle] = expectedNumberOfEventsInFirstBundle
		bundlesEventsNumberMap[secondExpectedBundleTitle] = expectedNumberOfEventsInSecondBundle

		bundlesEventsData := make(map[string][]string)
		bundlesEventsData[expectedBundleTitle] = []string{firstEventTitle, secondEventTitle}
		bundlesEventsData[secondExpectedBundleTitle] = []string{firstEventTitle, secondEventTitle}

		ctx := context.Background()

		app, err := fixtures.RegisterApplicationFromInput(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, appInput)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &app)
		require.NoError(t, err)

		secondApp, err := fixtures.RegisterApplicationFromInput(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, secondAppInput)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &secondApp)
		require.NoError(t, err)

		thirdApp, err := fixtures.RegisterApplicationFromInput(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, thirdAppInput)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &thirdApp)
		require.NoError(t, err)

		fourthApp, err := fixtures.RegisterApplicationFromInput(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, fourthAppInput)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &fourthApp)
		require.NoError(t, err)

		fifthApp, err := fixtures.RegisterApplicationFromInput(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, fifthAppInput)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, testConfig.DefaultTestTenant, &fifthApp)
		require.NoError(t, err)

		t.Log("Create integration system")
		intSys, err := fixtures.RegisterIntegrationSystem(t, ctx, dexGraphQLClient, "", "test-int-system")
		defer fixtures.CleanupIntegrationSystem(t, ctx, dexGraphQLClient, "", intSys)
		require.NoError(t, err)
		require.NotEmpty(t, intSys.ID)

		intSystemCredentials := fixtures.RequestClientCredentialsForIntegrationSystem(t, ctx, dexGraphQLClient, "", intSys.ID)
		defer fixtures.DeleteSystemAuthForIntegrationSystem(t, ctx, dexGraphQLClient, intSystemCredentials.ID)

		unsecuredHttpClient := http.DefaultClient
		unsecuredHttpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		oauthCredentialData, ok := intSystemCredentials.Auth.Credential.(*directorSchema.OAuthCredentialData)
		require.True(t, ok)

		conf := &clientcredentials.Config{
			ClientID:     oauthCredentialData.ClientID,
			ClientSecret: oauthCredentialData.ClientSecret,
			TokenURL:     oauthCredentialData.URL,
		}

		ctx = context.WithValue(ctx, oauth2.HTTPClient, unsecuredHttpClient)
		httpClient := conf.Client(ctx)
		httpClient.Timeout = 10 * time.Second

		scheduleTime, err := parseCronTime(testConfig.AggregatorSchedule)
		require.NoError(t, err)

		defaultTestTimeout := 2*scheduleTime + testTimeoutAdditionalBuffer
		defaultCheckInterval := defaultTestTimeout / 20

		err = verifyORDDocument(defaultCheckInterval, defaultTestTimeout, func() bool {
			var respBody string

			// Verify system instances
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/systemInstances?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})
			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfSystemInstances {
				t.Log("Missing System Instances...will try again")
				return false
			}
			assertions.AssertMultipleEntitiesFromORDService(t, respBody, systemInstancesMap, expectedNumberOfSystemInstances)

			// Verify packages
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/packages?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfPackages {
				t.Log("Missing Packages...will try again")
				return false
			}
			assertions.AssertSingleEntityFromORDService(t, respBody, expectedNumberOfPackages, expectedPackageTitle, expectedPackageDescription, descriptionField)
			t.Log("Successfully verified packages")

			// Verify bundles
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/consumptionBundles?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfBundles {
				t.Log("Missing Bundles...will try again")
				return false
			}
			assertions.AssertMultipleEntitiesFromORDService(t, respBody, bundlesMap, expectedNumberOfBundles)
			t.Log("Successfully verified bundles")

			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/consumptionBundles?$expand=apis&$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})
			assertions.AssertRelationBetweenBundleAndEntityFromORDService(t, respBody, apisField, bundlesAPIsNumberMap, bundlesAPIsData)
			t.Log("Successfully verified relation between apis and bundles")

			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/consumptionBundles?$expand=events&$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})
			assertions.AssertRelationBetweenBundleAndEntityFromORDService(t, respBody, eventsField, bundlesEventsNumberMap, bundlesEventsData)
			t.Log("Successfully verified relation between events and bundles")

			// Verify products
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/products?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfProducts {
				t.Log("Missing Products...will try again")
				return false
			}
			assertions.AssertSingleEntityFromORDService(t, respBody, expectedNumberOfProducts, expectedProductTitle, expectedProductShortDescription, shortDescriptionField)
			t.Log("Successfully verified products")

			// Verify apis
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/apis?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfAPIs {
				t.Log("Missing APIs...will try again")
				return false
			}
			// In the document there are actually 2 APIs but there is a tombstone for the second one so in the end there will be only one API
			assertions.AssertSingleEntityFromORDService(t, respBody, expectedNumberOfAPIs, firstAPIExpectedTitle, firstAPIExpectedDescription, descriptionField)
			t.Log("Successfully verified apis")

			// Verify the api spec
			specs := gjson.Get(respBody, fmt.Sprintf("value.%d.resourceDefinitions", 0)).Array()
			require.Equal(t, expectedNumberOfResourceDefinitionsPerAPI, len(specs))

			var specURL string
			for _, s := range specs {
				specType := s.Get("type").String()
				specFormat := s.Get("mediaType").String()
				if specType == expectedSpecType && specFormat == expectedSpecFormat {
					specURL = s.Get("url").String()
					break
				}
			}

			respBody = makeRequestWithHeaders(t, httpClient, specURL, map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})
			if len(respBody) == 0 || !strings.Contains(respBody, "swagger") {
				t.Logf("Spec %s not successfully fetched... will try again", specURL)
				return false
			}
			t.Log("Successfully verified api spec")

			// Verify events
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/events?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfEvents {
				t.Log("Missing Events...will try again")
				return false
			}
			assertions.AssertMultipleEntitiesFromORDService(t, respBody, eventsMap, expectedNumberOfEvents)
			t.Log("Successfully verified events")

			// Verify tombstones
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/tombstones?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfTombstones {
				t.Log("Missing Tombstones...will try again")
				return false
			}
			assertions.AssertTombstoneFromORDService(t, respBody, expectedNumberOfTombstones, expectedTombstoneOrdID)
			t.Log("Successfully verified tombstones")

			// Verify vendors
			respBody = makeRequestWithHeaders(t, httpClient, testConfig.ORDServiceURL+"/vendors?$format=json", map[string][]string{tenantHeader: {testConfig.DefaultTestTenant}})

			if len(gjson.Get(respBody, "value").Array()) < expectedNumberOfVendors {
				t.Log("Missing Vendors...will try again")
				return false
			}
			assertions.AssertVendorFromORDService(t, respBody, expectedNumberOfVendors, expectedVendorTitle)
			t.Log("Successfully verified vendors")

			return true
		})
		require.NoError(t, err)
	})
}

func verifyORDDocument(interval time.Duration, timeout time.Duration, conditionalFunc func() bool) error {
	done := time.After(timeout)
	ticker := time.Tick(interval)

	for {
		if conditionalFunc() {
			return nil
		}

		select {
		case <-done:
			return errors.New("timeout waiting for entities to be present in DB")
		case <-ticker:
		}
	}
}

func makeRequestWithHeaders(t *testing.T, httpClient *http.Client, url string, headers map[string][]string) string {
	return request.MakeRequestWithHeadersAndStatusExpect(t, httpClient, url, headers, http.StatusOK, testConfig.ORDServiceDefaultResponseType)
}

func toggleORDConfigSecurity(t *testing.T, url string, ordSecurityConfig *fixtures.ORDConfigSecurity) {
	body, err := json.Marshal(ordSecurityConfig)
	require.NoError(t, err)

	reader := bytes.NewReader(body)
	response, err := http.DefaultClient.Post(url, "application/json", reader)
	require.NoError(t, err)
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Logf("Could not close response body %s", err)
		}
	}()
	if response.StatusCode != http.StatusOK {
		bytes, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)
		t.Fatalf("Failed to toggle ORD Config security to %t: %s", ordSecurityConfig.Enabled, string(bytes))
	}
}
