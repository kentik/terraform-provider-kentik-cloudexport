package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/kentik/community_sdk_golang/kentikapi"
	"github.com/kentik/community_sdk_golang/kentikapi/models"
	"github.com/stretchr/testify/assert"
)

const (
	ceAWSDS   = "data.kentik-cloudexport_item.aws"
	ceAzureDS = "data.kentik-cloudexport_item.azure"
	ceGCPDS   = "data.kentik-cloudexport_item.gce"
	ceIBMDS   = "data.kentik-cloudexport_item.ibm"
)

func TestDataSourceCloudExportItemAWS(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestCloudExportDataSourceItems(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceAWSDS, "id", "1"),
					resource.TestCheckResourceAttr(ceAWSDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAWSDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "name", "test_terraform_aws_export"),
					resource.TestCheckResourceAttr(ceAWSDS, "description", "terraform aws cloud export"),
					resource.TestCheckResourceAttr(ceAWSDS, "plan_id", "11467"),
					resource.TestCheckResourceAttr(ceAWSDS, "cloud_provider", "aws"),
					resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.apply_bgp", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.use_bgp_device_id", "dummy-device-id"),
					resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.device_bgp_type", "dummy-device-bgp-type"),
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

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestCloudExportDataSourceItems(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceGCPDS, "id", "2"),
					resource.TestCheckResourceAttr(ceGCPDS, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceGCPDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceGCPDS, "name", "test_terraform_gce_export"),
					resource.TestCheckResourceAttr(ceGCPDS, "description", "terraform gce cloud export"),
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

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestCloudExportDataSourceItems(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceIBMDS, "id", "3"),
					resource.TestCheckResourceAttr(ceIBMDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceIBMDS, "enabled", "false"),
					resource.TestCheckResourceAttr(ceIBMDS, "name", "test_terraform_ibm_export"),
					resource.TestCheckResourceAttr(ceIBMDS, "description", "terraform ibm cloud export"),
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

func TestDataSourceCloudExportItemAzure(t *testing.T) {
	t.Parallel()

	server := newTestAPIServer(t, makeInitialCloudExports())
	server.Start()
	defer server.Stop()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestCloudExportDataSourceItems(server.URL()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceAzureDS, "id", "4"),
					resource.TestCheckResourceAttr(ceAzureDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAzureDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAzureDS, "name", "test_terraform_azure_export"),
					resource.TestCheckResourceAttr(ceAzureDS, "description", "terraform azure cloud export"),
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

func makeTestCloudExportDataSourceItems(apiURL string) string {
	return fmt.Sprintf(`
		provider "kentik-cloudexport" {
			apiurl = "%v"
			email = "joe.doe@example.com"
			token = "dummy-token"
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
	`,
		apiURL,
	)
}

func TestAccDataSourceCloudExportItemAWS(t *testing.T) {
	ce, err := createTestAccCloudExportItemAWS()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, deleteTestAccCloudExportItem(ce))
	}()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderAWS, ce),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceAWSDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceAWSDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "name", "acc_test_terraform_aws_export"),
					resource.TestCheckResourceAttr(ceAWSDS, "description", "terraform aws cloud export"),
					resource.TestCheckResourceAttr(ceAWSDS, "plan_id", "7512"),
					resource.TestCheckResourceAttr(ceAWSDS, "cloud_provider", "aws"),
					resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.apply_bgp", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.use_bgp_device_id", "dummy-device-id"),
					resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.device_bgp_type", "dummy-device-bgp-type"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.bucket", "terraform-aws-bucket"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.iam_role_arn", "dummy-iam-role-arn"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.region", "us-east-2"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.delete_after_read", "true"),
					resource.TestCheckResourceAttr(ceAWSDS, "aws.0.multiple_buckets", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceCloudExportItemGCE(t *testing.T) {
	ce, err := createTestAccCloudExportItemGCE()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, deleteTestAccCloudExportItem(ce))
	}()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderGCE, ce),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceGCPDS, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceGCPDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceGCPDS, "name", "acc_test_terraform_gce_export"),
					resource.TestCheckResourceAttr(ceGCPDS, "description", "terraform gce cloud export"),
					resource.TestCheckResourceAttr(ceGCPDS, "plan_id", "7512"),
					resource.TestCheckResourceAttr(ceGCPDS, "cloud_provider", "gce"),
					resource.TestCheckResourceAttr(ceGCPDS, "gce.0.project", "project gce"),
					resource.TestCheckResourceAttr(ceGCPDS, "gce.0.subscription", "subscription gce"),
					resource.TestCheckResourceAttr(ceGCPDS, "bgp.0.apply_bgp", "true"),
					resource.TestCheckResourceAttr(ceGCPDS, "bgp.0.use_bgp_device_id", "dummy-device-id"),
					resource.TestCheckResourceAttr(ceGCPDS, "bgp.0.device_bgp_type", "dummy-device-bgp-type"),
				),
			},
		},
	})
}

func TestAccDataSourceCloudExportItemIBM(t *testing.T) {
	ce, err := createTestAccCloudExportItemIBM()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, deleteTestAccCloudExportItem(ce))
	}()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderIBM, ce),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceIBMDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
					resource.TestCheckResourceAttr(ceIBMDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceIBMDS, "name", "acc_test_terraform_ibm_export"),
					resource.TestCheckResourceAttr(ceIBMDS, "description", "terraform ibm cloud export"),
					resource.TestCheckResourceAttr(ceIBMDS, "plan_id", "7512"),
					resource.TestCheckResourceAttr(ceIBMDS, "cloud_provider", "ibm"),
					resource.TestCheckResourceAttr(ceIBMDS, "ibm.0.bucket", "terraform-ibm-bucket"),
					resource.TestCheckResourceAttr(ceIBMDS, "bgp.0.apply_bgp", "true"),
					resource.TestCheckResourceAttr(ceIBMDS, "bgp.0.use_bgp_device_id", "dummy-device-id"),
					resource.TestCheckResourceAttr(ceIBMDS, "bgp.0.device_bgp_type", "dummy-device-bgp-type"),
				),
			},
		},
	})
}

func TestAccDataSourceCloudExportItemAzure(t *testing.T) {
	ce, err := createTestAccCloudExportItemAzure()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, deleteTestAccCloudExportItem(ce))
	}()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderAzure, ce),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ceAzureDS, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
					resource.TestCheckResourceAttr(ceAzureDS, "enabled", "true"),
					resource.TestCheckResourceAttr(ceAzureDS, "name", "acc_test_terraform_azure_export"),
					resource.TestCheckResourceAttr(ceAzureDS, "description", "terraform azure cloud export"),
					resource.TestCheckResourceAttr(ceAzureDS, "plan_id", "7512"),
					resource.TestCheckResourceAttr(ceAzureDS, "cloud_provider", "azure"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.location", "centralus"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.resource_group", "traffic-generator"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.storage_account", "dummy-sa"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.subscription_id", "dummy-sid"),
					resource.TestCheckResourceAttr(ceAzureDS, "azure.0.security_principal_enabled", "true"),
					resource.TestCheckNoResourceAttr(ceAzureDS, "bgp.0.apply_bgp"),
				),
			},
		},
	})
}

func makeTestAccCloudExportDataSourceItems(provider string, ce *models.CloudExport) string {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return ""
	}
	return fmt.Sprintf(`
		data "kentik-cloudexport_item" "%v" {
			id = "%v"
		}
	`,
		provider, ce.ID,
	)
}

func createTestAccCloudExportItemAWS() (*models.CloudExport, error) {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return nil, nil
	}
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewAWSCloudExport(models.CloudExportAWSRequiredFields{
		Name:   "acc_test_terraform_aws_export",
		PlanID: "7512",
		AWSProperties: models.AWSPropertiesRequiredFields{
			Bucket: "terraform-aws-bucket",
		},
	})
	ce.Type = models.CloudExportTypeKentikManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = "terraform aws cloud export"
	ce.GetAWSProperties().IAMRoleARN = "dummy-iam-role-arn"
	ce.GetAWSProperties().Region = "us-east-2"
	ce.GetAWSProperties().DeleteAfterRead = pointer.ToBool(true)
	ce.GetAWSProperties().MultipleBuckets = pointer.ToBool(true)
	ce.BGP = &models.BGPProperties{
		ApplyBGP:       pointer.ToBool(true),
		UseBGPDeviceID: "dummy-device-id",
		DeviceBGPType:  "dummy-device-bgp-type",
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func createTestAccCloudExportItemGCE() (*models.CloudExport, error) {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return nil, nil
	}
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewGCECloudExport(models.CloudExportGCERequiredFields{
		Name:   "acc_test_terraform_gce_export",
		PlanID: "7512",
		GCEProperties: models.GCEPropertiesRequiredFields{
			Project:      "project gce",
			Subscription: "subscription gce",
		},
	})
	ce.Type = models.CloudExportTypeCustomerManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = "terraform gce cloud export"
	ce.BGP = &models.BGPProperties{
		ApplyBGP:       pointer.ToBool(true),
		UseBGPDeviceID: "dummy-device-id",
		DeviceBGPType:  "dummy-device-bgp-type",
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func createTestAccCloudExportItemIBM() (*models.CloudExport, error) {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return nil, nil
	}
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewIBMCloudExport(models.CloudExportIBMRequiredFields{
		Name:   "acc_test_terraform_ibm_export",
		PlanID: "7512",
		IBMProperties: models.IBMPropertiesRequiredFields{
			Bucket: "terraform-ibm-bucket",
		},
	})
	ce.Type = models.CloudExportTypeKentikManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = "terraform ibm cloud export"
	ce.BGP = &models.BGPProperties{
		ApplyBGP:       pointer.ToBool(true),
		UseBGPDeviceID: "dummy-device-id",
		DeviceBGPType:  "dummy-device-bgp-type",
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func createTestAccCloudExportItemAzure() (*models.CloudExport, error) {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return nil, nil
	}
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewAzureCloudExport(models.CloudExportAzureRequiredFields{
		Name:   "acc_test_terraform_azure_export",
		PlanID: "7512",
		AzureProperties: models.AzurePropertiesRequiredFields{
			Location:       "centralus",
			ResourceGroup:  "traffic-generator",
			StorageAccount: "dummy-sa",
			SubscriptionID: "dummy-sid",
		},
	})
	ce.Type = models.CloudExportTypeCustomerManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = "terraform azure cloud export"
	ce.GetAzureProperties().SecurityPrincipalEnabled = pointer.ToBool(true)
	ce.BGP = &models.BGPProperties{
		ApplyBGP: pointer.ToBool(false),
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func deleteTestAccCloudExportItem(ce *models.CloudExport) error {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return nil
	}
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return err
	}
	err = client.CloudExports.Delete(ctx, ce.ID)
	if err != nil {
		return fmt.Errorf("client.CloudExports.Delete: %w", err)
	}
	return nil
}

func newClient() (*kentikapi.Client, error) {
	authEmail, _ := os.LookupEnv("KTAPI_AUTH_EMAIL")
	authToken, _ := os.LookupEnv("KTAPI_AUTH_TOKEN")
	client, err := kentikapi.NewClient(kentikapi.Config{
		AuthEmail:   authEmail,
		AuthToken:   authToken,
		LogPayloads: false,
	})
	if err != nil {
		return nil, fmt.Errorf("newClient: %w", err)
	}
	return client, nil
}
