// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataSource struct {
	helpers.TypeInfo
	CapacityClient client
}

type DataSourceModel struct {
	TenantId         types.String                    `tfsdk:"tenant_id"`
	LicenseModelType types.String                    `tfsdk:"license_model_type"`
	TenantCapacities []TenantCapacityDataSourceModel `tfsdk:"tenant_capacities"`
}

type TenantCapacityDataSourceModel struct {
	CapacityType  types.String               `tfsdk:"capacity_type"`
	CapacityUnits types.String               `tfsdk:"capacity_units"`
	TotalCapacity types.Float32              `tfsdk:"total_capacity"`
	MaxCapacity   types.Float32              `tfsdk:"max_capacity"`
	Consumption   ConsumptionDataSourceModel `tfsdk:"consumption"`
	Status        types.String               `tfsdk:"status"`
}

type ConsumptionDataSourceModel struct {
	Actual          types.Float32 `tfsdk:"actual"`
	Rated           types.Float32 `tfsdk:"rated"`
	ActualUpdatedOn types.String  `tfsdk:"actual_updated_on"`
	RatedUpdatedOn  types.String  `tfsdk:"rated_updated_on"`
}
