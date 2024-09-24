// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type SharesDataSource struct {
	helpers.TypeInfo
	ConnectionsClient client
}

type SharesListDataSourceModel struct {
	Timeouts      timeouts.Value          `tfsdk:"timeouts"`
	EnvironmentId types.String            `tfsdk:"environment_id"`
	ConnectorName types.String            `tfsdk:"connector_name"`
	ConnectionId  types.String            `tfsdk:"connection_id"`
	Shares        []SharesDataSourceModel `tfsdk:"shares"`
}

type SharesDataSourceModel struct {
	Id        types.String                   `tfsdk:"id"`
	RoleName  types.String                   `tfsdk:"role_name"`
	Principal SharesPrincipalDataSourceModel `tfsdk:"principal"`
}

type SharesPrincipalDataSourceModel struct {
	EntraId     types.String `tfsdk:"entra_object_id"`
	DisplayName types.String `tfsdk:"display_name"`
}

type ConnectionsDataSource struct {
	helpers.TypeInfo
	ConnectionsClient client
}

type ConnectionsListDataSourceModel struct {
	Timeouts      timeouts.Value               `tfsdk:"timeouts"`
	EnvironmentId types.String                 `tfsdk:"environment_id"`
	Connections   []ConnectionsDataSourceModel `tfsdk:"connections"`
}

type ConnectionsDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	DisplayName             types.String `tfsdk:"display_name"`
	Status                  []string     `tfsdk:"status"`
	ConnectionParameters    types.String `tfsdk:"connection_parameters"`
	ConnectionParametersSet types.String `tfsdk:"connection_parameters_set"`
}

type ShareResource struct {
	helpers.TypeInfo
	ConnectionsClient client
}

type ShareResourceModel struct {
	Timeouts      timeouts.Value              `tfsdk:"timeouts"`
	Id            types.String                `tfsdk:"id"`
	EnvironmentId types.String                `tfsdk:"environment_id"`
	ConnectorName types.String                `tfsdk:"connector_name"`
	ConnectionId  types.String                `tfsdk:"connection_id"`
	RoleName      types.String                `tfsdk:"role_name"`
	Principal     SharePrincipalResourceModel `tfsdk:"principal"`
}

type SharePrincipalResourceModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}

type Resource struct {
	helpers.TypeInfo
	ConnectionsClient client
}

type ResourceModel struct {
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
	Id                      types.String   `tfsdk:"id"`
	Name                    types.String   `tfsdk:"name"`
	EnvironmentId           types.String   `tfsdk:"environment_id"`
	DisplayName             types.String   `tfsdk:"display_name"`
	Status                  types.Set      `tfsdk:"status"`
	ConnectionParameters    types.String   `tfsdk:"connection_parameters"`
	ConnectionParametersSet types.String   `tfsdk:"connection_parameters_set"`
}
