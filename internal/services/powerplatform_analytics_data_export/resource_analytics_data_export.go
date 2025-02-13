// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &ResourceAnalyticsDataExport{}

func NewBillingPolicyResource() resource.Resource {
	return &ResourceAnalyticsDataExport{
		TypeInfo: helpers.TypeInfo{
			TypeName: "analytics_data_export",
		},
	}
}

// Create implements resource.Resource.
func (r *ResourceAnalyticsDataExport) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (r *ResourceAnalyticsDataExport) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Metadata implements resource.Resource.
func (r *ResourceAnalyticsDataExport) Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse) {
	panic("unimplemented")
}

// Read implements resource.Resource.
func (r *ResourceAnalyticsDataExport) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements resource.Resource.
func (r *ResourceAnalyticsDataExport) Schema(context.Context, resource.SchemaRequest, *resource.SchemaResponse) {
	panic("unimplemented")
}

// Update implements resource.Resource.
func (r *ResourceAnalyticsDataExport) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
