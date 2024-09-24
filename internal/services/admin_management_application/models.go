// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package admin_management_application

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type AdminManagementApplicationResource struct {
	helpers.TypeInfo
	AdminManagementApplicationClient client
}

type AdminManagementApplicationResourceModel struct {
	Timeouts timeouts.Value        `tfsdk:"timeouts"`
	Id       customtypes.UUIDValue `tfsdk:"id"`
}
