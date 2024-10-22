// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
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
	Rules              types.Object `tfsdk:"rules"`
}

type environmentGroupRuleSetSharingControlsModel struct {
	ShareMode     types.String `tfsdk:"share_mode"`
	ShareMaxLimit types.Number `tfsdk:"share_max_limit"`
}

type environmentGroupRuleSetUsageInsightsModel struct {
	InsightsEnabled types.Bool `tfsdk:"insights_enabled"`
}

type environmentGroupRuleSetMakerWelcomeContentModel struct {
	MakerOnboardingUrl      types.String `tfsdk:"maker_onboarding_url"`
	MakerOnboardingMarkdown types.String `tfsdk:"maker_onboarding_markdown"`
}

type environmentGroupRuleSetSolutionCheckerEnforcementModel struct {
	SolutionCheckerMode types.String `tfsdk:"solution_checker_mode"`
	SendEmailsEnabled   types.Bool   `tfsdk:"send_emails_enabled"`
}

type environmentGroupRuleSetBackupRetentionModel struct {
	PeriodInDays types.Int32 `tfsdk:"period_in_days"`
}

type environmentGroupRuleSetAiGeneratedDescriptionsModel struct {
	AiDescriptionEnabled types.Bool `tfsdk:"ai_description_enabled"`
}

type environmentGroupRuleSetAiGenerativeSettingsModel struct {
	MoveDataAcrossRegionsEnabled types.Bool `tfsdk:"move_data_across_regions_enabled"`
	BingSearchEnabled            types.Bool `tfsdk:"bing_search_enabled"`
}
