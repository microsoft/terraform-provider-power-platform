// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataSource struct {
	helpers.TypeInfo
	TenantClient Client
}

type DataSourceModel struct {
	TenantId                         types.String `tfsdk:"tenant_id"`
	State                            types.String `tfsdk:"state"`
	Location                         types.String `tfsdk:"location"`
	AadCountryGeo                    types.String `tfsdk:"aad_country_geo"`
	DataStorageGeo                   types.String `tfsdk:"data_storage_geo"`
	DefaultEnvironmentGeo            types.String `tfsdk:"default_environment_geo"`
	AadDataBoundary                  types.String `tfsdk:"aad_data_boundary"`
	FedRAMPHighCertificationRequired types.Bool   `tfsdk:"fed_ramp_high_certification_required"`
}
