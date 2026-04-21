// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type ApplicationUserResource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type ApplicationUserResourceModel struct {
	Timeouts       timeouts.Value `tfsdk:"timeouts"`
	Id             types.String   `tfsdk:"id"`
	EnvironmentId  types.String   `tfsdk:"environment_id"`
	ApplicationId  types.String   `tfsdk:"application_id"`
	SystemUserId   types.String   `tfsdk:"system_user_id"`
	BusinessUnitId types.String   `tfsdk:"business_unit_id"`
}
