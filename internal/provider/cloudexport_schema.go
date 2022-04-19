package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/kentikapi/models"
)

// CloudExportSchema reflects CloudExport type and defines a Cloud Export item used in terraform .tf files
// Note: currently, nesting an object is only possible by using single-item List element (Terraform limitation).
func makeCloudExportSchema(mode schemaMode) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    mode == create || mode == readList, // provided by server on creating/listing items
			Required:    mode == readSingle,                 // provided by user in order to read single item
			Description: "The internal cloud export identifier. This is Read-only and assigned by Kentik",
		},
		"type": {
			Type:     schema.TypeString,
			Computed: mode == readSingle || mode == readList, // provided by server on read
			Required: mode == create,                         // provided by user on create
			Description: "CLOUD_EXPORT_TYPE_UNSPECIFIED: Invalid or incomplete exports. " +
				"CLOUD_EXPORT_TYPE_KENTIK_MANAGED: Cloud exports that are managed by Kentik. " +
				"CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED: Exports that are managed by Kentik customers " +
				"(eg. by running an agent)",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Computed:    mode == readSingle || mode == readList, // provided by server on read
			Required:    mode == create,                         // provided by user on create
			Description: "Whether this task is enabled and intended to run, or disabled",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    mode == readSingle || mode == readList, // provided by server on read
			Required:    mode == create,                         // provided by user on create
			Description: "A short name for this export",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    mode == readSingle || mode == readList, // provided by server on read
			Optional:    mode == create,                         // optionally provided by user on create
			Description: "An optional, longer description",
		},
		"plan_id": {
			Type:        schema.TypeString,
			Computed:    mode == readSingle || mode == readList, // provided by server on read
			Required:    mode == create,                         // provided by user on create
			Description: "The identifier of the Kentik plan associated with this task",
		},
		"cloud_provider": {
			Type:        schema.TypeString,
			Computed:    mode == readSingle || mode == readList, // provided by server on read
			Required:    mode == create,                         // provided by user on create
			Description: "The cloud provider targeted by this export (aws, azure, gce, ibm)",
		},
		"aws":            makeAWSSchema(mode),
		"azure":          makeAzureSchema(mode),
		"bgp":            makeBGPSchema(mode),
		"gce":            makeGCESchema(mode),
		"ibm":            makeIBMSchema(mode),
		"current_status": makeCurrentStatusSchema(),
	}
}

func makeAWSSchema(mode schemaMode) *schema.Schema {
	return &schema.Schema{
		// nested object
		Type:        schema.TypeList,
		Computed:    mode == readSingle || mode == readList, // provided by server on read
		Optional:    mode == create,                         // optionally provided by user on create
		Description: "Properties specific to Amazon Web Services \"vpc flow logs\" exports",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bucket": {
					Type:        schema.TypeString,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "Source S3 bucket to fetch vpc flow logs from",
				},
				"iam_role_arn": {
					Type:        schema.TypeString,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "ARN for the IAM role to assume when fetching data or making AWS calls for this export",
				},
				"region": {
					Type:        schema.TypeString,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "AWS region where this bucket resides",
				},
				"delete_after_read": {
					Type:        schema.TypeBool,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "If true, attempt to delete vpc flow log chunks from S3 after they've been read",
				},
				"multiple_buckets": {
					Type:     schema.TypeBool,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
			},
		},
	}
}

func makeAzureSchema(mode schemaMode) *schema.Schema {
	return &schema.Schema{
		// nested object
		Type:        schema.TypeList,
		Computed:    mode == readSingle || mode == readList, // provided by server on read
		Optional:    mode == create,                         // optionally provided by user on create
		Description: "Properties specific to Azure exports",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"location": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
				"resource_group": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
				"storage_account": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
				"subscription_id": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
				"security_principal_enabled": {
					Type:     schema.TypeBool,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
			},
		},
	}
}

func makeBGPSchema(mode schemaMode) *schema.Schema {
	return &schema.Schema{
		// nested object
		Type:        schema.TypeList,
		Computed:    mode == readSingle || mode == readList, // provided by server on read
		Optional:    mode == create,                         // optionally provided by user on create
		Description: "Optional BGP related settings",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"apply_bgp": {
					Type:        schema.TypeBool,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "If true, apply BGP data discovered via another device to the flow from this export",
				},
				"use_bgp_device_id": {
					Type:        schema.TypeString,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "Which other device to get BGP data from",
				},
				"device_bgp_type": {
					Type:        schema.TypeString,
					Computed:    mode == readSingle || mode == readList, // provided by server on read
					Required:    mode == create,                         // provided by user on create
					Description: "device, other_device, none",
				},
			},
		},
	}
}

func makeGCESchema(mode schemaMode) *schema.Schema {
	return &schema.Schema{
		// nested object
		Type:        schema.TypeList,
		Computed:    mode == readSingle || mode == readList, // provided by server on read
		Optional:    mode == create,                         // optionally provided by user on create
		Description: "Properties specific to Google Cloud export",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"project": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
				"subscription": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
			},
		},
	}
}

func makeIBMSchema(mode schemaMode) *schema.Schema {
	return &schema.Schema{
		// nested object
		Type:        schema.TypeList,
		Computed:    mode == readSingle || mode == readList, // provided by server on read
		Optional:    mode == create,                         // optionally provided by user on create
		Description: "Properties specific to IBM Cloud exports",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bucket": {
					Type:     schema.TypeString,
					Computed: mode == readSingle || mode == readList, // provided by server on read
					Required: mode == create,                         // provided by user on create
				},
			},
		},
	}
}

func makeCurrentStatusSchema() *schema.Schema {
	return &schema.Schema{
		// nested object
		Type:        schema.TypeList,
		Computed:    true, // always provided by server
		Description: "Export task status",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"status": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "OK, ERROR or other short and descriptive status",
				},
				"error_message": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "If not empty, the current error",
				},
				"flow_found": {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "If true, we found flow logs",
				},
				"api_access": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"storage_account_access": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

// cloudExportToMap is used for API get operation to fill terraform resource from cloudexport item.
func cloudExportToMap(e *models.CloudExport) map[string]interface{} {
	o := make(map[string]interface{})
	if e == nil {
		return o
	}

	o["id"] = e.ID
	o["type"] = e.Type
	o["enabled"] = e.Enabled
	o["name"] = e.Name
	o["description"] = e.Description
	o["plan_id"] = e.PlanID
	o["cloud_provider"] = e.CloudProvider

	if e.GetAWSProperties() != nil {
		aws := make(map[string]interface{})
		aws["bucket"] = e.GetAWSProperties().Bucket
		aws["iam_role_arn"] = e.GetAWSProperties().IAMRoleARN
		aws["region"] = e.GetAWSProperties().Region
		aws["delete_after_read"] = e.GetAWSProperties().DeleteAfterRead
		aws["multiple_buckets"] = e.GetAWSProperties().MultipleBuckets
		o["aws"] = []interface{}{aws}
	}

	if e.GetAzureProperties() != nil {
		azure := make(map[string]interface{})
		azure["location"] = e.GetAzureProperties().Location
		azure["resource_group"] = e.GetAzureProperties().ResourceGroup
		azure["storage_account"] = e.GetAzureProperties().StorageAccount
		azure["subscription_id"] = e.GetAzureProperties().SubscriptionID
		azure["security_principal_enabled"] = e.GetAzureProperties().SecurityPrincipalEnabled
		o["azure"] = []interface{}{azure}
	}

	if e.BGP != nil {
		bgp := make(map[string]interface{})
		bgp["apply_bgp"] = e.BGP.ApplyBGP
		bgp["use_bgp_device_id"] = e.BGP.UseBGPDeviceID
		bgp["device_bgp_type"] = e.BGP.DeviceBGPType
		o["bgp"] = []interface{}{bgp}
	}

	if e.GetGCEProperties() != nil {
		gce := make(map[string]interface{})
		gce["project"] = e.GetGCEProperties().Project
		gce["subscription"] = e.GetGCEProperties().Subscription
		o["gce"] = []interface{}{gce}
	}

	if e.GetIBMProperties() != nil {
		ibm := make(map[string]interface{})
		ibm["bucket"] = e.GetIBMProperties().Bucket
		o["ibm"] = []interface{}{ibm}
	}

	if e.CurrentStatus != nil {
		cs := make(map[string]interface{})
		cs["status"] = e.CurrentStatus.Status
		cs["error_message"] = e.CurrentStatus.ErrorMessage
		cs["flow_found"] = e.CurrentStatus.FlowFound
		cs["api_access"] = e.CurrentStatus.APIAccess
		cs["storage_account_access"] = e.CurrentStatus.StorageAccountAccess
		o["current_status"] = []interface{}{cs}
	}

	return o
}

// resourceDataToCloudExport is used for API create/update operations to fill cloudexport item from terraform resource.
//nolint: gocyclo
func resourceDataToCloudExport(d *schema.ResourceData) (*models.CloudExport, error) {
	// Note: only set the user-writable attributes, read-only attributes that are only generated on server side:
	// CurrentStatus, are left with nil values and so are not serialized and not sent to API server
	export := models.CloudExport{}

	ceID := d.Get("id")
	if v, ok := ceID.(string); ok {
		export.ID = v
	}

	// required
	t := d.Get("type")
	export.Type = models.CloudExportType(t.(string))

	// required
	enabled := d.Get("enabled")
	if v, ok := enabled.(bool); ok {
		export.Enabled = &v
	}

	// required
	name := d.Get("name")
	if v, ok := name.(string); ok {
		export.Name = v
	}

	// optional
	if description, ok := d.GetOk("description"); ok {
		if v, ok := description.(string); ok {
			export.Description = v
		}
	}

	// required
	planID := d.Get("plan_id")
	if v, ok := planID.(string); ok {
		export.PlanID = v
	}

	// required
	cloudProvider := d.Get("cloud_provider").(string) //nolint: errcheck, forcetypeassert // type enforced by schema
	export.CloudProvider = models.CloudProvider(cloudProvider)

	// validation: for any given cloud_provider, there should also be an object of the same name,
	// containing configuration details, e.g. for cloud_provider="ibm", ibm{...} object should be defined
	providerObj, ok := d.GetOk(cloudProvider)
	if !ok {
		return nil, fmt.Errorf("for cloud_provider=%[1]s, there should also be %[1]s{...} attribute provided", cloudProvider)
	}
	providerDef := providerObj.([]interface{})[0]       // extract nested object under index 0. Terraform clumsiness
	providerMap := providerDef.(map[string]interface{}) //nolint: errcheck, forcetypeassert // type enforced by schema
	switch cloudProvider {
	case "aws":
		{
			aws := models.AWSProperties{
				Bucket:     providerMap["bucket"].(string),
				IAMRoleARN: providerMap["iam_role_arn"].(string),
				Region:     providerMap["region"].(string),
			}
			deleteAfterRead := providerMap["delete_after_read"]
			if v, ok := deleteAfterRead.(bool); ok {
				aws.DeleteAfterRead = &v
			}
			multipleBuckets := providerMap["multiple_buckets"]
			if v, ok := multipleBuckets.(bool); ok {
				aws.MultipleBuckets = &v
			}
			export.Properties = &aws
		}
	case "azure":
		{
			azure := models.AzureProperties{
				Location:       providerMap["location"].(string),
				ResourceGroup:  providerMap["resource_group"].(string),
				StorageAccount: providerMap["storage_account"].(string),
				SubscriptionID: providerMap["subscription_id"].(string),
			}
			securityPrincipalEnabled := providerMap["security_principal_enabled"]
			if v, ok := securityPrincipalEnabled.(bool); ok {
				azure.SecurityPrincipalEnabled = &v
			}
			export.Properties = &azure
		}
	case "gce":
		{
			gce := models.GCEProperties{
				Project:      providerMap["project"].(string),
				Subscription: providerMap["subscription"].(string),
			}
			export.Properties = &gce
		}
	case "ibm":
		{
			ibm := models.IBMProperties{
				Bucket: providerMap["bucket"].(string),
			}
			export.Properties = &ibm
		}
	default:
		return nil, fmt.Errorf("cloud_provider should be one of [aws, azure, ibm, gce], got: %q", cloudProvider)
	}
	if cloudProvider == "bgp" {
		bgp := models.BGPProperties{
			UseBGPDeviceID: providerMap["use_bgp_device_id"].(string),
			DeviceBGPType:  providerMap["device_bgp_type"].(string),
		}
		applyBGP := providerMap["apply_bgp"]
		if v, ok := applyBGP.(bool); ok {
			bgp.ApplyBGP = &v
		}
		export.BGP = &bgp
	}
	return &export, nil
}
