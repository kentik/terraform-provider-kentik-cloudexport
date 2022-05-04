//go:build tools

package tools

//nolint:typecheck
import (
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs" // documentation generation
	_ "mvdan.cc/gofumpt"                                            // code formatting
)
