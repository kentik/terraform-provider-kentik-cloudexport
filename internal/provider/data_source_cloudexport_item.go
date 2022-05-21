package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/kentikapi"
)

func dataSourceCloudExportItem() *schema.Resource {
	return &schema.Resource{
		Description: "Data source representing single cloud export item",
		ReadContext: dataSourceCloudExportItemRead,
		Schema:      makeCloudExportSchema(readSingle),
	}
}

func dataSourceCloudExportItemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Get cloud export Kentik API request", map[string]interface{}{"ID": d.Get("id").(string)})
	export, err := m.(*kentikapi.Client).CloudExports.Get(ctx, d.Get("id").(string))
	tflog.Debug(ctx, "Get cloud export Kentik API response", map[string]interface{}{"response": export})
	if err != nil {
		return detailedDiagError("Failed to read cloud export item", err)
	}

	mapExport := cloudExportToMap(export)
	for k, v := range mapExport {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(export.ID)

	return nil
}
