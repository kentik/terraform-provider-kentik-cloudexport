package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/kentikapi"
)

func dataSourceCloudExportList() *schema.Resource {
	return &schema.Resource{
		Description: "Data source representing list of cloud exports",
		ReadContext: dataSourceCloudExportListRead,
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: makeCloudExportSchema(readList),
				},
			},
		},
	}
}

func dataSourceCloudExportListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "List cloud export Kentik API request")
	listResp, httpResp, err := m.(*kentikapi.Client).CloudExportAdminServiceAPI.
		ExportList(ctx).
		Execute()
	tflog.Debug(ctx, "List cloud export Kentik API response", map[string]interface{}{"response": listResp})
	if err != nil {
		return detailedDiagError("Failed to read cloud export list", err, httpResp)
	}

	if listResp.Exports != nil {
		numExports := len(*listResp.Exports)
		exports := make([]interface{}, numExports)
		for i, e := range *listResp.Exports {
			ee := e // avoid implicit memory aliasing in for loop (G601)
			exports[i] = cloudExportToMap(&ee)
		}

		if err = d.Set("items", exports); err != nil {
			return diag.FromErr(err)
		}
	}

	// use UNIX time as ID to force list update every time Terraform asks for the list
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
