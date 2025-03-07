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
	"context"
	"os"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/tests/pkg/subscription"
	"github.com/kyma-incubator/compass/tests/pkg/util"

	"github.com/kyma-incubator/compass/components/director/pkg/certloader"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

type config struct {
	DefaultTestTenant                                     string
	DefaultTestSubaccount                                 string
	DirectorExternalCertSecuredURL                        string
	DirectorGraphqlOauthURL                               string
	ORDAggregatorURL                                      string
	ORDServiceURL                                         string
	ORDAggregatorContainerName                            string `envconfig:"ORD_AGGREGATOR_CONTAINER_NAME"`
	ExternalServicesMockBaseURL                           string
	ExternalServicesMockUnsecuredURL                      string
	ExternalServicesMockUnsecuredWithAdditionalContentURL string
	ExternalServicesMockAbsoluteURL                       string
	ExternalServicesMockOrdCertSecuredURL                 string
	ExternalServicesMockUnsecuredMultiTenantURL           string
	ExternalServicesMockBasicURL                          string
	ExternalServicesMockOauthURL                          string
	ClientID                                              string
	ClientSecret                                          string
	BasicUsername                                         string
	BasicPassword                                         string
	ORDServiceDefaultResponseType                         string
	GlobalRegistryURL                                     string
	SubscriptionProviderAppNameValue                      string `envconfig:"APP_SUBSCRIPTION_PROVIDER_APP_NAME_VALUE"`
	TestConsumerSubaccountID                              string
	TestProviderSubaccountID                              string
	TokenPath                                             string
	ExternalClientCertSecretName                          string `envconfig:"APP_EXTERNAL_CLIENT_CERT_SECRET_NAME"`
	ExtSvcClientCertSecretName                            string `envconfig:"APP_EXT_SVC_CLIENT_CERT_SECRET_NAME"`
	SubscriptionProviderAppNameProperty                   string `envconfig:"APP_TENANT_PROVIDER_SUBSCRIPTION_PROVIDER_APP_NAME_PROPERTY"`
	CertLoaderConfig                                      certloader.Config
	ClientTimeout                                         time.Duration `envconfig:"default=60s"`
	SkipSSLValidation                                     bool          `envconfig:"default=false"`
	SubscriptionConfig                                    subscription.Config
	ORDWebhookMappings                                    string `envconfig:"APP_ORD_WEBHOOK_MAPPINGS"`
	ProxyApplicationTemplateName                          string `envconfig:"APP_PROXY_APPLICATION_TEMPLATE_NAME"`
}

var (
	testConfig config

	certSecuredGraphQLClient *graphql.Client
	certCache                certloader.Cache
)

func TestMain(m *testing.M) {
	err := envconfig.Init(&testConfig)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "while initializing envconfig"))
	}

	ctx := context.Background()

	certCache, err = certloader.StartCertLoader(ctx, testConfig.CertLoaderConfig)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "while starting cert cache"))
	}

	if err := util.WaitForCache(certCache); err != nil {
		log.D().Fatal(err)
	}

	certSecuredGraphQLClient = gql.NewCertAuthorizedGraphQLClientWithCustomURL(testConfig.DirectorExternalCertSecuredURL, certCache.Get()[testConfig.ExternalClientCertSecretName].PrivateKey, certCache.Get()[testConfig.ExternalClientCertSecretName].Certificate, testConfig.SkipSSLValidation)

	exitVal := m.Run()
	os.Exit(exitVal)

}
