// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type EnvironmentTemplatesDataSource struct {
	helpers.TypeInfo
	EnvironmentTemplatesClient client
}

type EnvironmentTemplatesDataSourceModel struct {
	Timeouts  timeouts.Value                  `tfsdk:"timeouts"`
	Location  types.String                    `tfsdk:"location"`
	Templates []EnvironmentTemplatesDataModel `tfsdk:"environment_templates"`
}

type EnvironmentTemplatesDataModel struct {
	Category                     string `tfsdk:"category"`
	ID                           string `tfsdk:"id"`
	Name                         string `tfsdk:"name"`
	DisplayName                  string `tfsdk:"display_name"`
	Location                     string `tfsdk:"location"`
	IsDisabled                   bool   `tfsdk:"is_disabled"`
	DisabledReasonCode           string `tfsdk:"disabled_reason_code"`
	DisabledReasonMessage        string `tfsdk:"disabled_reason_message"`
	IsCustomerEngagement         bool   `tfsdk:"is_customer_engagement"`
	IsSupportedForResetOperation bool   `tfsdk:"is_supported_for_reset_operation"`
}
