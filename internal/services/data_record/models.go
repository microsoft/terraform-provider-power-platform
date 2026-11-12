// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package data_record

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
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
	DataRecordClient  client
	EnvironmentClient environment.Client
}

// Values of the data record resource's relationship_removal_behavior attribute.
const (
	// RelationshipRemovalBehaviorDeleteRelatedRecord deletes the related row when it is removed from a
	// collection-valued relationship. This is the default and the provider's original behavior.
	RelationshipRemovalBehaviorDeleteRelatedRecord = "delete_related_record"
	// RelationshipRemovalBehaviorDisassociate only removes the link (a $ref DELETE) and keeps the related row.
	RelationshipRemovalBehaviorDisassociate = "disassociate"
)

type DataRecordResourceModel struct {
	Timeouts                    timeouts.Value `tfsdk:"timeouts"`
	DisableOnDestroy            types.Bool     `tfsdk:"disable_on_destroy"`
	RelationshipRemovalBehavior types.String   `tfsdk:"relationship_removal_behavior"`
	Id                          types.String   `tfsdk:"id"`
	EnvironmentId               types.String   `tfsdk:"environment_id"`
	TableLogicalName            types.String   `tfsdk:"table_logical_name"`
	Columns                     types.Dynamic  `tfsdk:"columns"`
}

// relationshipRemovalBehavior returns the configured relationship removal behavior, or the default when it is not set.
func (m *DataRecordResourceModel) relationshipRemovalBehavior() string {
	if m.RelationshipRemovalBehavior.IsNull() || m.RelationshipRemovalBehavior.IsUnknown() || m.RelationshipRemovalBehavior.ValueString() == "" {
		return RelationshipRemovalBehaviorDeleteRelatedRecord
	}
	return m.RelationshipRemovalBehavior.ValueString()
}
