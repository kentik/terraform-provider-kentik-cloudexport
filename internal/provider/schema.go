package provider

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// schemaMode determines if we want a schema for:
// - reading single item - we need to provide "id" of the item to read, everything else is provided by the server,
// - reading list of items - we don't need to provide a thing, everything is provided by the server,
// - creating new item - we need to provide a bunch of obligatory attributes, the rest is provided by the server.
type schemaMode int

const (
	readSingle schemaMode = iota
	readList
	create
)

func skipOnReadDiagFunc(mode schemaMode, diagFunc schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	if mode == readSingle || mode == readList {
		return nil
	}
	return diagFunc
}

func skipOnReadOneOf(mode schemaMode, providers []string) []string {
	if mode == readSingle || mode == readList {
		return nil
	}
	return providers
}
