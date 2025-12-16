// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type ManagedEnvironmentResource struct {
	helpers.TypeInfo
	ManagedEnvironmentClient client
}

type ManagedEnvironmentResourceModel struct {
	Timeouts                               timeouts.Value `tfsdk:"timeouts"`
	Id                                     types.String   `tfsdk:"id"`
	EnvironmentId                          types.String   `tfsdk:"environment_id"`
	ProtectionLevel                        types.String   `tfsdk:"protection_level"`
	IsUsageInsightsDisabled                types.Bool     `tfsdk:"is_usage_insights_disabled"`
	IsGroupSharingDisabled                 types.Bool     `tfsdk:"is_group_sharing_disabled"`
	MaxLimitUserSharing                    types.Int64    `tfsdk:"max_limit_user_sharing"`
	LimitSharingMode                       types.String   `tfsdk:"limit_sharing_mode"`
	SolutionCheckerMode                    types.String   `tfsdk:"solution_checker_mode"`
	SuppressValidationEmails               types.Bool     `tfsdk:"suppress_validation_emails"`
	SolutionCheckerRuleOverrides           types.Set      `tfsdk:"solution_checker_rule_overrides"`
	PowerAutomateIsSharingDisabled         types.Bool     `tfsdk:"power_automate_is_sharing_disabled"`
	CopilotAllowGrantPermissionsWhenShared types.Bool     `tfsdk:"copilot_allow_grant_editor_permissions_when_shared"`
	CopilotLimitSharingMode                types.String   `tfsdk:"copilot_limit_sharing_mode"` // noLimit, excludeSharingToSecurityGroups
	CopilotMaxLimitUserSharing             types.Int64    `tfsdk:"copilot_max_limit_user_sharing"`
}

type SolutionCheckerRulesArrayDto struct {
	Value []SolutionCheckerRuleDto `json:"value"`
}

type SolutionCheckerRuleDto struct {
	Description     string `json:"description,omitempty"`
	GuidanceUrl     string `json:"guidanceUrl,omitempty"`
	Include         string `json:"include,omitempty"`
	Code            string `json:"code,omitempty"`
	Summary         string `json:"summary,omitempty"`
	ComponetType    string `json:"componetType,omitempty"`
	PrimaryCategory string `json:"primaryCategory,omitempty"`
	Severity        string `json:"severity,omitempty"`
	HowToFix        string `json:"howToFix,omitempty"`
}
