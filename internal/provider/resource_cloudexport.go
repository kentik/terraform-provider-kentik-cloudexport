package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/kentikapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func resourceCloudExport() *schema.Resource {
	return &schema.Resource{
		Description:   "Resource representing cloud export item",
		CreateContext: resourceCloudExportCreate,
		ReadContext:   resourceCloudExportRead,
		UpdateContext: resourceCloudExportUpdate,
		DeleteContext: resourceCloudExportDelete,
		Schema:        makeCloudExportSchema(create),
	}
}

func resourceCloudExportCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	export, err := resourceDataToCloudExport(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Create cloud export Kentik API request", map[string]interface{}{"request": export})

	export, err = m.(*kentikapi.Client).CloudExports.Create(ctx, export)
	tflog.Debug(ctx, "Create cloud export Kentik API response", map[string]interface{}{"response": export})
	if err != nil {
		return detailedDiagError("Failed to create cloud export", err)
	}

	err = d.Set("id", export.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(export.ID) // create the resource in TF state

	// read back the just-created resource to handle the case when server applies modifications to provided data
	return resourceCloudExportRead(ctx, d, m)
}

func resourceCloudExportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Get cloud export Kentik API request", map[string]interface{}{"ID": d.Get("id").(string)})
	export, err := m.(*kentikapi.Client).CloudExports.Get(ctx, d.Get("id").(string))

	tflog.Debug(ctx, "Get cloud export Kentik API response", map[string]interface{}{"response": export})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				d.SetId("") // delete the resource in TF state
				return nil
			}
		}
		return detailedDiagError("Failed to read cloud export", err)
	}
	mapExport := cloudExportToMap(export)
	for k, v := range mapExport {
		if err = d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceCloudExportUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// check if any attribute has changed
	if d.HasChange("") {
		export, err := resourceDataToCloudExport(d)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, "Update cloud export Kentik API request", map[string]interface{}{"request": export})
		resp, err := m.(*kentikapi.Client).CloudExports.Update(ctx, export)
		tflog.Debug(ctx, "Update cloud export Kentik API response", map[string]interface{}{"response": resp})
		if err != nil {
			return detailedDiagError("Failed to update cloud export", err)
		}
	}

	// read back the just-updated resource to handle the case when server applies modifications to provided data
	return resourceCloudExportRead(ctx, d, m)
}

func resourceCloudExportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Delete cloud export Kentik API request", map[string]interface{}{"ID": d.Get("id").(string)})
	err := m.(*kentikapi.Client).CloudExports.Delete(ctx, d.Get("id").(string))
	if err != nil {
		return detailedDiagError("Failed to delete cloud export", err)
	}
	tflog.Debug(ctx, "Deleted cloud export in Kentik", map[string]interface{}{"ID": d.Get("id").(string)})
	return nil
}
