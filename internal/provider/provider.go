package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/apiv6/kentikapi"
)

const (
	apiURLKey = "apiurl"
	emailKey  = "email"
	tokenKey  = "token"
)

func init() {
	// Set descriptions to support Markdown syntax, this will be used in document generation and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			apiURLKey: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KTAPI_URL", nil),
				Description: "CloudExport API server URL (optional). Can also be specified with KTAPI_URL environment variable.",
			},
			emailKey: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KTAPI_AUTH_EMAIL", nil),
				Description: "Authorization email (required). Can also be specified with KTAPI_AUTH_EMAIL environment variable.",
			},
			tokenKey: {
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KTAPI_AUTH_TOKEN", nil),
				Description: "Authorization token (required). Can also be specified with KTAPI_AUTH_TOKEN environment variable.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kentik-cloudexport_list": dataSourceCloudExportList(),
			"kentik-cloudexport_item": dataSourceCloudExportItem(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"kentik-cloudexport_item": resourceCloudExport(),
		},
		ConfigureContextFunc: configure,
	}
}

func configure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// TODO: comment why required here
	authEmail, ok := d.GetOk(emailKey)
	if !ok {
		return nil, diag.Errorf("Missing required %v argument", emailKey)
	}

	authToken := d.Get(tokenKey)
	if !ok {
		return nil, diag.Errorf("Missing required %v argument", tokenKey)
	}

	return kentikapi.NewClient(kentikapi.Config{
		CloudExportAPIURL: getURL(d),
		AuthEmail:         authEmail.(string),
		AuthToken:         authToken.(string),
	}), nil
}

func getURL(d *schema.ResourceData) string {
	var url string
	apiURL, ok := d.GetOk(apiURLKey)
	if ok {
		url = apiURL.(string)
	}
	return url
}
