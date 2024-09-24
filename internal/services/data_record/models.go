// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package data_record

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type ExpandModel struct {
	NavigationProperty types.String  `tfsdk:"navigation_property"`
	Select             []string      `tfsdk:"select"`
	Filter             types.String  `tfsdk:"filter"`
	OrderBy            types.String  `tfsdk:"order_by"`
	Top                types.Int64   `tfsdk:"top"`
	Expand             []ExpandModel `tfsdk:"expand"`
}

type DataRecordListDataSourceModel struct {
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
	EnvironmentId               types.String   `tfsdk:"environment_id"`
	EntityCollection            types.String   `tfsdk:"entity_collection"`
	Select                      []string       `tfsdk:"select"`
	Filter                      types.String   `tfsdk:"filter"`
	Apply                       types.String   `tfsdk:"apply"`
	OrderBy                     types.String   `tfsdk:"order_by"`
	Top                         types.Int64    `tfsdk:"top"`
	ReturnTotalRowsCount        types.Bool     `tfsdk:"return_total_rows_count"`
	TotalRowsCount              types.Int64    `tfsdk:"total_rows_count"`
	TotalRowsCountLimitExceeded types.Bool     `tfsdk:"total_rows_count_limit_exceeded"`
	SavedQuery                  types.String   `tfsdk:"saved_query"`
	UserQuery                   types.String   `tfsdk:"user_query"`
	Expand                      []ExpandModel  `tfsdk:"expand"`
	Rows                        types.Dynamic  `tfsdk:"rows"`
}

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataRecordResourceModel struct {
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id               types.String   `tfsdk:"id"`
	EnvironmentId    types.String   `tfsdk:"environment_id"`
	TableLogicalName types.String   `tfsdk:"table_logical_name"`
	Columns          types.Dynamic  `tfsdk:"columns"`
}
