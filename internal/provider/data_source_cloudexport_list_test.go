package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/kentik/community_sdk_golang/kentikapi/models"
	"github.com/stretchr/testify/assert"
)

const (
	exportsDS = "data.kentik-cloudexport_list.exports"
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

func TestAccDataSourceCloudExportList(t *testing.T) {
	exports, err := createTestAccCloudExportList()
	assert.NoError(t, err)
	defer func() {
		for _, ce := range exports {
			assert.NoError(t, deleteTestAccCloudExportItem(ce))
		}
	}()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: makeTestAccCloudExportDataSourceList(),
				Check: resource.ComposeTestCheckFunc(
					// more properties are verified in TestAccDataSourceCloudExportItem* tests
					resource.TestCheckResourceAttrSet(exportsDS, "items.0.name"),
					resource.TestCheckResourceAttrSet(exportsDS, "items.1.name"),
				),
			},
		},
	})
}

func makeTestAccCloudExportDataSourceList() string {
	return `
		data "kentik-cloudexport_list" "exports" {}
	`
}

func createTestAccCloudExportList() ([]*models.CloudExport, error) {
	_, accTest := os.LookupEnv("TF_ACC")
	if !accTest {
		return nil, nil
	}
	ctx := context.Background()
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	ceAWS := models.NewAWSCloudExport(models.CloudExportAWSRequiredFields{
		Name:   "acc_test_terraform_aws_export_list",
		PlanID: "11467",
		AWSProperties: models.AWSPropertiesRequiredFields{
			Bucket: "terraform-aws-bucket",
		},
	})
	ceAWS.Type = models.CloudExportTypeKentikManaged
	ceAWS.Description = "terraform aws cloud export"
	ceAWS.GetAWSProperties().IAMRoleARN = "dummy-iam-role-arn"
	ceAWS.GetAWSProperties().Region = "us-east-2"
	ceAWS.GetAWSProperties().DeleteAfterRead = pointer.ToBool(true)
	ceAWS.GetAWSProperties().MultipleBuckets = pointer.ToBool(true)
	ceAWS, err = client.CloudExports.Create(ctx, ceAWS)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}

	ceGCE := models.NewGCECloudExport(models.CloudExportGCERequiredFields{
		Name:   "acc_test_terraform_gce_export_list",
		PlanID: "21600",
		GCEProperties: models.GCEPropertiesRequiredFields{
			Project:      "project gce",
			Subscription: "subscription gce",
		},
	})
	ceGCE.Type = models.CloudExportTypeCustomerManaged
	ceGCE.Description = "terraform gce cloud export"
	ceGCE, err = client.CloudExports.Create(ctx, ceGCE)
	if err != nil {
		return nil, fmt.Errorf("client.CloudExports.Create: %w", err)
	}

	return []*models.CloudExport{ceAWS, ceGCE}, nil
}
