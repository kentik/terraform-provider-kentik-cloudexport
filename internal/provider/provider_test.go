package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/kentik/terraform-provider-kentik-cloudexport/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	t.Parallel()
	err := provider.New().InternalValidate()
	assert.NoError(t, err)
}

func TestProvider_Configure_MinimalConfig(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{{
			Config: makeMinimalConfig(server.URL()),
		}},
	})
}

func TestProvider_Configure_CustomRetryConfig(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{{
			Config: makeCustomRetryConfig(server.URL()),
		}},
	})
}

func TestProvider_Configure_InvalidRetryConfig(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{{
			Config:      makeInvalidRetryConfig(server.URL()),
			ExpectError: regexp.MustCompile("parse max_delay duration"),
		}},
	})
}

func makeMinimalConfig(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		// Trigger arbitrary action
		data "kentik-cloudexport_item" "aws" {
			id = "1"
		}
	`,
		apiURL,
	)
}

func makeCustomRetryConfig(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
			retry {
				max_attempts = 66
				min_delay = "100ms"
				max_delay = "1m"
			}
			log_payloads = true
		}
		
		// Trigger arbitrary action
		data "kentik-cloudexport_item" "aws" {
			id = "1"
		}
	`,
		apiURL,
	)
}

func makeInvalidRetryConfig(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
			retry {
				max_attempts = 66
				min_delay = "100ms"
				max_delay = "invalid-delay"
			}
			log_payloads = true
		}
		
		// Trigger arbitrary action
		data "kentik-cloudexport_item" "aws" {
			id = "1"
		}
	`,
		apiURL,
	)
}
