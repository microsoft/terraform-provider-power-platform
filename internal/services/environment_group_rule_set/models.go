// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type EnvironmentGroupRuleSetResource struct {
	helpers.TypeInfo
	EnvironmentGroupRuleSetClient client
}

type EnvironmentGroupRuleSetResourceModel struct {
	// Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String `tfsdk:"id"`
	EnvironmentGroupId types.String `tfsdk:"environment_group_id"`
	Rules              types.Set    `tfsdk:"rules"`
}

type EnvironmentGroupRuleSetRuleResourceModel struct {
	Type         types.String `tfsdk:"type"`
	ResourceType types.String `tfsdk:"resource_type"`
}
