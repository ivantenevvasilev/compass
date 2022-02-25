package config

import "github.com/pkg/errors"

type HandlerConfig struct {
	TenantMappingEndpoint         string `envconfig:"default=/tenant-mapping"`
	AuthenticationMappingEndpoint string `envconfig:"default=/authn-mapping/{authenticator}"`
	ValidationIstioCertEndpoint   string `envconfig:"default=/v1/certificate/data/resolve"`
	RuntimeMappingEndpoint        string `envconfig:"default=/runtime-mapping"`
}

// Validate ensures the constructed Config contains valid property values
func (c *HandlerConfig) Validate() error {
	if c.AuthenticationMappingEndpoint == "" || c.TenantMappingEndpoint == "" || c.ValidationIstioCertEndpoint == "" {
		return errors.New("Missing configuration")
	}

	return nil
}
