package provider_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/terraform-provider-kentik-cloudexport/internal/provider"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
func providerFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"kentik-cloudexport": func() (*schema.Provider, error) { //nolint: unparam
			return provider.New(), nil
		},
	}
}
