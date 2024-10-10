// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type environmentGroupRuleSetResource struct {
	helpers.TypeInfo
	EnvironmentGroupRuleSetClient client
}

type environmentGroupRuleSetResourceModel struct {
	// Timeouts         timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String `tfsdk:"id"`
	EnvironmentGroupId types.String `tfsdk:"environment_group_id"`
	Rules              types.Set    `tfsdk:"rules"`
}

type environmentGroupRuleSetRuleResourceModel struct {
	Type types.String `tfsdk:"type"`
	//ResourceType types.String `tfsdk:"resource_type"`
	Values types.Set `tfsdk:"values"`
}

type environmentRulesValuesResourceModel struct {
	ShareMode                    types.String `tfsdk:"share_mode"`
	ShareMaxLimit                types.Int32  `tfsdk:"share_max_limit"`
	InsightsEnabled              types.Bool   `tfsdk:"insights_enabled"`
	OnboardingUrl                types.String `tfsdk:"onboarding_url"`
	OnboardingMarkdown           types.String `tfsdk:"onboarding_markdown"`
	SolutionCheckerMode          types.String `tfsdk:"solution_checker_mode"`
	SendEmailsEnabled            types.Bool   `tfsdk:"send_emails_enabled"`
	PeriodInDays                 types.Int32  `tfsdk:"period_in_days"`
	AiDescriptionEnabled         types.Bool   `tfsdk:"ai_description_enabled"`
	MoveDataAcrossRegionsEnabled types.Bool   `tfsdk:"move_data_across_regions_enabled"`
	BingSearchEnabled            types.Bool   `tfsdk:"bing_search_enabled"`
}

var ruleSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"type": types.StringType,
		//"resource_type": types.StringType,
		"values": types.SetType{ElemType: environmentRulesValuesObjectType},
	},
}

var environmentRulesValuesObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"share_mode":                       types.StringType,
		"share_max_limit":                  types.Int32Type,
		"insights_enabled":                 types.BoolType,
		"onboarding_url":                   types.StringType,
		"onboarding_markdown":              types.StringType,
		"solution_checker_mode":            types.StringType,
		"send_emails_enabled":              types.BoolType,
		"period_in_days":                   types.Int32Type,
		"ai_description_enabled":           types.BoolType,
		"move_data_across_regions_enabled": types.BoolType,
		"bing_search_enabled":              types.BoolType,
	},
}
