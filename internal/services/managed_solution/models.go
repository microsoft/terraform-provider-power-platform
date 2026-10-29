// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type Resource struct {
	helpers.TypeInfo
	Client Client
}

type ResourceModel struct {
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
	Id                            types.String   `tfsdk:"id"`
	EnvironmentId                 types.String   `tfsdk:"environment_id"`
	UniqueName                    types.String   `tfsdk:"unique_name"`
	Version                       types.String   `tfsdk:"version"`
	ConnectionReferences          types.Map      `tfsdk:"connection_references"`
	SkipProductUpdateDependencies types.Bool     `tfsdk:"skip_product_update_dependencies"`
	PublishAllCustomizations      types.Bool     `tfsdk:"publish_all_customizations"`
	Source                        *SourceModel   `tfsdk:"source"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	SolutionId                    types.String   `tfsdk:"solution_id"`
}

type SourceModel struct {
	Path types.String `tfsdk:"path"`
	URL  types.String `tfsdk:"url"`
}
