// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connectors

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataSource struct {
	helpers.TypeInfo
	ConnectorsClient client
}

type ListDataSourceModel struct {
	Timeouts      timeouts.Value    `tfsdk:"timeouts"`
	EnvironmentId types.String      `tfsdk:"environment_id"`
	Connectors    []DataSourceModel `tfsdk:"connectors"`
}

type DataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Tier        types.String `tfsdk:"tier"`
	Publisher   types.String `tfsdk:"publisher"`
	Unblockable types.Bool   `tfsdk:"unblockable"`
}

func convertFromConnectorDto(connectorDto connectorDto) DataSourceModel {
	return DataSourceModel{
		Id:          types.StringValue(connectorDto.Id),
		Name:        types.StringValue(connectorDto.Name),
		Type:        types.StringValue(connectorDto.Type),
		Description: types.StringValue(connectorDto.Properties.Description),
		DisplayName: types.StringValue(connectorDto.Properties.DisplayName),
		Tier:        types.StringValue(connectorDto.Properties.Tier),
		Publisher:   types.StringValue(connectorDto.Properties.Publisher),
		Unblockable: types.BoolValue(connectorDto.Properties.Unblockable),
	}
}
