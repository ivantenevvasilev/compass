package claims

import (
	"context"
	"fmt"

	labelutils "github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/tenantmapping"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"
)

// RuntimeService is used to interact with runtimes
//go:generate mockery --name=RuntimeService --output=automock --outpkg=automock --case=underscore
type RuntimeService interface {
	GetLabel(context.Context, string, string) (*model.Label, error)
	ListByFilters(context.Context, []*labelfilter.LabelFilter) ([]*model.Runtime, error)
}

type validator struct {
	runtimesSvc                   RuntimeService
	transact                      persistence.Transactioner
	subscriptionProviderLabelKey  string
	consumerSubaccountIDsLabelKey string
}

// NewValidator creates new claims validator
func NewValidator(transact persistence.Transactioner, runtimesSvc RuntimeService, subscriptionProviderLabelKey, consumerSubaccountIDsLabelKey string) *validator {
	return &validator{
		transact:                      transact,
		runtimesSvc:                   runtimesSvc,
		subscriptionProviderLabelKey:  subscriptionProviderLabelKey,
		consumerSubaccountIDsLabelKey: consumerSubaccountIDsLabelKey,
	}
}

// Validate validates given id_token claims
func (v *validator) Validate(ctx context.Context, claims Claims) error {
	if err := claims.Valid(); err != nil {
		return errors.Wrapf(err, "while validating claims")
	}

	if claims.Tenant[tenantmapping.ConsumerTenantKey] == "" && claims.Tenant[tenantmapping.ExternalTenantKey] != "" {
		return apperrors.NewTenantNotFoundError(claims.Tenant[tenantmapping.ExternalTenantKey])
	}

	if claims.OnBehalfOf == "" {
		return nil
	}

	log.C(ctx).Infof("Consumer-Provider call by %s on behalf of %s. Proceeding with double authentication crosscheck...", claims.Tenant[tenantmapping.ProviderTenantKey], claims.Tenant[tenantmapping.ConsumerTenantKey])
	if len(claims.TokenClientID) == 0 {
		log.C(ctx).Errorf("Could not find consumer token client ID")
		return apperrors.NewUnauthorizedError("could not find consumer token client ID")
	}
	if len(claims.Region) == 0 {
		log.C(ctx).Errorf("Could not determine consumer token's region")
		return apperrors.NewUnauthorizedError("could not determine token's region")
	}

	filters := []*labelfilter.LabelFilter{
		labelfilter.NewForKeyWithQuery(v.subscriptionProviderLabelKey, fmt.Sprintf("\"%s\"", claims.TokenClientID)),
		labelfilter.NewForKeyWithQuery(tenant.RegionLabelKey, fmt.Sprintf("\"%s\"", claims.Region)),
	}

	ctxWithProviderTenant := tenant.SaveToContext(ctx, claims.Tenant[tenantmapping.ProviderTenantKey], claims.Tenant[tenantmapping.ProviderExternalTenantKey])

	log.C(ctx).Infof("Listing runtimes in provider tenant %s for labels %s: %s and %s: %s", claims.Tenant[tenantmapping.ProviderTenantKey], tenant.RegionLabelKey, claims.Region, v.subscriptionProviderLabelKey, claims.TokenClientID)

	tx, err := v.transact.Begin()
	if err != nil {
		return errors.Wrap(err, "while opening db transaction")
	}
	defer v.transact.RollbackUnlessCommitted(ctx, tx)
	ctxWithProviderTenant = persistence.SaveToContext(ctxWithProviderTenant, tx)

	runtimes, err := v.runtimesSvc.ListByFilters(ctxWithProviderTenant, filters)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Error while listing runtimes in provider tenant %s for labels %s: %s and %s: %s: %v", claims.Tenant[tenantmapping.ProviderTenantKey], tenant.RegionLabelKey, claims.Region, v.subscriptionProviderLabelKey, claims.TokenClientID, err)
		return errors.Wrapf(err, "failed to get runtimes in tenant %s for labels %s: %s and %s: %s", claims.Tenant[tenantmapping.ProviderTenantKey], tenant.RegionLabelKey, claims.Region, v.subscriptionProviderLabelKey, claims.TokenClientID)
	}

	log.C(ctx).Infof("Found %d runtimes in provider tenant %s for labels %s: %s and %s: %s", len(runtimes), claims.Tenant[tenantmapping.ProviderTenantKey], tenant.RegionLabelKey, claims.Region, v.subscriptionProviderLabelKey, claims.TokenClientID)

	expectedConsumerTenant := claims.Tenant[tenantmapping.ExternalTenantKey]
	found := false
	for _, runtime := range runtimes {
		label, err := v.runtimesSvc.GetLabel(ctxWithProviderTenant, runtime.ID, v.consumerSubaccountIDsLabelKey)
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				continue
			}
			return errors.Wrapf(err, "failed to get label %s for runtime with ID %s", v.consumerSubaccountIDsLabelKey, runtime.ID)
		}
		labelValues, err := labelutils.ValueToStringsSlice(label.Value)
		if err != nil {
			return err
		}

		for _, val := range labelValues {
			if val == expectedConsumerTenant {
				found = true
				break
			}
		}

		if found {
			break
		}
	}

	if !found {
		log.C(ctx).Errorf("Consumer's external tenant %s was not found in the %s label of any runtime in the provider tenant %s", expectedConsumerTenant, v.consumerSubaccountIDsLabelKey, claims.Tenant[tenantmapping.ProviderTenantKey])
		return apperrors.NewUnauthorizedError(fmt.Sprintf("Consumer's external tenant %s was not found in the %s label of any runtime in the provider tenant %s", expectedConsumerTenant, v.consumerSubaccountIDsLabelKey, claims.Tenant[tenantmapping.ProviderTenantKey]))
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "while committing db transaction")
	}

	return nil
}
