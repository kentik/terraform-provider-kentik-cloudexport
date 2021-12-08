package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Note: values checked in tests below are provided by fake API Server from test-data.json
// (running in background).

//nolint: dupl
func TestDataSourceCloudExportItemAWS(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { checkAPIServerConnection(t) },
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCloudExportDataSourceItems,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceAWSDS, "id", "1"),
					resource.TestCheckResourceAttr(ceAWSDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAWSDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "name", "test_terraform_aws_export"),
					resource.TestCheckResourceAttr(ceAWSDS, "description", "terraform aws cloud export"),
					resource.TestCheckResourceAttr(ceAWSDS, "api_root", "http://localhost:8080/api"),
					resource.TestCheckResourceAttr(ceAWSDS, "flow_dest", "http://localhost:8080/flow"),
					resource.TestCheckResourceAttr(ceAWSDS, "plan_id", "11467"),
					resource.TestCheckResourceAttr(ceAWSDS, "cloud_provider", "aws"),
					resource.TestCheckResourceAttr(ceAWSDS, "current_status.0.status", "OK"),
					resource.TestCheckResourceAttr(ceAWSDS, "current_status.0.error_message", "No errors"),
					resource.TestCheckResourceAttr(ceAWSDS, "current_status.0.flow_found", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "current_status.0.api_access", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "current_status.0.storage_account_access", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.bucket", "terraform-aws-bucket"),
					resource.TestCheckResourceAttr(
						ceAWSDS, "aws.0.iam_role_arn", "arn:aws:iam::003740049406:role/trafficTerraformIngestRole",
					),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.region", "us-east-2"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.delete_after_read", "false"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.multiple_buckets", "false"),
				),
			},
		},
	})
}

func TestDataSourceCloudExportItemGCE(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { checkAPIServerConnection(t) },
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCloudExportDataSourceItems,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceGCPDS, "id", "2"),
					resource.TestCheckResourceAttr(ceGCPDS, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceGCPDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceGCPDS, "name", "test_terraform_gce_export"),
					resource.TestCheckResourceAttr(ceGCPDS, "description", "terraform gce cloud export"),
					resource.TestCheckResourceAttr(ceGCPDS, "api_root", "http://localhost:8080/api"),
					resource.TestCheckResourceAttr(ceGCPDS, "flow_dest", "http://localhost:8080/flow"),
					resource.TestCheckResourceAttr(ceGCPDS, "plan_id", "21600"),
					resource.TestCheckResourceAttr(ceGCPDS, "cloud_provider", "gce"),
					resource.TestCheckResourceAttr(ceGCPDS, "current_status.0.status", "NOK"),
					resource.TestCheckResourceAttr(ceGCPDS, "current_status.0.error_message", "Timeout"),
					resource.TestCheckResourceAttr(ceGCPDS, "current_status.0.flow_found", "false"),
					resource.TestCheckResourceAttr(ceGCPDS, "current_status.0.api_access", "false"),
					resource.TestCheckResourceAttr(ceGCPDS, "current_status.0.storage_account_access", "false"),
					resource.TestCheckResourceAttr(ceGCPDS, "gce.0.project", "project gce"),
					resource.TestCheckResourceAttr(ceGCPDS, "gce.0.subscription", "subscription gce"),
				),
			},
		},
	})
}

func TestDataSourceCloudExportItemIBM(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { checkAPIServerConnection(t) },
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCloudExportDataSourceItems,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceIBMDS, "id", "3"),
					resource.TestCheckResourceAttr(ceIBMDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceIBMDS, "enabled", "false"),
					resource.TestCheckResourceAttr(ceIBMDS, "name", "test_terraform_ibm_export"),
					resource.TestCheckResourceAttr(ceIBMDS, "description", "terraform ibm cloud export"),
					resource.TestCheckResourceAttr(ceIBMDS, "api_root", "http://localhost:8080/api"),
					resource.TestCheckResourceAttr(ceIBMDS, "flow_dest", "http://localhost:8080/flow"),
					resource.TestCheckResourceAttr(ceIBMDS, "plan_id", "11467"),
					resource.TestCheckResourceAttr(ceIBMDS, "cloud_provider", "ibm"),
					resource.TestCheckResourceAttr(ceIBMDS, "current_status.0.status", "OK"),
					resource.TestCheckResourceAttr(ceIBMDS, "current_status.0.error_message", "No errors"),
					resource.TestCheckResourceAttr(ceIBMDS, "current_status.0.flow_found", "false"),
					resource.TestCheckResourceAttr(ceIBMDS, "current_status.0.api_access", "false"),
					resource.TestCheckResourceAttr(ceIBMDS, "current_status.0.storage_account_access", "false"),
					resource.TestCheckResourceAttr(ceIBMDS, "ibm.0.bucket", "terraform-ibm-bucket"),
				),
			},
		},
	})
}

//nolint: dupl
func TestDataSourceCloudExportItemAzure(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { checkAPIServerConnection(t) },
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCloudExportDataSourceItems,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceAzureDS, "id", "4"),
					resource.TestCheckResourceAttr(ceAzureDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAzureDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAzureDS, "name", "test_terraform_azure_export"),
					resource.TestCheckResourceAttr(ceAzureDS, "description", "terraform azure cloud export"),
					resource.TestCheckResourceAttr(ceAzureDS, "api_root", "http://localhost:8080/api"),
					resource.TestCheckResourceAttr(ceAzureDS, "flow_dest", "http://localhost:8080/flow"),
					resource.TestCheckResourceAttr(ceAzureDS, "plan_id", "11467"),
					resource.TestCheckResourceAttr(ceAzureDS, "cloud_provider", "azure"),
					resource.TestCheckResourceAttr(ceAzureDS, "current_status.0.status", "OK"),
					resource.TestCheckResourceAttr(ceAzureDS, "current_status.0.error_message", "No errors"),
					resource.TestCheckResourceAttr(ceAzureDS, "current_status.0.flow_found", "false"),
					resource.TestCheckResourceAttr(ceAzureDS, "current_status.0.api_access", "false"),
					resource.TestCheckResourceAttr(ceAzureDS, "current_status.0.storage_account_access", "false"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.location", "centralus"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.resource_group", "traffic-generator"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.storage_account", "kentikstorage"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.subscription_id", "784bd5ec-122b-41b7-9719-22f23d5b49c8"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.security_principal_enabled", "true"),
				),
			},
		},
	})
}

func TestDataSourceCloudExportItemBGP(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { checkAPIServerConnection(t) },
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testCloudExportDataSourceItems,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceBGPDS, "id", "5"),
					resource.TestCheckResourceAttr(ceBGPDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceBGPDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceBGPDS, "name", "test_terraform_bgp_export"),
					resource.TestCheckResourceAttr(ceBGPDS, "description", "terraform bgp cloud export"),
					resource.TestCheckResourceAttr(ceBGPDS, "api_root", "http://localhost:8080/api"),
					resource.TestCheckResourceAttr(ceBGPDS, "flow_dest", "http://localhost:8080/flow"),
					resource.TestCheckResourceAttr(ceBGPDS, "plan_id", "11467"),
					resource.TestCheckResourceAttr(ceBGPDS, "cloud_provider", "bgp"),
					resource.TestCheckResourceAttr(ceBGPDS, "current_status.0.status", "OK"),
					resource.TestCheckResourceAttr(ceBGPDS, "current_status.0.error_message", "No errors"),
					resource.TestCheckResourceAttr(ceBGPDS, "current_status.0.flow_found", "false"),
					resource.TestCheckResourceAttr(ceBGPDS, "current_status.0.api_access", "false"),
					resource.TestCheckResourceAttr(ceBGPDS, "current_status.0.storage_account_access", "false"),
					resource.TestCheckResourceAttr(ceBGPDS, "bgp.0.apply_bgp", "true"),
					resource.TestCheckResourceAttr(ceBGPDS, "bgp.0.use_bgp_device_id", "1324"),
					resource.TestCheckResourceAttr(ceBGPDS, "bgp.0.device_bgp_type", "router"),
				),
			},
		},
	})
}

const (
	ceAWSDS   = "data.kentik-cloudexport_item.aws"
	ceAzureDS = "data.kentik-cloudexport_item.azure"
	ceBGPDS   = "data.kentik-cloudexport_item.bgp"
	ceGCPDS   = "data.kentik-cloudexport_item.gce"
	ceIBMDS   = "data.kentik-cloudexport_item.ibm"

	testCloudExportDataSourceItems = `
		provider "kentik-cloudexport" {
			# apiurl = "http://localhost:8080" # KTAPI_URL env variable used instead
		}
		  
		data "kentik-cloudexport_item" "aws" {
			id = "1"
		}
		
		data "kentik-cloudexport_item" "gce" {
			id = "2"
		}
		
		data "kentik-cloudexport_item" "ibm" {
			id = "3"
		}
		
		data "kentik-cloudexport_item" "azure" {
			id = "4"
		}
		
		data "kentik-cloudexport_item" "bgp" {
			id = "5"
		}
	`
)
