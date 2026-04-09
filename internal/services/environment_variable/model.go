// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environmentvariable

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type Resource struct {
	helpers.TypeInfo
	Client client
}

type ResourceModel struct {
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
	Id                            types.String   `tfsdk:"id"`
	EnvironmentId                 types.String   `tfsdk:"environment_id"`
	SchemaName                    types.String   `tfsdk:"schema_name"`
	Value                         types.String   `tfsdk:"value"`
	EnvironmentVariableDefinition types.String   `tfsdk:"environment_variable_definition_id"`
	EnvironmentVariableValue      types.String   `tfsdk:"environment_variable_value_id"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	Description                   types.String   `tfsdk:"description"`
	Type                          types.String   `tfsdk:"type"`
	ValueSchema                   types.String   `tfsdk:"value_schema"`
	SecretStore                   types.String   `tfsdk:"secret_store"`
}
