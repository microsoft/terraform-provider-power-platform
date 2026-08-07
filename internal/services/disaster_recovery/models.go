// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package disaster_recovery

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

type DisasterRecoveryResource struct {
	helpers.TypeInfo
	client client
}

type DisasterRecoveryResourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	Enabled       types.Bool     `tfsdk:"enabled"`
}

func convertDtoToModel(environmentId string, env *environment.EnvironmentDto) DisasterRecoveryResourceModel {
	model := DisasterRecoveryResourceModel{
		Id:            types.StringValue(environmentId),
		EnvironmentId: types.StringValue(environmentId),
		Enabled:       types.BoolValue(false),
	}

	if env.Properties != nil && env.Properties.States != nil && env.Properties.States.DisasterRecovery != nil {
		model.Enabled = types.BoolValue(env.Properties.States.DisasterRecovery.Id == "Enabled")
	}

	return model
}
