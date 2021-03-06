//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/kentik/terraform-provider-kentik-cloudexport/internal/provider"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like Delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New, Debug: debugMode}

	plugin.Serve(opts)
}
