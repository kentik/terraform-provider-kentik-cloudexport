package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Note: values checked in tests below are provided by fake API Server from test-data.json
// (running in background).

func TestDataSourceCloudExportList(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { checkAPIServerConnection(t) },
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCloudExportDataSourceList,
				Check: resource.ComposeTestCheckFunc(
					// more properties are verified in TestDataSourceCloudExportItem* tests
					resource.TestCheckResourceAttr(exportsDS, "items.0.name", "test_terraform_aws_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.1.name", "test_terraform_gce_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.2.name", "test_terraform_ibm_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.3.name", "test_terraform_azure_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.4.name", "test_terraform_bgp_export"),
				),
			},
		},
	})
}

const (
	exportsDS                     = "data.kentik-cloudexport_list.exports"
	testCloudExportDataSourceList = `
		provider "kentik-cloudexport" {
			# apiurl = "http://localhost:8080" # KTAPI_URL env variable used instead
		}
		  
		data "kentik-cloudexport_list" "exports" {}
	`
)
