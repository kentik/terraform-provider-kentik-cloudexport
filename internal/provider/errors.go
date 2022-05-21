package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func detailedDiagError(summary string, err error) diag.Diagnostics {
	diags := diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   err.Error(),
	}
	return diag.Diagnostics{diags}
}
