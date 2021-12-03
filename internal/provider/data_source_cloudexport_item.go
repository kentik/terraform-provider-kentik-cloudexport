package provider

import (
	"context"

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
	getResp, httpResp, err := m.(*kentikapi.Client).CloudExportAdminServiceAPI.
		ExportGet(ctx, d.Get("id").(string)).
		Execute()
	if err != nil {
		return detailedDiagError("Failed to read cloud export item", err, httpResp)
	}

	mapExport := cloudExportToMap(getResp.Export)
	for k, v := range mapExport {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(*getResp.Export.Id)

	return nil
}
