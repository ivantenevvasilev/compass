package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/certs"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/kyma-incubator/compass/tests/pkg/tenant"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSystemAccess(t *testing.T) {
	t.Run("TestDirectorCertificateAccess Integration System consumer: manage account tenant entities", func(t *testing.T) {
		ctx := context.Background()
		defaultTenantId := tenant.TestTenants.GetDefaultTenantID()

		// Build graphql director client configured with certificate
		clientKey, rawCertChain := certs.ClientCertPair(t, conf.ExternalCA.Certificate, conf.ExternalCA.Key)
		directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, clientKey, rawCertChain, conf.SkipSSLValidation)
		managedTenant := tenant.TestTenants.GetIDByName(t, tenant.TestIntegrationSystemManagedAccount)

		t.Log(fmt.Sprintf("Trying to create applications in account tenant %s", managedTenant))
		app, err := fixtures.RegisterApplication(t, ctx, directorCertSecuredClient, "managed-app", managedTenant)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, managedTenant, &app)
		require.NoError(t, err)
		require.NotEmpty(t, app.ID)

		t.Log(fmt.Sprintf("Trying to list applications in account tenant %s", managedTenant))
		apps := fixtures.GetApplicationPage(t, ctx, directorCertSecuredClient, managedTenant)
		require.Len(t, apps.Data, 1)
		require.Equal(t, app.Name, apps.Data[0].Name)

		t.Log(fmt.Sprintf("Trying to register runtime in account tenant %s", managedTenant))
		rt, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, directorCertSecuredClient, managedTenant, &graphql.RuntimeInput{Name: "managed-runtime"})
		defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, managedTenant, &rt)
		require.NoError(t, err)
		require.NotEmpty(t, rt.ID)

		t.Log(fmt.Sprintf("Trying to create application template in account tenant %s via client certificate", managedTenant))
		at, err := fixtures.CreateApplicationTemplate(t, ctx, directorCertSecuredClient, defaultTenantId, "test-app-template")
		defer fixtures.CleanupApplicationTemplate(t, ctx, dexGraphQLClient, managedTenant, &at)
		require.NoError(t, err)
		require.NotEmpty(t, at.ID)
	})

	t.Run("TestDirectorCertificateAccess Integration System consumer: cannot manage entities in non-managed tenant types", func(t *testing.T) {
		ctx := context.Background()

		// Build graphql director client configured with certificate
		clientKey, rawCertChain := certs.ClientCertPair(t, conf.ExternalCA.Certificate, conf.ExternalCA.Key)
		directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, clientKey, rawCertChain, conf.SkipSSLValidation)

		nonManagedTenant := tenant.TestTenants.GetIDByName(t, tenant.TestProviderSubaccount)

		t.Log(fmt.Sprintf("Trying to create applications in subaccount tenant %s", nonManagedTenant))
		app, err := fixtures.RegisterApplication(t, ctx, directorCertSecuredClient, "non-managed-app", nonManagedTenant)
		defer fixtures.CleanupApplication(t, ctx, dexGraphQLClient, nonManagedTenant, &app)
		require.Error(t, err)

		t.Log(fmt.Sprintf("Trying to list applications in subaccount tenant %s", nonManagedTenant))
		getAppReq := fixtures.FixGetApplicationsRequestWithPagination()
		apps := graphql.ApplicationPage{}
		err = testctx.Tc.RunOperationWithCustomTenant(ctx, directorCertSecuredClient, nonManagedTenant, getAppReq, &apps)
		require.Error(t, err)

		t.Log(fmt.Sprintf("Trying to register runtime in account tenant %s", nonManagedTenant))
		rt, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, directorCertSecuredClient, nonManagedTenant, &graphql.RuntimeInput{Name: "non-managed-runtime"})
		defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, nonManagedTenant, &rt)
		require.Error(t, err)

		t.Log(fmt.Sprintf("Trying to create application template in account tenant %s", nonManagedTenant))
		at, err := fixtures.CreateApplicationTemplate(t, ctx, directorCertSecuredClient, nonManagedTenant, "non-managed-app-template")
		defer fixtures.CleanupApplicationTemplate(t, ctx, dexGraphQLClient, nonManagedTenant, &at)
		require.Error(t, err)
	})
}
