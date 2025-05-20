// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package environment_application_admin

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type EnvironmentApplicationAdminResource struct {
	helpers.TypeInfo
	EnvironmentApplicationAdminClient client
}

type EnvironmentApplicationAdminResourceModel struct {
	Timeouts      timeouts.Value        `tfsdk:"timeouts"`
	EnvironmentId customtypes.UUIDValue `tfsdk:"environment_id"`
	ApplicationId customtypes.UUIDValue `tfsdk:"application_id"`
	Id            types.String          `tfsdk:"id"`
}
