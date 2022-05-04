package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Note: we only check the user-provided values as we don't control the server-provided ones

func TestResourceCloudExportAWS(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccResourceCloudExportCreateAWS(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceAWSResource, "id"),
					resource.TestCheckResourceAttr(ceAWSResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAWSResource, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAWSResource, "name", "resource_test_terraform_aws_export"),
					resource.TestCheckResourceAttr(ceAWSResource, "description", "resource test aws export"),
					resource.TestCheckResourceAttr(ceAWSResource, "plan_id", "9948"),
					resource.TestCheckResourceAttr(ceAWSResource, "cloud_provider", "aws"),
					resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.apply_bgp", "true"),
					resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.use_bgp_device_id", "1234"),
					resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.device_bgp_type", "router"),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.bucket", "resource-terraform-aws-bucket"),
					resource.TestCheckResourceAttr(
						ceAWSResource, "aws.0.iam_role_arn", "arn:aws:iam::003740049406:role/trafficTerraformIngestRole",
					),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.region", "eu-central-1"),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.delete_after_read", "true"),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.multiple_buckets", "true"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportUpdateAWS(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceAWSResource, "id"),
					resource.TestCheckResourceAttr(ceAWSResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceAWSResource, "enabled", "false"),
					resource.TestCheckResourceAttr(ceAWSResource, "name", "resource_test_terraform_aws_export_updated"),
					resource.TestCheckResourceAttr(ceAWSResource, "description", "resource test aws export updated"),
					resource.TestCheckResourceAttr(ceAWSResource, "plan_id", "3333"),
					resource.TestCheckResourceAttr(ceAWSResource, "cloud_provider", "aws"),
					resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.apply_bgp", "false"),
					resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.use_bgp_device_id", "4444"),
					resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.device_bgp_type", "dns"),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.bucket", "resource-terraform-aws-bucket-updated"),
					resource.TestCheckResourceAttr(
						ceAWSResource,
						"aws.0.iam_role_arn",
						"arn:aws:iam::003740049406:role/trafficTerraformIngestRole_updated",
					),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.region", "eu-central-1-updated"),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.delete_after_read", "false"),
					resource.TestCheckResourceAttr(ceAWSResource, "aws.0.multiple_buckets", "false"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportDestroy(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					testResourceDoesntExists(ceAWSResource),
				),
			},
		},
	})
}

func TestResourceCloudExportGCE(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccResourceCloudExportCreateGCE(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceGCEResource, "id"),
					resource.TestCheckResourceAttr(ceGCEResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceGCEResource, "enabled", "true"),
					resource.TestCheckResourceAttr(ceGCEResource, "name", "resource_test_terraform_gce_export"),
					resource.TestCheckResourceAttr(ceGCEResource, "description", "resource test gce export"),
					resource.TestCheckResourceAttr(ceGCEResource, "plan_id", "9948"),
					resource.TestCheckResourceAttr(ceGCEResource, "cloud_provider", "gce"),
					resource.TestCheckResourceAttr(ceGCEResource, "gce.0.project", "gce project"),
					resource.TestCheckResourceAttr(ceGCEResource, "gce.0.subscription", "gce subscription"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportUpdateGCE(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceGCEResource, "id"),
					resource.TestCheckResourceAttr(ceGCEResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceGCEResource, "enabled", "false"),
					resource.TestCheckResourceAttr(ceGCEResource, "name", "resource_test_terraform_gce_export_updated"),
					resource.TestCheckResourceAttr(ceGCEResource, "description", "resource test gce export updated"),
					resource.TestCheckResourceAttr(ceGCEResource, "plan_id", "3333"),
					resource.TestCheckResourceAttr(ceGCEResource, "cloud_provider", "gce"),
					resource.TestCheckResourceAttr(ceGCEResource, "gce.0.project", "gce project updated"),
					resource.TestCheckResourceAttr(ceGCEResource, "gce.0.subscription", "gce subscription updated"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportDestroy(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					testResourceDoesntExists(ceGCEResource),
				),
			},
		},
	})
}

func TestResourceCloudExportIBM(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccResourceCloudExportCreateIBM(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceIBMResource, "id"),
					resource.TestCheckResourceAttr(ceIBMResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceIBMResource, "enabled", "true"),
					resource.TestCheckResourceAttr(ceIBMResource, "name", "resource_test_terraform_ibm_export"),
					resource.TestCheckResourceAttr(ceIBMResource, "description", "resource test ibm export"),
					resource.TestCheckResourceAttr(ceIBMResource, "plan_id", "9948"),
					resource.TestCheckResourceAttr(ceIBMResource, "cloud_provider", "ibm"),
					resource.TestCheckResourceAttr(ceIBMResource, "ibm.0.bucket", "ibm-bucket"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportUpdateIBM(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceIBMResource, "id"),
					resource.TestCheckResourceAttr(ceIBMResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceIBMResource, "enabled", "false"),
					resource.TestCheckResourceAttr(ceIBMResource, "name", "resource_test_terraform_ibm_export_updated"),
					resource.TestCheckResourceAttr(ceIBMResource, "description", "resource test ibm export updated"),
					resource.TestCheckResourceAttr(ceIBMResource, "plan_id", "3333"),
					resource.TestCheckResourceAttr(ceIBMResource, "cloud_provider", "ibm"),
					resource.TestCheckResourceAttr(ceIBMResource, "ibm.0.bucket", "ibm-bucket-updated"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportDestroy(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					testResourceDoesntExists(ceIBMResource),
				),
			},
		},
	})
}

func TestResourceCloudExportAzure(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccResourceCloudExportCreateAzure(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceAzureResource, "id"),
					resource.TestCheckResourceAttr(ceAzureResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAzureResource, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAzureResource, "name", "resource_test_terraform_azure_export"),
					resource.TestCheckResourceAttr(ceAzureResource, "description", "resource test azure export"),
					resource.TestCheckResourceAttr(ceAzureResource, "plan_id", "9948"),
					resource.TestCheckResourceAttr(ceAzureResource, "cloud_provider", "azure"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.location", "centralus"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.resource_group", "traffic-generator"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.storage_account", "kentikstorage"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.subscription_id", "7777"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.security_principal_enabled", "true"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportUpdateAzure(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceAzureResource, "id"),
					resource.TestCheckResourceAttr(ceAzureResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceAzureResource, "enabled", "false"),
					resource.TestCheckResourceAttr(ceAzureResource, "name", "resource_test_terraform_azure_export_updated"),
					resource.TestCheckResourceAttr(ceAzureResource, "description", "resource test azure export updated"),
					resource.TestCheckResourceAttr(ceAzureResource, "plan_id", "3333"),
					resource.TestCheckResourceAttr(ceAzureResource, "cloud_provider", "azure"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.location", "centralus-updated"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.resource_group", "traffic-generator-updated"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.storage_account", "kentikstorage-updated"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.subscription_id", "8888"),
					resource.TestCheckResourceAttr(ceAzureResource, "azure.0.security_principal_enabled", "false"),
				),
			},
			{
				Config: makeTestAccResourceCloudExportDestroy(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					testResourceDoesntExists(ceAzureResource),
				),
			},
		},
	})
}

func testResourceDoesntExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		_, exists := s.RootModule().Resources[name]
		if exists {
			return fmt.Errorf("Resource %q found when not expected", name)
		}

		return nil
	}
}

const (
	ceAWSResource   = "kentik-cloudexport_item.test_aws"
	ceAzureResource = "kentik-cloudexport_item.test_azure"
	ceGCEResource   = "kentik-cloudexport_item.test_gce"
	ceIBMResource   = "kentik-cloudexport_item.test_ibm"
)

func makeTestAccResourceCloudExportCreateAWS(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_aws" {
			name= "resource_test_terraform_aws_export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "resource test aws export"
			plan_id= "9948"
			cloud_provider= "aws"
			bgp {
				apply_bgp= true
				use_bgp_device_id= "1234"
				device_bgp_type= "router"
			}
			aws {
				bucket= "resource-terraform-aws-bucket"
				iam_role_arn= "arn:aws:iam::003740049406:role/trafficTerraformIngestRole"
				region= "eu-central-1"
				delete_after_read= true
				multiple_buckets= true
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportUpdateAWS(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_aws" {
			name= "resource_test_terraform_aws_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=false
			description= "resource test aws export updated"
			plan_id= "3333"
			cloud_provider= "aws"
			bgp {
				apply_bgp= false
				use_bgp_device_id= "4444"
				device_bgp_type= "dns"
			}
			aws {
				bucket= "resource-terraform-aws-bucket-updated"
				iam_role_arn= "arn:aws:iam::003740049406:role/trafficTerraformIngestRole_updated"
				region= "eu-central-1-updated"
				delete_after_read= false
				multiple_buckets= false
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportCreateGCE(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_gce" {
			name= "resource_test_terraform_gce_export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "resource test gce export"
			plan_id= "9948"
			cloud_provider= "gce"
			gce {
				project= "gce project"
				subscription= "gce subscription"
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportUpdateGCE(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_gce" {
			name= "resource_test_terraform_gce_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=false
			description= "resource test gce export updated"
			plan_id= "3333"
			cloud_provider= "gce"
			gce {
				project= "gce project updated"
				subscription= "gce subscription updated"
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportCreateIBM(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_ibm" {
			name= "resource_test_terraform_ibm_export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "resource test ibm export"
			plan_id= "9948"
			cloud_provider= "ibm"
			ibm {
				bucket= "ibm-bucket"
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportUpdateIBM(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_ibm" {
			name= "resource_test_terraform_ibm_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=false
			description= "resource test ibm export updated"
			plan_id= "3333"
			cloud_provider= "ibm"
			ibm {
				bucket= "ibm-bucket-updated"
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportCreateAzure(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_azure" {
			name= "resource_test_terraform_azure_export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "resource test azure export"
			plan_id= "9948"
			cloud_provider= "azure"
			azure {
				location= "centralus"
				resource_group= "traffic-generator"
				storage_account= "kentikstorage"
				subscription_id= "7777"
				security_principal_enabled=true
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportUpdateAzure(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		
		resource "kentik-cloudexport_item" "test_azure" {
			name= "resource_test_terraform_azure_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=false
			description= "resource test azure export updated"
			plan_id= "3333"
			cloud_provider= "azure"
			azure {
				location= "centralus-updated"
				resource_group= "traffic-generator-updated"
				storage_account= "kentikstorage-updated"
				subscription_id= "8888"
				security_principal_enabled=false
			}
		  }
		`,
		apiURL,
	)
}

func makeTestAccResourceCloudExportDestroy(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		`, apiURL,
	)
}
