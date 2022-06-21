package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Note: we only check the user-provided values as we don't control the server-provided ones

const (
	ceAWSResource   = "kentik-cloudexport_item.test_aws"
	ceAzureResource = "kentik-cloudexport_item.test_azure"
	ceGCEResource   = "kentik-cloudexport_item.test_gce"
	ceIBMResource   = "kentik-cloudexport_item.test_ibm"
)

func TestResourceCloudExportAWS(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestResourceCloudExportCreateAWS(server.URL()),
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
				Config: makeTestResourceCloudExportUpdateAWS(server.URL()),
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
				Config: makeTestResourceCloudExportDestroy(server.URL()),
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
				Config: makeTestResourceCloudExportCreateGCE(server.URL()),
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
				Config: makeTestResourceCloudExportUpdateGCE(server.URL()),
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
				Config: makeTestResourceCloudExportDestroy(server.URL()),
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
				Config: makeTestResourceCloudExportCreateIBM(server.URL()),
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
				Config: makeTestResourceCloudExportUpdateIBM(server.URL()),
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
				Config: makeTestResourceCloudExportDestroy(server.URL()),
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
				Config: makeTestResourceCloudExportCreateAzure(server.URL()),
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
				Config: makeTestResourceCloudExportUpdateAzure(server.URL()),
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
				Config: makeTestResourceCloudExportDestroy(server.URL()),
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

func makeTestResourceCloudExportCreateAWS(apiURL string) string {
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

func makeTestResourceCloudExportUpdateAWS(apiURL string) string {
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

func makeTestResourceCloudExportCreateGCE(apiURL string) string {
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

func makeTestResourceCloudExportUpdateGCE(apiURL string) string {
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

func makeTestResourceCloudExportCreateIBM(apiURL string) string {
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

func makeTestResourceCloudExportUpdateIBM(apiURL string) string {
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

func makeTestResourceCloudExportCreateAzure(apiURL string) string {
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

func makeTestResourceCloudExportUpdateAzure(apiURL string) string {
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

func makeTestResourceCloudExportDestroy(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
		}
		`, apiURL,
	)
}

func TestAccResourceCloudExportAWS(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccResourceCloudExportCreateAWS(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceAWSResource, "id"),
						resource.TestCheckResourceAttr(ceAWSResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
						resource.TestCheckResourceAttr(ceAWSResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceAWSResource, "name", fmt.Sprintf("%s-aws-export", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSResource, "description", fmt.Sprintf("%s-description", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSResource, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceAWSResource, "cloud_provider", "aws"),
						resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.apply_bgp", "true"),
						resource.TestCheckResourceAttr(
							ceAWSResource,
							"bgp.0.use_bgp_device_id",
							fmt.Sprintf("%s-bgp-id", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.device_bgp_type", "router"),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.bucket", fmt.Sprintf("%s-aws-bucket", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAWSResource,
							"aws.0.iam_role_arn",
							fmt.Sprintf("%s-iam-role-arn", getAccTestPrefix()),
						),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.region", "eu-central-1"),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.delete_after_read", "true"),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.multiple_buckets", "true"),
					),
				},
				{
					Config: makeTestAccResourceCloudExportUpdateAWS(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceAWSResource, "id"),
						resource.TestCheckResourceAttr(ceAWSResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
						resource.TestCheckResourceAttr(ceAWSResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceAWSResource, "name", fmt.Sprintf("%s-aws-export-update", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAWSResource,
							"description",
							fmt.Sprintf("%s-description-update", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSResource, "cloud_provider", "aws"),
						resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.apply_bgp", "false"),
						resource.TestCheckResourceAttr(
							ceAWSResource,
							"bgp.0.use_bgp_device_id",
							fmt.Sprintf("%s-bgp-id-update", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSResource, "bgp.0.device_bgp_type", "dns"),
						resource.TestCheckResourceAttr(
							ceAWSResource,
							"aws.0.bucket",
							fmt.Sprintf("%s-aws-bucket-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAWSResource,
							"aws.0.iam_role_arn",
							fmt.Sprintf("%s-iam-role-arn-updated", getAccTestPrefix()),
						),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.region", "eu-central-1-updated"),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.delete_after_read", "false"),
						resource.TestCheckResourceAttr(ceAWSResource, "aws.0.multiple_buckets", "false"),
					),
				},
			},
		})
	}
}

func TestAccResourceCloudExportGCE(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccResourceCloudExportCreateGCE(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceGCEResource, "id"),
						resource.TestCheckResourceAttr(ceGCEResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
						resource.TestCheckResourceAttr(ceGCEResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceGCEResource, "name", fmt.Sprintf("%s-gce-export", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceGCEResource, "description", fmt.Sprintf("%s-description", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceGCEResource, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceGCEResource, "cloud_provider", "gce"),
						resource.TestCheckResourceAttr(ceGCEResource, "gce.0.project", fmt.Sprintf("%s-gce project", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceGCEResource,
							"gce.0.subscription",
							fmt.Sprintf("%s-gce subscription", getAccTestPrefix())),
					),
				},
				{
					Config: makeTestAccResourceCloudExportUpdateGCE(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceGCEResource, "id"),
						resource.TestCheckResourceAttr(ceGCEResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
						resource.TestCheckResourceAttr(ceGCEResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceGCEResource, "name", fmt.Sprintf("%s-gce-export-update", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceGCEResource,
							"description",
							fmt.Sprintf("%s-description-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceGCEResource, "cloud_provider", "gce"),
						resource.TestCheckResourceAttr(
							ceGCEResource,
							"gce.0.project",
							fmt.Sprintf("%s-gce project updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceGCEResource,
							"gce.0.subscription",
							fmt.Sprintf("%s-gce subscription updated", getAccTestPrefix())),
					),
				},
			},
		})
	}
}

func TestAccResourceCloudExportIBM(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccResourceCloudExportCreateIBM(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceIBMResource, "id"),
						resource.TestCheckResourceAttr(ceIBMResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
						resource.TestCheckResourceAttr(ceIBMResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceIBMResource, "name", fmt.Sprintf("%s-ibm-export", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceIBMResource, "description", fmt.Sprintf("%s-description", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceIBMResource, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceIBMResource, "cloud_provider", "ibm"),
						resource.TestCheckResourceAttr(ceIBMResource, "ibm.0.bucket", fmt.Sprintf("%s-ibm-bucket", getAccTestPrefix())),
					),
				},
				{
					Config: makeTestAccResourceCloudExportUpdateIBM(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceIBMResource, "id"),
						resource.TestCheckResourceAttr(ceIBMResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
						resource.TestCheckResourceAttr(ceIBMResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceIBMResource, "name", fmt.Sprintf("%s-ibm-export-update", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceIBMResource,
							"description",
							fmt.Sprintf("%s-description-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceIBMResource, "cloud_provider", "ibm"),
						resource.TestCheckResourceAttr(
							ceIBMResource,
							"ibm.0.bucket",
							fmt.Sprintf("%s-ibm-bucket-updated", getAccTestPrefix())),
					),
				},
			},
		})
	}
}

func TestAccResourceCloudExportAzure(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccResourceCloudExportCreateAzure(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceAzureResource, "id"),
						resource.TestCheckResourceAttr(ceAzureResource, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
						resource.TestCheckResourceAttr(ceAzureResource, "enabled", "true"),
						resource.TestCheckResourceAttr(ceAzureResource, "name", fmt.Sprintf("%s-azure-export", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureResource, "description", fmt.Sprintf("%s-description", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureResource, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceAzureResource, "cloud_provider", "azure"),
						resource.TestCheckResourceAttr(ceAzureResource, "azure.0.location", "centralus"),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"azure.0.resource_group",
							fmt.Sprintf("%s-traffic-generator", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"azure.0.storage_account",
							fmt.Sprintf("%s-kentikstorage", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"azure.0.subscription_id",
							fmt.Sprintf("%s-sub-id", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureResource, "azure.0.security_principal_enabled", "true"),
					),
				},
				{
					Config: makeTestAccResourceCloudExportUpdateAzure(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(ceAzureResource, "id"),
						resource.TestCheckResourceAttr(ceAzureResource, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
						resource.TestCheckResourceAttr(ceAzureResource, "enabled", "true"),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"name",
							fmt.Sprintf("%s-azure-export-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"description",
							fmt.Sprintf("%s-description-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureResource, "cloud_provider", "azure"),
						resource.TestCheckResourceAttr(ceAzureResource, "azure.0.location", "centralus-updated"),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"azure.0.resource_group",
							fmt.Sprintf("%s-traffic-generator-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"azure.0.storage_account",
							fmt.Sprintf("%s-kentik-storage-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAzureResource,
							"azure.0.subscription_id",
							fmt.Sprintf("%s-sub-id-updated", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureResource, "azure.0.security_principal_enabled", "false"),
					),
				},
			},
		})
	}
}

func makeTestAccResourceCloudExportCreateAWS() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_aws" {
			name= "%[1]s-aws-export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "%[1]s-description"
			plan_id= %[2]s
			cloud_provider= "aws"
			bgp {
				apply_bgp= true
				use_bgp_device_id= "%[1]s-bgp-id"
				device_bgp_type= "router"
			}
			aws {
				bucket= "%[1]s-aws-bucket"
				iam_role_arn= "%[1]s-iam-role-arn"
				region= "eu-central-1"
				delete_after_read= true
				multiple_buckets= true
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportUpdateAWS() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_aws" {
			name= "%[1]s-aws-export-update"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
			description= "%[1]s-description-update"
			plan_id= %[2]s
			cloud_provider= "aws"
			bgp {
				apply_bgp= false
				use_bgp_device_id= "%[1]s-bgp-id-update"
				device_bgp_type= "dns"
			}
			aws {
				bucket= "%[1]s-aws-bucket-updated"
				iam_role_arn= "%[1]s-iam-role-arn-updated"
				region= "eu-central-1-updated"
				delete_after_read= false
				multiple_buckets= false
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportCreateGCE() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_gce" {
			name= "%[1]s-gce-export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "%[1]s-description"
			plan_id= %[2]s
			cloud_provider= "gce"
			gce {
				project= "%[1]s-gce project"
				subscription= "%[1]s-gce subscription"
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportUpdateGCE() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_gce" {
			name= "%[1]s-gce-export-update"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
			description= "%[1]s-description-updated"
			plan_id= %[2]s
			cloud_provider= "gce"
			gce {
				project= "%[1]s-gce project updated"
				subscription= "%[1]s-gce subscription updated"
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportCreateIBM() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_ibm" {
			name= "%[1]s-ibm-export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "%[1]s-description"
			plan_id= %[2]s
			cloud_provider= "ibm"
			ibm {
				bucket= "%[1]s-ibm-bucket"
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportUpdateIBM() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_ibm" {
			name= "%[1]s-ibm-export-update"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
			description= "%[1]s-description-updated"
			plan_id= %[2]s
			cloud_provider= "ibm"
			ibm {
				bucket= "%[1]s-ibm-bucket-updated"
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportCreateAzure() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_azure" {
			name= "%[1]s-azure-export"
			type= "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"
			enabled=true
			description= "%[1]s-description"
			plan_id= %[2]s
			cloud_provider= "azure"
			azure {
				location= "centralus"
				resource_group= "%[1]s-traffic-generator"
				storage_account= "%[1]s-kentikstorage"
				subscription_id= "%[1]s-sub-id"
				security_principal_enabled=true
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}

func makeTestAccResourceCloudExportUpdateAzure() string {
	return fmt.Sprintf(`
		resource "kentik-cloudexport_item" "test_azure" {
			name= "%[1]s-azure-export-updated"
			type= "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"
			enabled=true
			description= "%[1]s-description-updated"
			plan_id= %[2]s
			cloud_provider= "azure"
			azure {
				location= "centralus-updated"
				resource_group= "%[1]s-traffic-generator-updated"
				storage_account= "%[1]s-kentik-storage-updated"
				subscription_id= "%[1]s-sub-id-updated"
				security_principal_enabled=false
			}
		  }
		`, getAccTestPrefix(), getKentikPlanIDAccTests())
}
