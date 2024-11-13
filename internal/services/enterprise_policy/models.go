// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type Resource struct {
	helpers.TypeInfo
	EnterprisePolicyClient Client
}

type sourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	SystemId      types.String   `tfsdk:"system_id"`
	PolicyType    types.String   `tfsdk:"policy_type"`
}
