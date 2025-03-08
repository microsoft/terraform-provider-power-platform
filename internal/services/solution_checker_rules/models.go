// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution_checker_rules

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataSource struct {
	helpers.TypeInfo
	SolutionCheckerRulesClient client
}

type DataSourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	Rules         []RuleModel    `tfsdk:"rules"`
}

type RuleModel struct {
	Code                       types.String `tfsdk:"code"`
	Description                types.String `tfsdk:"description"`
	Summary                    types.String `tfsdk:"summary"`
	HowToFix                   types.String `tfsdk:"how_to_fix"`
	GuidanceUrl                types.String `tfsdk:"guidance_url"`
	ComponentType              types.Int64  `tfsdk:"component_type"`
	PrimaryCategory            types.Int64  `tfsdk:"primary_category"`
	PrimaryCategoryDescription types.String `tfsdk:"primary_category_description"`
	Include                    types.Bool   `tfsdk:"include"`
	Severity                   types.Int64  `tfsdk:"severity"`
}

// Helper function to convert from DTO to Model.
func convertFromRuleDto(rule ruleDto) RuleModel {
	return RuleModel{
		Code:                       types.StringValue(rule.Code),
		Description:                types.StringValue(rule.Description),
		Summary:                    types.StringValue(rule.Summary),
		HowToFix:                   types.StringValue(rule.HowToFix),
		GuidanceUrl:                types.StringValue(rule.GuidanceUrl),
		ComponentType:              types.Int64Value(int64(rule.ComponentType)),
		PrimaryCategory:            types.Int64Value(int64(rule.PrimaryCategory)),
		PrimaryCategoryDescription: types.StringValue(getPrimaryCategoryDescription(rule.PrimaryCategory)),
		Include:                    types.BoolValue(rule.Include),
		Severity:                   types.Int64Value(int64(rule.Severity)),
	}
}

// Helper function to get primary category description.
func getPrimaryCategoryDescription(primaryCategory int) string {
	switch primaryCategory {
	case 0:
		return "Error"
	case 1:
		return "Performance"
	case 2:
		return "Security"
	case 3:
		return "Design"
	case 4:
		return "Usage"
	default:
		return "Unknown"
	}
}
