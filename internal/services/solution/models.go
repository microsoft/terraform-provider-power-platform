// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataSource struct {
	helpers.TypeInfo
	SolutionClient Client
}

type ListDataSourceModel struct {
	Timeouts      timeouts.Value    `tfsdk:"timeouts"`
	EnvironmentId types.String      `tfsdk:"environment_id"`
	Solutions     []DataSourceModel `tfsdk:"solutions"`
}

type DataSourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
	Id            types.String `tfsdk:"id"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	InstallTime   types.String `tfsdk:"install_time"`
	Version       types.String `tfsdk:"version"`
	IsManaged     types.Bool   `tfsdk:"is_managed"`
}

func convertFromSolutionDto(solutionDto SolutionDto) DataSourceModel {
	return DataSourceModel{
		EnvironmentId: types.StringValue(solutionDto.EnvironmentId),
		DisplayName:   types.StringValue(solutionDto.DisplayName),
		Name:          types.StringValue(solutionDto.Name),
		CreatedTime:   types.StringValue(solutionDto.CreatedTime),
		Id:            types.StringValue(solutionDto.Id),
		ModifiedTime:  types.StringValue(solutionDto.ModifiedTime),
		InstallTime:   types.StringValue(solutionDto.InstallTime),
		Version:       types.StringValue(solutionDto.Version),
		IsManaged:     types.BoolValue(solutionDto.IsManaged),
	}
}

type Resource struct {
	helpers.TypeInfo
	SolutionClient Client
}

type ResourceModel struct {
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
	Id                   types.String   `tfsdk:"id"`
	SolutionFileChecksum types.String   `tfsdk:"solution_file_checksum"`
	SettingsFileChecksum types.String   `tfsdk:"settings_file_checksum"`
	EnvironmentId        types.String   `tfsdk:"environment_id"`
	SolutionVersion      types.String   `tfsdk:"solution_version"`
	SolutionFile         types.String   `tfsdk:"solution_file"`
	SettingsFile         types.String   `tfsdk:"settings_file"`
	IsManaged            types.Bool     `tfsdk:"is_managed"`
	DisplayName          types.String   `tfsdk:"display_name"`
}
