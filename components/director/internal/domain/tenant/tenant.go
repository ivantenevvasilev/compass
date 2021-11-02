package tenant

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
)

type key int

// TenantContextKey is the key under which the TenantCtx is saved in a given context.Context.
const TenantContextKey key = iota

// TenantCtx is the structure can be saved in a request context. It is used to determine the tenant context in which the request is being executed.
type TenantCtx struct {
	InternalID string
	ExternalID string
}

// LoadFromContext retrieves the internal tenant ID from the provided context. It returns error if such ID cannot be found.
func LoadFromContext(ctx context.Context) (string, error) {
	tenant, ok := ctx.Value(TenantContextKey).(TenantCtx)

	if !ok {
		return "", apperrors.NewCannotReadTenantError()
	}

	if tenant.InternalID == "" {
		return "", apperrors.NewTenantRequiredError()
	}

	return tenant.InternalID, nil
}

// LoadExternalFromContext retrieves the external tenant ID from the provided context. It returns error if such ID cannot be found.
func LoadExternalFromContext(ctx context.Context) (string, error) {
	tenant, ok := ctx.Value(TenantContextKey).(TenantCtx)

	if !ok {
		return "", apperrors.NewCannotReadTenantError()
	}

	if tenant.ExternalID == "" {
		return "", apperrors.NewTenantRequiredError()
	}

	return tenant.ExternalID, nil
}

// SaveToContext returns a child context of the provided context, including the provided tenant information.
// The internal tenant ID can be later retrieved from the context by calling LoadFromContext.
func SaveToContext(ctx context.Context, internalID, externalID string) context.Context {
	tenantCtx := TenantCtx{InternalID: internalID, ExternalID: externalID}
	return context.WithValue(ctx, TenantContextKey, tenantCtx)
}
