package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceCloudExportList(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestCloudExportDataSourceList(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					// more properties are verified in TestDataSourceCloudExportItem* tests
					resource.TestCheckResourceAttr(exportsDS, "items.0.name", "test_terraform_aws_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.1.name", "test_terraform_gce_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.2.name", "test_terraform_ibm_export"),
					resource.TestCheckResourceAttr(exportsDS, "items.3.name", "test_terraform_azure_export"),
				),
			},
		},
	})
}

const (
	exportsDS = "data.kentik-cloudexport_list.exports"
)

func makeTestCloudExportDataSourceList(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		  
		data "kentik-cloudexport_list" "exports" {}
	`,
		apiURL,
	)
}
