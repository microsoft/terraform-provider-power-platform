// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Value    []DataModel    `tfsdk:"locations"`
}

type DataModel struct {
	ID                                     string   `tfsdk:"id"`
	Name                                   string   `tfsdk:"name"`
	DisplayName                            string   `tfsdk:"display_name"`
	Code                                   string   `tfsdk:"code"`
	IsDefault                              bool     `tfsdk:"is_default"`
	IsDisabled                             bool     `tfsdk:"is_disabled"`
	CanProvisionDatabase                   bool     `tfsdk:"can_provision_database"`
	CanProvisionCustomerEngagementDatabase bool     `tfsdk:"can_provision_customer_engagement_database"`
	AzureRegions                           []string `tfsdk:"azure_regions"`
}

type macroRegionsDataSourceModel struct {
	Timeouts     timeouts.Value         `tfsdk:"timeouts"`
	MacroRegions []macroRegionDataModel `tfsdk:"macro_regions"`
}

type macroRegionDataModel struct {
	MacroRegionId     types.String `tfsdk:"macro_region_id"`
	DisplayName       types.String `tfsdk:"display_name"`
	DataResidencyNote types.String `tfsdk:"data_residency_note"`
}
