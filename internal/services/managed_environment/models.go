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
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`
	Id                           types.String   `tfsdk:"id"`
	EnvironmentId                types.String   `tfsdk:"environment_id"`
	ProtectionLevel              types.String   `tfsdk:"protection_level"`
	IsUsageInsightsDisabled      types.Bool     `tfsdk:"is_usage_insights_disabled"`
	IsGroupSharingDisabled       types.Bool     `tfsdk:"is_group_sharing_disabled"`
	MaxLimitUserSharing          types.Int64    `tfsdk:"max_limit_user_sharing"`
	LimitSharingMode             types.String   `tfsdk:"limit_sharing_mode"`
	SolutionCheckerMode          types.String   `tfsdk:"solution_checker_mode"`
	SuppressValidationEmails     types.Bool     `tfsdk:"suppress_validation_emails"`
	MakerOnboardingUrl           types.String   `tfsdk:"maker_onboarding_url"`
	MakerOnboardingMarkdown      types.String   `tfsdk:"maker_onboarding_markdown"`
	SolutionCheckerRuleOverrides types.Set      `tfsdk:"solution_checker_rule_overrides"`
}
