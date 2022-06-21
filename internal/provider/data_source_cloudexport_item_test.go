package provider_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/kentik/community_sdk_golang/kentikapi"
	"github.com/kentik/community_sdk_golang/kentikapi/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ceAWSDS   = "data.kentik-cloudexport_item.aws"
	ceAzureDS = "data.kentik-cloudexport_item.azure"
	ceGCPDS   = "data.kentik-cloudexport_item.gce"
	ceIBMDS   = "data.kentik-cloudexport_item.ibm"
)

//nolint: gochecknoinits
func init() {
	resource.AddTestSweepers("kentik_tf_integ_test", &resource.Sweeper{
		Name: "kentik_tf_integ_test",
		F: func(region string) error {
			ctx := context.Background()
			client, err := newClient()
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}

			getAll, err := client.CloudExports.GetAll(ctx)
			if err != nil {
				return fmt.Errorf("error getting CloudExports: %s", err)
			}

			for _, ce := range getAll.CloudExports {
				if strings.HasPrefix(ce.Name, getAccTestPrefix()) {
					err := client.CloudExports.Delete(ctx, ce.ID)
					if err != nil {
						log.Printf("Error destroying %s during sweep: %s", ce.Name, err)
					}
				}
			}
			return nil
		},
	})
}

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
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		ce, err := createTestAccCloudExportItemAWS()
		assert.NoError(t, err)

		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderAWS, ce),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ceAWSDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
						resource.TestCheckResourceAttr(ceAWSDS, "enabled", "true"),
						resource.TestCheckResourceAttr(ceAWSDS, "name", fmt.Sprintf("%s-aws-export-item", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSDS, "description", fmt.Sprintf("%s-description-aws", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSDS, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceAWSDS, "cloud_provider", "aws"),
						resource.TestCheckResourceAttr(ceAWSDS, "bgp.0.apply_bgp", "true"),
						resource.TestCheckResourceAttr(
							ceAWSDS,
							"bgp.0.use_bgp_device_id",
							fmt.Sprintf("%s-device-id", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAWSDS,
							"bgp.0.device_bgp_type",
							fmt.Sprintf("%s-device-bgp-type", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceAWSDS,
							"aws.0.bucket",
							fmt.Sprintf("%s-terraform-aws-bucket", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSDS, "aws.0.iam_role_arn", fmt.Sprintf("%s-iam-role-arn", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAWSDS, "aws.0.region", "us-east-2"),
						resource.TestCheckResourceAttr(ceAWSDS, "aws.0.delete_after_read", "true"),
						resource.TestCheckResourceAttr(ceAWSDS, "aws.0.multiple_buckets", "true"),
					),
				},
			},
		})
	}
}

func TestAccDataSourceCloudExportItemGCE(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		ce, err := createTestAccCloudExportItemGCE()
		assert.NoError(t, err)

		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderGCE, ce),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ceGCPDS, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
						resource.TestCheckResourceAttr(ceGCPDS, "enabled", "true"),
						resource.TestCheckResourceAttr(ceGCPDS, "name", fmt.Sprintf("%s-gce-export-item", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceGCPDS, "description", fmt.Sprintf("%s-description-gce", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceGCPDS, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceGCPDS, "cloud_provider", "gce"),
						resource.TestCheckResourceAttr(
							ceGCPDS,
							"gce.0.project",
							fmt.Sprintf("%s-gce-export-list-project gce", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceGCPDS,
							"gce.0.subscription",
							fmt.Sprintf("%s-subscription gce", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceGCPDS,
							"bgp.0.apply_bgp",
							"true"),
						resource.TestCheckResourceAttr(
							ceGCPDS,
							"bgp.0.use_bgp_device_id",
							fmt.Sprintf("%s-device-id", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceGCPDS,
							"bgp.0.device_bgp_type",
							fmt.Sprintf("%s-device-bgp-type", getAccTestPrefix())),
					),
				},
			},
		})
	}
}

func TestAccDataSourceCloudExportItemIBM(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		ce, err := createTestAccCloudExportItemIBM()
		assert.NoError(t, err)

		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderIBM, ce),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ceIBMDS, "type", "CLOUD_EXPORT_TYPE_KENTIK_MANAGED"),
						resource.TestCheckResourceAttr(ceIBMDS, "enabled", "true"),
						resource.TestCheckResourceAttr(ceIBMDS, "name", fmt.Sprintf("%s-ibm-export-item", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceIBMDS, "description", fmt.Sprintf("%s-description-ibm", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceIBMDS, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceIBMDS, "cloud_provider", "ibm"),
						resource.TestCheckResourceAttr(
							ceIBMDS,
							"ibm.0.bucket",
							fmt.Sprintf("%s-terraform-ibm-bucket", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceIBMDS, "bgp.0.apply_bgp", "true"),
						resource.TestCheckResourceAttr(
							ceIBMDS,
							"bgp.0.use_bgp_device_id",
							fmt.Sprintf("%s-device-id", getAccTestPrefix())),
						resource.TestCheckResourceAttr(
							ceIBMDS,
							"bgp.0.device_bgp_type",
							fmt.Sprintf("%s-device-bgp-type", getAccTestPrefix())),
					),
				},
			},
		})
	}
}

func TestAccDataSourceCloudExportItemAzure(t *testing.T) {
	if skipIfNotAcceptance() {
		checkRequiredEnvVariables(t)
		ce, err := createTestAccCloudExportItemAzure()
		assert.NoError(t, err)

		resource.ParallelTest(t, resource.TestCase{
			ProviderFactories: providerFactories(),
			Steps: []resource.TestStep{
				{
					Config: makeTestAccCloudExportDataSourceItems(models.CloudProviderAzure, ce),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(ceAzureDS, "type", "CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED"),
						resource.TestCheckResourceAttr(ceAzureDS, "enabled", "true"),
						resource.TestCheckResourceAttr(ceAzureDS, "name", fmt.Sprintf("%s-azure-export-item", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureDS, "description", fmt.Sprintf("%s-description-azure", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureDS, "plan_id", getKentikPlanIDAccTests()),
						resource.TestCheckResourceAttr(ceAzureDS, "cloud_provider", "azure"),
						resource.TestCheckResourceAttr(ceAzureDS, "azure.0.location", "centralus"),
						resource.TestCheckResourceAttr(
							ceAzureDS,
							"azure.0.resource_group",
							fmt.Sprintf("%s-traffic-generator", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureDS, "azure.0.storage_account", fmt.Sprintf("%s-sa", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureDS, "azure.0.subscription_id", fmt.Sprintf("%s-sid", getAccTestPrefix())),
						resource.TestCheckResourceAttr(ceAzureDS, "azure.0.security_principal_enabled", "true"),
						resource.TestCheckNoResourceAttr(ceAzureDS, "bgp.0.apply_bgp"),
					),
				},
			},
		})
	}
}

func skipIfNotAcceptance() bool {
	_, accTest := os.LookupEnv(resource.EnvTfAcc)
	return accTest
}

func getAccTestPrefix() string {
	return fmt.Sprintf("kentik_tf_integ_test_%s", os.Getenv("TF_ACC_PREFIX"))
}

func checkRequiredEnvVariables(t *testing.T) {
	_, ok := os.LookupEnv("KTAPI_AUTH_EMAIL")
	require.True(t, ok, "KTAPI_AUTH_EMAIL env variable not set")
	_, ok = os.LookupEnv("KTAPI_AUTH_TOKEN")
	require.True(t, ok, "KTAPI_AUTH_TOKEN env variable not set")
	_, ok = os.LookupEnv("KTAPI_URL")
	require.True(t, ok, "KTAPI_URL env variable not set")
	_, ok = os.LookupEnv("KENTIK_PLAN_ID")
	require.True(t, ok, "KENTIK_PLAN_ID env variable not set")
}

func makeTestAccCloudExportDataSourceItems(provider string, ce *models.CloudExport) string {
	return fmt.Sprintf(`
		data "kentik-cloudexport_item" "%v" {
			id = "%v"
		}
	`,
		provider, ce.ID,
	)
}

func createTestAccCloudExportItemAWS() (*models.CloudExport, error) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewAWSCloudExport(models.CloudExportAWSRequiredFields{
		Name:   fmt.Sprintf("%s-aws-export-item", getAccTestPrefix()),
		PlanID: getKentikPlanIDAccTests(),
		AWSProperties: models.AWSPropertiesRequiredFields{
			Bucket: fmt.Sprintf("%s-terraform-aws-bucket", getAccTestPrefix()),
		},
	})
	ce.Type = models.CloudExportTypeKentikManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = fmt.Sprintf("%s-description-aws", getAccTestPrefix())
	ce.GetAWSProperties().IAMRoleARN = fmt.Sprintf("%s-iam-role-arn", getAccTestPrefix())
	ce.GetAWSProperties().Region = "us-east-2"
	ce.GetAWSProperties().DeleteAfterRead = pointer.ToBool(true)
	ce.GetAWSProperties().MultipleBuckets = pointer.ToBool(true)
	ce.BGP = &models.BGPProperties{
		ApplyBGP:       pointer.ToBool(true),
		UseBGPDeviceID: fmt.Sprintf("%s-device-id", getAccTestPrefix()),
		DeviceBGPType:  fmt.Sprintf("%s-device-bgp-type", getAccTestPrefix()),
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func createTestAccCloudExportItemGCE() (*models.CloudExport, error) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewGCECloudExport(models.CloudExportGCERequiredFields{
		Name:   fmt.Sprintf("%s-gce-export-item", getAccTestPrefix()),
		PlanID: getKentikPlanIDAccTests(),
		GCEProperties: models.GCEPropertiesRequiredFields{
			Project:      fmt.Sprintf("%s-gce-export-list-project gce", getAccTestPrefix()),
			Subscription: fmt.Sprintf("%s-subscription gce", getAccTestPrefix()),
		},
	})
	ce.Type = models.CloudExportTypeCustomerManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = fmt.Sprintf("%s-description-gce", getAccTestPrefix())
	ce.BGP = &models.BGPProperties{
		ApplyBGP:       pointer.ToBool(true),
		UseBGPDeviceID: fmt.Sprintf("%s-device-id", getAccTestPrefix()),
		DeviceBGPType:  fmt.Sprintf("%s-device-bgp-type", getAccTestPrefix()),
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func createTestAccCloudExportItemIBM() (*models.CloudExport, error) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewIBMCloudExport(models.CloudExportIBMRequiredFields{
		Name:   fmt.Sprintf("%s-ibm-export-item", getAccTestPrefix()),
		PlanID: getKentikPlanIDAccTests(),
		IBMProperties: models.IBMPropertiesRequiredFields{
			Bucket: fmt.Sprintf("%s-terraform-ibm-bucket", getAccTestPrefix()),
		},
	})
	ce.Type = models.CloudExportTypeKentikManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = fmt.Sprintf("%s-description-ibm", getAccTestPrefix())
	ce.BGP = &models.BGPProperties{
		ApplyBGP:       pointer.ToBool(true),
		UseBGPDeviceID: fmt.Sprintf("%s-device-id", getAccTestPrefix()),
		DeviceBGPType:  fmt.Sprintf("%s-device-bgp-type", getAccTestPrefix()),
	}
	ce, err = client.CloudExports.Create(ctx, ce)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}
	return ce, nil
}

func createTestAccCloudExportItemAzure() (*models.CloudExport, error) {
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ce := models.NewAzureCloudExport(models.CloudExportAzureRequiredFields{
		Name:   fmt.Sprintf("%s-azure-export-item", getAccTestPrefix()),
		PlanID: getKentikPlanIDAccTests(),
		AzureProperties: models.AzurePropertiesRequiredFields{
			Location:       "centralus",
			ResourceGroup:  fmt.Sprintf("%s-traffic-generator", getAccTestPrefix()),
			StorageAccount: fmt.Sprintf("%s-sa", getAccTestPrefix()),
			SubscriptionID: fmt.Sprintf("%s-sid", getAccTestPrefix()),
		},
	})
	ce.Type = models.CloudExportTypeCustomerManaged
	ce.Enabled = pointer.ToBool(true)
	ce.Description = fmt.Sprintf("%s-description-azure", getAccTestPrefix())
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

func getKentikPlanIDAccTests() string {
	planID, _ := os.LookupEnv("KENTIK_PLAN_ID")
	return planID
}

func newClient() (*kentikapi.Client, error) {
	authEmail, _ := os.LookupEnv("KTAPI_AUTH_EMAIL")
	authToken, _ := os.LookupEnv("KTAPI_AUTH_TOKEN")
	apiURL, _ := os.LookupEnv("KTAPI_URL")
	if authEmail == "" {
		return nil, fmt.Errorf("authEmail variable is empty")
	}
	if authToken == "" {
		return nil, fmt.Errorf("authToken variable is empty")
	}
	client, err := kentikapi.NewClient(kentikapi.Config{
		APIURL:      apiURL,
		AuthEmail:   authEmail,
		AuthToken:   authToken,
		LogPayloads: false,
	})
	if err != nil {
		return nil, fmt.Errorf("newClient: %w", err)
	}
	return client, nil
}
