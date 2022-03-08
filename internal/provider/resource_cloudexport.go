package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/kentikapi"
	"github.com/kentik/community_sdk_golang/kentikapi/cloudexport"
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

	req := *cloudexport.NewV202101beta1CreateCloudExportRequest()
	req.Export = export
	tflog.Debug(ctx, "Create cloud export Kentik API request", "request", req)
	resp, httpResp, err := m.(*kentikapi.Client).CloudExportAdminServiceAPI.
		ExportCreate(ctx).
		Body(req).
		Execute()
	tflog.Debug(ctx, "Create cloud export Kentik API response", "response", resp)
	if err != nil {
		return detailedDiagError("Failed to create cloud export", err, httpResp)
	}

	err = d.Set("id", resp.Export.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Export.GetId()) // create the resource in TF state

	// read back the just-created resource to handle the case when server applies modifications to provided data
	return resourceCloudExportRead(ctx, d, m)
}

func resourceCloudExportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Get cloud export Kentik API request", "ID", d.Get("id").(string))
	resp, httpResp, err := m.(*kentikapi.Client).CloudExportAdminServiceAPI.
		ExportGet(ctx, d.Get("id").(string)).
		Execute()
	tflog.Debug(ctx, "Get cloud export Kentik API response", "response", resp)
	if err != nil {
		if httpResp.StatusCode == http.StatusNotFound {
			d.SetId("") // delete the resource in TF state
			return nil
		}
		return detailedDiagError("Failed to read cloud export", err, httpResp)
	}
	mapExport := cloudExportToMap(resp.Export)
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

		req := *cloudexport.NewV202101beta1UpdateCloudExportRequest()
		req.Export = export
		tflog.Debug(ctx, "Update cloud export Kentik API request", "request", req)
		_, httpResp, err := m.(*kentikapi.Client).CloudExportAdminServiceAPI.
			ExportUpdate(ctx, d.Get("id").(string)).
			Body(req).
			Execute()
		tflog.Debug(ctx, "Update cloud export Kentik API response", "response", httpResp.Body)
		if err != nil {
			return detailedDiagError("Failed to update cloud export", err, httpResp)
		}
	}

	// read back the just-updated resource to handle the case when server applies modifications to provided data
	return resourceCloudExportRead(ctx, d, m)
}

func resourceCloudExportDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Delete cloud export Kentik API request", "ID", d.Get("id").(string))
	_, httpResp, err := m.(*kentikapi.Client).CloudExportAdminServiceAPI.
		ExportDelete(ctx, d.Get("id").(string)).
		Execute()
	tflog.Debug(ctx, "Delete cloud export Kentik API response", "response", httpResp.Body)
	if err != nil {
		return detailedDiagError("Failed to delete cloud export", err, httpResp)
	}

	return nil
}
