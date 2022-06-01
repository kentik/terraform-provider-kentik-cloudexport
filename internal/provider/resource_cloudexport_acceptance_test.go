package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudExportItemAWS(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudExportCreateAWS(),
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
				Config: testAccResourceCloudExportUpdateAWS(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceAWSResource, "id"),
					resource.TestCheckResourceAttr(ceAWSResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceAWSResource, "enabled", "true"),
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
		},
	})
}

func TestAccResourceCloudExportGCE(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudExportCreateGCE(),
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
				Config: testAccResourceCloudExportUpdateGCE(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceGCEResource, "id"),
					resource.TestCheckResourceAttr(ceGCEResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceGCEResource, "enabled", "true"),
					resource.TestCheckResourceAttr(ceGCEResource, "name", "resource_test_terraform_gce_export_updated"),
					resource.TestCheckResourceAttr(ceGCEResource, "description", "resource test gce export updated"),
					resource.TestCheckResourceAttr(ceGCEResource, "plan_id", "3333"),
					resource.TestCheckResourceAttr(ceGCEResource, "cloud_provider", "gce"),
					resource.TestCheckResourceAttr(ceGCEResource, "gce.0.project", "gce project updated"),
					resource.TestCheckResourceAttr(ceGCEResource, "gce.0.subscription", "gce subscription updated"),
				),
			},
		},
	})
}

func TestAccResourceCloudExportIBM(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudExportCreateIBM(),
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
				Config: testAccResourceCloudExportUpdateIBM(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceIBMResource, "id"),
					resource.TestCheckResourceAttr(ceIBMResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceIBMResource, "enabled", "true"),
					resource.TestCheckResourceAttr(ceIBMResource, "name", "resource_test_terraform_ibm_export_updated"),
					resource.TestCheckResourceAttr(ceIBMResource, "description", "resource test ibm export updated"),
					resource.TestCheckResourceAttr(ceIBMResource, "plan_id", "3333"),
					resource.TestCheckResourceAttr(ceIBMResource, "cloud_provider", "ibm"),
					resource.TestCheckResourceAttr(ceIBMResource, "ibm.0.bucket", "ibm-bucket-updated"),
				),
			},
		},
	})
}

func TestAccResourceCloudExportAzure(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudExportCreateAzure(),
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
				Config: testAccResourceCloudExportUpdateAzure(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(ceAzureResource, "id"),
					resource.TestCheckResourceAttr(ceAzureResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceAzureResource, "enabled", "true"),
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
		},
	})
}

func testAccResourceCloudExportCreateAWS() string {
	return `
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
		`
}

func testAccResourceCloudExportUpdateAWS() string {
	return `
		resource "kentik-cloudexport_item" "test_aws" {
			name= "resource_test_terraform_aws_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
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
		`
}

func testAccResourceCloudExportCreateGCE() string {
	return `
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
		`
}

func testAccResourceCloudExportUpdateGCE() string {
	return `
		resource "kentik-cloudexport_item" "test_gce" {
			name= "resource_test_terraform_gce_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
			description= "resource test gce export updated"
			plan_id= "3333"
			cloud_provider= "gce"
			gce {
				project= "gce project updated"
				subscription= "gce subscription updated"
			}
		  }
		`
}

func testAccResourceCloudExportCreateIBM() string {
	return `
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
		`
}

func testAccResourceCloudExportUpdateIBM() string {
	return `
		resource "kentik-cloudexport_item" "test_ibm" {
			name= "resource_test_terraform_ibm_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
			description= "resource test ibm export updated"
			plan_id= "3333"
			cloud_provider= "ibm"
			ibm {
				bucket= "ibm-bucket-updated"
			}
		  }
		`
}

func testAccResourceCloudExportCreateAzure() string {
	return `
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
		`
}

func testAccResourceCloudExportUpdateAzure() string {
	return `
		resource "kentik-cloudexport_item" "test_azure" {
			name= "resource_test_terraform_azure_export_updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
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
		`
}
