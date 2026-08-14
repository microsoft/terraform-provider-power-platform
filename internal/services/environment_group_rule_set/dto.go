// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type environmentGroupRuleSetDto struct {
	Value []EnvironmentGroupRuleSetValueSetDto `json:"value"`
}

type EnvironmentGroupRuleSetValueSetDto struct {
	Parameters        []*environmentGroupRuleSetParameterDto       `json:"parameters"`
	Id                *string                                      `json:"id,omitempty"`
	LastModified      *string                                      `json:"lastModified,omitempty"`
	EnvironmentFilter *environmentGroupRuleSetEnvironmentFilterDto `json:"environmentFilter,omitempty"`
}

type environmentGroupRuleSetEnvironmentFilterDto struct {
	Type  string                                `json:"type"`
	Value []environmentGroupRuleSetValueTypeDto `json:"values"`
}

type environmentGroupRuleSetValueTypeDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type environmentGroupRuleSetValueDto struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type environmentGroupRuleSetParameterDto struct {
	HasStagedChanges *bool                             `json:"hasStagedChanges,omitempty"`
	Type             string                            `json:"type"`
	ResourceType     string                            `json:"resourceType"`
	Value            []environmentGroupRuleSetValueDto `json:"value"`
}

func convertEnvironmentGroupRuleSetResourceModelToDto(ctx context.Context, model environmentGroupRuleSetResourceModel) (EnvironmentGroupRuleSetValueSetDto, error) {
	dto := EnvironmentGroupRuleSetValueSetDto{}

	if !model.Id.IsUnknown() && !model.Id.IsNull() {
		dto.Id = model.Id.ValueStringPointer()
	}
	currentTime := time.Now().Format(time.RFC3339Nano)
	dto.LastModified = &currentTime
	dto.EnvironmentFilter = &environmentGroupRuleSetEnvironmentFilterDto{}
	dto.EnvironmentFilter.Type = "Include"
	dto.EnvironmentFilter.Value = append(dto.EnvironmentFilter.Value, environmentGroupRuleSetValueTypeDto{
		Id:   model.EnvironmentGroupId.ValueString(),
		Type: "EnvironmentGroup",
	})

	if !model.Rules.IsNull() && !model.Rules.IsUnknown() {
		ruleAttrs := model.Rules.Attributes()
		convertSharingControls(ctx, ruleAttrs, &dto)
		convertUsageInsights(ctx, ruleAttrs, &dto)
		convertSolutionCheckerEnforcement(ctx, ruleAttrs, &dto)
		if err := convertBackupRetention(ctx, ruleAttrs, &dto); err != nil {
			return dto, err
		}
		if err := convertAiGeneratedDesc(ctx, ruleAttrs, &dto); err != nil {
			return dto, err
		}
		if err := convertAiGenerativeSettings(ctx, ruleAttrs, &dto); err != nil {
			return dto, err
		}
	}
	return dto, nil
}

func convertAiGenerativeSettings(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) error {
	aiGenerativeSettingsObj := attrs["ai_generative_settings"]
	if !aiGenerativeSettingsObj.IsNull() && !aiGenerativeSettingsObj.IsUnknown() {
		objectValue, ok := aiGenerativeSettingsObj.(basetypes.ObjectValue)
		if !ok {
			return fmt.Errorf("expected ai_generative_settings to be of type ObjectValue, got %T", aiGenerativeSettingsObj)
		}

		var aiGenerativeSettings environmentGroupRuleSetAiGenerativeSettingsModel
		if diags := objectValue.As(ctx, &aiGenerativeSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
			return fmt.Errorf("failed to convert ai generative settings: %v", diags)
		}

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             AI_GENERATIVE_SETTINGS,
			ResourceType:     NOT_SPECIFIED,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    CROSS_GEO_COPILOT_DATA_MOVEMENT_ENABLED,
			Value: strconv.FormatBool(aiGenerativeSettings.MoveDataAcrossRegionsEnabled.ValueBool()),
		})
		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    BING_CHAT_ENABLED,
			Value: strconv.FormatBool(aiGenerativeSettings.BingSearchEnabled.ValueBool()),
		})
	}
	return nil
}

func convertAiGeneratedDesc(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) error {
	aiGeneratedDescObj := attrs["ai_generated_descriptions"]
	if !aiGeneratedDescObj.IsNull() && !aiGeneratedDescObj.IsUnknown() {
		objectValue, ok := aiGeneratedDescObj.(basetypes.ObjectValue)
		if !ok {
			return fmt.Errorf("expected ai_generated_descriptions to be of type ObjectValue, got %T", aiGeneratedDescObj)
		}

		var aiGeneratedDesc environmentGroupRuleSetAiGeneratedDescriptionsModel
		if diags := objectValue.As(ctx, &aiGeneratedDesc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
			return fmt.Errorf("failed to convert ai generated desc: %v", diags)
		}

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             AI_GENERATED_DESC,
			ResourceType:     APP,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    DISABLE_AI_GENERATED_DESCRIPTIONS,
			Value: strconv.FormatBool(!aiGeneratedDesc.AiDescriptionEnabled.ValueBool()),
		})
	}
	return nil
}

func convertBackupRetention(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) error {
	backupRetentionObj := attrs["backup_retention"]
	if !backupRetentionObj.IsNull() && !backupRetentionObj.IsUnknown() {
		objectValue, ok := backupRetentionObj.(basetypes.ObjectValue)
		if !ok {
			return fmt.Errorf("expected backup_retention to be of type ObjectValue, got %T", backupRetentionObj)
		}

		var backupRetention environmentGroupRuleSetBackupRetentionModel
		if diags := objectValue.As(ctx, &backupRetention, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
			return fmt.Errorf("failed to convert backup retention: %v", diags)
		}

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             BACKUP_RETENTION,
			ResourceType:     NOT_SPECIFIED,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    RETENTION_PERIOD,
			Value: fmt.Sprintf("%d.00:00:00", backupRetention.PeriodInDays.ValueInt32()),
		})
	}
	return nil
}

func convertSolutionCheckerEnforcement(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) {
	solutionCheckerObj := attrs["solution_checker_enforcement"]
	if !solutionCheckerObj.IsNull() && !solutionCheckerObj.IsUnknown() {
		objectValue, ok := solutionCheckerObj.(basetypes.ObjectValue)
		if !ok {
			return // Skip conversion if type assertion fails
		}

		var solutionChecker environmentGroupRuleSetSolutionCheckerEnforcementModel
		objectValue.As(ctx, &solutionChecker, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             SOLUTION_CHECKER_ENFORCEMENT,
			ResourceType:     NOT_SPECIFIED,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    SOLUTION_CHECKER_MODE,
			Value: solutionChecker.SolutionCheckerMode.ValueString(),
		})
		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    SUPPRESS_VALIDATION_EMAILS,
			Value: strconv.FormatBool(solutionChecker.SendEmailsEnabled.ValueBool()),
		})
		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    SOLUTION_CHECKER_RULE_OVERRIDES,
			Value: "",
		})
	}
}

func convertUsageInsights(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) {
	usageInsightsObj := attrs["usage_insights"]
	if !usageInsightsObj.IsNull() && !usageInsightsObj.IsUnknown() {
		objectValue, ok := usageInsightsObj.(basetypes.ObjectValue)
		if !ok {
			return // Skip conversion if type assertion fails
		}

		var usageInsights environmentGroupRuleSetUsageInsightsModel
		objectValue.As(ctx, &usageInsights, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             USAGE_INSIGHTS,
			ResourceType:     NOT_SPECIFIED,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    INCLUDE_ON_HOME_PAGE_INSIGHTS,
			Value: "false",
		})
		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    EXCLUDE_ENVIRONMENT_FROM_ANALYSIS,
			Value: strconv.FormatBool(!usageInsights.InsightsEnabled.ValueBool()),
		})
	}
}

func convertSharingControls(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) {
	sharingControlObj := attrs["sharing_controls"]
	if !sharingControlObj.IsNull() && !sharingControlObj.IsUnknown() {
		objectValue, ok := sharingControlObj.(basetypes.ObjectValue)
		if !ok {
			return // Skip conversion if type assertion fails
		}

		var sharingControl environmentGroupRuleSetSharingControlsModel
		objectValue.As(ctx, &sharingControl, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             SHARING,
			ResourceType:     APP,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		if sharingControl.ShareMode.ValueString() == "no limit" {
			rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
				Id:    CAN_SHARE_WITH_SECURITY_GROUPS,
				Value: NO_LIMIT,
			})
			rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
				Id:    IS_GROUP_SHARING_DISABLED,
				Value: "false",
			})
			rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
				Id:    MAXIMUM_SHARE_LIMIT,
				Value: "-1",
			})
		} else {
			rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
				Id:    CAN_SHARE_WITH_SECURITY_GROUPS,
				Value: EXCLUDE_SHARING_TO_SECURITY_GROUPS,
			})
			rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
				Id:    IS_GROUP_SHARING_DISABLED,
				Value: "true",
			})
			rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
				Id:    MAXIMUM_SHARE_LIMIT,
				Value: sharingControl.ShareMaxLimit.String(),
			})
		}
	}
}

func convertEnvironmentGroupRuleSetDtoToModel(environmentGroupId string, ruleSet *EnvironmentGroupRuleSetValueSetDto, policy *ruleBasedPolicyDto) (*environmentGroupRuleSetResourceModel, error) {
	rulesModel, err := convertRulesDtoToModel(ruleSet, policy)
	if err != nil {
		return nil, err
	}

	model := environmentGroupRuleSetResourceModel{
		EnvironmentGroupId: types.StringValue(environmentGroupId),
		Id:                 types.StringNull(),
		PolicyId:           types.StringNull(),
		Rules:              rulesModel,
	}

	if ruleSet != nil {
		model.Id = types.StringPointerValue(ruleSet.Id)
	}
	if policy != nil {
		model.PolicyId = types.StringValue(policy.Id)
	}

	return &model, nil
}

func convertRulesDtoToModel(dto *EnvironmentGroupRuleSetValueSetDto, policy *ruleBasedPolicyDto) (basetypes.ObjectValue, error) {
	sharingControlType, sharingControlValue, err := convertSharingControlsDtoToModel(getParameterByType(dto, SHARING))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	usageInsightsType, usageInsightsValue, err := convertUsageInsightsDtoToModel(getParameterByType(dto, USAGE_INSIGHTS))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	makerWelcomeContentType, makerWelcomeContentValue, err := convertMakerWelcomeContentDtoToModel(policyRuleSetsOf(policy), getParameterByType(dto, MAKER_WELCOME_CONTENT))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	solutionChekerEnforcementType, solutionChekerEnforcementValue, err := convertSolutionCheckerEnforcementDtoToModel(getParameterByType(dto, SOLUTION_CHECKER_ENFORCEMENT))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	backupRetentionType, backupRetentionValue, err := convertBackupRetentionDtoToModel(getParameterByType(dto, BACKUP_RETENTION))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	aiGeneratedDescType, aiGeneratedDescValue, err := convertAiGeneratedDescDtoToModel(getParameterByType(dto, AI_GENERATED_DESC))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	aiGenerativeSettingsType, aiGenerativeSettingsValue, err := convertAiGenerativeSettingsDtoToModel(getParameterByType(dto, AI_GENERATIVE_SETTINGS))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	policyRuleSets := []policyRuleSetDto{}
	if policy != nil {
		policyRuleSets = policy.RuleSets
	}
	advancedConnectorOnlyType, advancedConnectorOnlyValue, err := convertAdvancedConnectorPoliciesOnlyDtoToModel(policyRuleSets)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	cspType, cspValue, err := convertContentSecurityPolicyDtoToModel(policyRuleSets)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	connMgmtType, connMgmtValue, err := convertConnectorManagementDtoToModel(policyRuleSets)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	atr := map[string]attr.Type{
		"sharing_controls":                 sharingControlType,
		"usage_insights":                   usageInsightsType,
		"maker_welcome_content":            makerWelcomeContentType,
		"solution_checker_enforcement":     solutionChekerEnforcementType,
		"backup_retention":                 backupRetentionType,
		"ai_generated_descriptions":        aiGeneratedDescType,
		"ai_generative_settings":           aiGenerativeSettingsType,
		"advanced_connector_policies_only": advancedConnectorOnlyType,
		"content_security_policy":          cspType,
		"advanced_connector_policies":      connMgmtType,
	}

	hasLegacyRules := dto != nil && len(dto.Parameters) > 0
	hasPolicyRules := len(policyRuleSets) > 0
	if !hasLegacyRules && !hasPolicyRules {
		return types.ObjectNull(atr), nil
	}

	attrV := map[string]attr.Value{
		"sharing_controls":                 sharingControlValue,
		"usage_insights":                   usageInsightsValue,
		"maker_welcome_content":            makerWelcomeContentValue,
		"solution_checker_enforcement":     solutionChekerEnforcementValue,
		"backup_retention":                 backupRetentionValue,
		"ai_generated_descriptions":        aiGeneratedDescValue,
		"ai_generative_settings":           aiGenerativeSettingsValue,
		"advanced_connector_policies_only": advancedConnectorOnlyValue,
		"content_security_policy":          cspValue,
		"advanced_connector_policies":      connMgmtValue,
	}

	return types.ObjectValueMust(atr, attrV), nil
}

func convertAiGenerativeSettingsDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"move_data_across_regions_enabled": types.BoolType,
		"bing_search_enabled":              types.BoolType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	crossGeoCopilotDataMovementEnabled := tryGetRuleValueFromDto(dto.Value, CROSS_GEO_COPILOT_DATA_MOVEMENT_ENABLED)
	if crossGeoCopilotDataMovementEnabled == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", CROSS_GEO_COPILOT_DATA_MOVEMENT_ENABLED)
	}
	bingSearchEnabled := tryGetRuleValueFromDto(dto.Value, BING_CHAT_ENABLED)
	if bingSearchEnabled == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", BING_CHAT_ENABLED)
	}

	attrValue := map[string]attr.Value{}
	attrValue["move_data_across_regions_enabled"] = types.BoolValue(crossGeoCopilotDataMovementEnabled.Value == "true")
	attrValue["bing_search_enabled"] = types.BoolValue(bingSearchEnabled.Value == "true")

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

func convertAiGeneratedDescDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"ai_description_enabled": types.BoolType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	aiDescriptionEnabled := tryGetRuleValueFromDto(dto.Value, DISABLE_AI_GENERATED_DESCRIPTIONS)
	if aiDescriptionEnabled == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", DISABLE_AI_GENERATED_DESCRIPTIONS)
	}

	attrValue := map[string]attr.Value{}
	attrValue["ai_description_enabled"] = types.BoolValue(aiDescriptionEnabled.Value == "false")

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

func convertBackupRetentionDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"period_in_days": types.Int32Type,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	periodInDays := tryGetRuleValueFromDto(dto.Value, RETENTION_PERIOD)
	if periodInDays == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", RETENTION_PERIOD)
	}

	attrValue := map[string]attr.Value{}

	// expected format for 14 days would be 14.00:00:00
	v := strings.Split(periodInDays.Value, ".")
	if len(v) != 2 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), errors.New("invalid format for period in days")
	}
	period, err := strconv.ParseInt(v[0], 10, 32)
	if err != nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), errors.New("invalid format when parsing period in days")
	}
	attrValue["period_in_days"] = types.Int32Value(int32(period))

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

func convertSolutionCheckerEnforcementDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"solution_checker_mode": types.StringType,
		"send_emails_enabled":   types.BoolType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	solutionCheckerMode := tryGetRuleValueFromDto(dto.Value, SOLUTION_CHECKER_MODE)
	if solutionCheckerMode == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", SOLUTION_CHECKER_MODE)
	}
	sendEmailsEnabled := tryGetRuleValueFromDto(dto.Value, SUPPRESS_VALIDATION_EMAILS)
	if sendEmailsEnabled == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", SUPPRESS_VALIDATION_EMAILS)
	}

	attrValue := map[string]attr.Value{}
	attrValue["solution_checker_mode"] = types.StringValue(solutionCheckerMode.Value)
	attrValue["send_emails_enabled"] = types.BoolValue(sendEmailsEnabled.Value == "true")

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

// convertMakerWelcomeContentDtoToModel prefers the policy rule set, falling back to the legacy
// parameter for rule sets written before maker welcome content moved to the policies API.
func convertMakerWelcomeContentDtoToModel(policyRuleSets []policyRuleSetDto, dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"maker_onboarding_url":      types.StringType,
		"maker_onboarding_markdown": types.StringType,
	}
	objType := types.ObjectType{AttrTypes: attrType}

	if ruleSet := findPolicyRuleSet(policyRuleSets, MAKER_ONBOARDING_CONTENT_ID); ruleSet != nil {
		asString := func(key string) types.String {
			if s, ok := ruleSet.Inputs[key].(string); ok {
				return types.StringValue(s)
			}
			return types.StringValue("")
		}
		return objType, types.ObjectValueMust(attrType, map[string]attr.Value{
			"maker_onboarding_url":      asString(MAKER_ONBOARDING_URL),
			"maker_onboarding_markdown": asString(MAKER_ONBOARDING_MARKDOWN),
		}), nil
	}

	if dto == nil || len(dto.Value) == 0 {
		return objType, types.ObjectNull(attrType), nil
	}

	makerOnboardingUrl := tryGetRuleValueFromDto(dto.Value, MAKER_ONBOARDING_URL)
	if makerOnboardingUrl == nil {
		return objType, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", MAKER_ONBOARDING_URL)
	}
	makerOnboardingMarkdown := tryGetRuleValueFromDto(dto.Value, MAKER_ONBOARDING_MARKDOWN)
	if makerOnboardingMarkdown == nil {
		return objType, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", MAKER_ONBOARDING_MARKDOWN)
	}

	return objType, types.ObjectValueMust(attrType, map[string]attr.Value{
		"maker_onboarding_url":      types.StringValue(makerOnboardingUrl.Value),
		"maker_onboarding_markdown": types.StringValue(makerOnboardingMarkdown.Value),
	}), nil
}

func convertUsageInsightsDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"insights_enabled": types.BoolType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	excludeEnvFromAnalisys := tryGetRuleValueFromDto(dto.Value, EXCLUDE_ENVIRONMENT_FROM_ANALYSIS)
	if excludeEnvFromAnalisys == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", EXCLUDE_ENVIRONMENT_FROM_ANALYSIS)
	}
	attrValue := map[string]attr.Value{}
	attrValue["insights_enabled"] = types.BoolValue(excludeEnvFromAnalisys.Value == "false")

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

func convertSharingControlsDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"share_mode":      types.StringType,
		"share_max_limit": types.NumberType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	canShareWithSecurityGroups := tryGetRuleValueFromDto(dto.Value, CAN_SHARE_WITH_SECURITY_GROUPS)
	if canShareWithSecurityGroups == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", CAN_SHARE_WITH_SECURITY_GROUPS)
	}
	maximumShareLimit := tryGetRuleValueFromDto(dto.Value, MAXIMUM_SHARE_LIMIT)
	if maximumShareLimit == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", MAXIMUM_SHARE_LIMIT)
	}
	maxLimitValue, _ := strconv.ParseFloat(maximumShareLimit.Value, 64)

	attrValue := map[string]attr.Value{}
	if canShareWithSecurityGroups.Value == NO_LIMIT {
		attrValue["share_mode"] = types.StringValue("no limit")
		attrValue["share_max_limit"] = types.NumberNull() // -1 value is considered as null in terraform
	} else {
		attrValue["share_mode"] = types.StringValue("exclude sharing with security groups")
		attrValue["share_max_limit"] = types.NumberValue(big.NewFloat(maxLimitValue))
	}

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

func tryGetRuleValueFromDto(values []environmentGroupRuleSetValueDto, valueId string) *environmentGroupRuleSetValueDto {
	for _, value := range values {
		if value.Id == valueId {
			return &value
		}
	}
	return nil
}

func getParameterByType(params *EnvironmentGroupRuleSetValueSetDto, paramType string) *environmentGroupRuleSetParameterDto {
	if params == nil {
		return nil
	}
	for paramInx := range params.Parameters {
		if params.Parameters[paramInx].Type == paramType {
			return params.Parameters[paramInx]
		}
	}
	return nil
}

// --- Rule-based policy DTOs. These rule sets are served by the rule-based policies API. ---

type ruleBasedPolicyDto struct {
	Id           string             `json:"id,omitempty"`
	Name         string             `json:"name,omitempty"`
	TenantId     string             `json:"tenantId,omitempty"`
	LastModified string             `json:"lastModified,omitempty"`
	RuleSetCount int                `json:"ruleSetCount,omitempty"`
	RuleSets     []policyRuleSetDto `json:"ruleSets,omitempty"`
}

type ruleBasedPolicyRequestDto struct {
	Name     string             `json:"name"`
	RuleSets []policyRuleSetDto `json:"ruleSets"`
}

type policyRuleSetDto struct {
	Id      string         `json:"id"`
	Version string         `json:"version"`
	Inputs  map[string]any `json:"inputs"`
}

type ruleAssignmentDto struct {
	PolicyId     string `json:"policyId"`
	ResourceId   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	RuleSetCount int    `json:"ruleSetCount"`
	TenantId     string `json:"tenantId"`
}

type ruleAssignmentsResponseDto struct {
	Value []ruleAssignmentDto `json:"value"`
}

type policyAssignmentRequestDto struct {
	AssignmentOverrides []any `json:"assignmentOverrides"`
}

type removeRuleRequestDto struct {
	Name     string             `json:"name"`
	RuleSets []removeRuleSetDto `json:"ruleSets"`
}

type removeRuleSetDto struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type allowedConnectorDto struct {
	AllowedActionsMode         string   `json:"AllowedActionsMode"`
	AllowedConnectionTypesMode string   `json:"AllowedConnectionTypesMode"`
	AllowedConnector           string   `json:"AllowedConnector"`
	AllowedActions             []string `json:"AllowedActions,omitempty"`
}

// managedPolicyRuleSetIds maps policy rule set IDs to the versions this resource manages.
var managedPolicyRuleSetIds = map[string]string{
	ADVANCED_CONNECTOR_POLICIES_ONLY_ID: ADVANCED_CONNECTOR_POLICIES_VERSION,
	CONTENT_SECURITY_POLICY_ID:          CONTENT_SECURITY_POLICY_VERSION,
	CONNECTOR_MANAGEMENT_ID:             CONNECTOR_MANAGEMENT_VERSION,
	MAKER_ONBOARDING_CONTENT_ID:         MAKER_ONBOARDING_CONTENT_VERSION,
}

// policyRuleAttributeNames are the rules under `rules` that are backed by the rule-based policies API.
var policyRuleAttributeNames = []string{"advanced_connector_policies_only", "content_security_policy", "advanced_connector_policies", "maker_welcome_content"}

// hasPolicyRules reports whether the configuration sets any rule backed by the rule-based policies API.
func hasPolicyRules(rules basetypes.ObjectValue) bool {
	if rules.IsNull() || rules.IsUnknown() {
		return false
	}
	attrs := rules.Attributes()
	for _, name := range policyRuleAttributeNames {
		if v, ok := attrs[name]; ok && v != nil && !v.IsNull() && !v.IsUnknown() {
			return true
		}
	}
	return false
}

// hasLegacyRules reports whether the configuration sets any rule backed by the legacy rule sets API.
func hasLegacyRules(rules basetypes.ObjectValue) bool {
	if rules.IsNull() || rules.IsUnknown() {
		return false
	}
	policyNames := map[string]bool{}
	for _, name := range policyRuleAttributeNames {
		policyNames[name] = true
	}
	for name, v := range rules.Attributes() {
		if policyNames[name] {
			continue
		}
		if v != nil && !v.IsNull() && !v.IsUnknown() {
			return true
		}
	}
	return false
}

// mergePolicyRuleSets replaces the rule sets managed by this resource while preserving any others.
func mergePolicyRuleSets(existing, desired []policyRuleSetDto) []policyRuleSetDto {
	merged := make([]policyRuleSetDto, 0, len(existing)+len(desired))
	for _, rs := range existing {
		if _, managed := managedPolicyRuleSetIds[rs.Id]; !managed {
			merged = append(merged, rs)
		}
	}
	return append(merged, desired...)
}

// terraformManagedPolicyName is the policy name this resource assigns when it creates a policy.
// A policy carrying this name was created by this resource and can safely be deleted with it.
func terraformManagedPolicyName(environmentGroupId string) string {
	return fmt.Sprintf("terraform-managed-%s", environmentGroupId)
}

func convertRulesModelToPolicyDto(ctx context.Context, model environmentGroupRuleSetResourceModel) (ruleBasedPolicyRequestDto, error) {
	dto := ruleBasedPolicyRequestDto{
		Name:     terraformManagedPolicyName(model.EnvironmentGroupId.ValueString()),
		RuleSets: make([]policyRuleSetDto, 0),
	}

	if model.Rules.IsNull() || model.Rules.IsUnknown() {
		return dto, nil
	}

	attrs := model.Rules.Attributes()

	if obj, ok := attrs["advanced_connector_policies_only"].(basetypes.ObjectValue); ok && !obj.IsNull() && !obj.IsUnknown() {
		var advancedConnector environmentGroupRuleSetAdvancedConnectorPoliciesOnlyModel
		if diags := obj.As(ctx, &advancedConnector, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
			return dto, fmt.Errorf("failed to convert advanced_connector_policies_only: %v", diags)
		}
		dto.RuleSets = append(dto.RuleSets, policyRuleSetDto{
			Id:      ADVANCED_CONNECTOR_POLICIES_ONLY_ID,
			Version: ADVANCED_CONNECTOR_POLICIES_VERSION,
			Inputs:  map[string]any{ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY: advancedConnector.Enabled.ValueBool()},
		})
	}

	if obj, ok := attrs["content_security_policy"].(basetypes.ObjectValue); ok && !obj.IsNull() && !obj.IsUnknown() {
		cspRuleSet, err := convertContentSecurityPolicyModelToDto(ctx, obj)
		if err != nil {
			return dto, fmt.Errorf("failed to convert content_security_policy: %w", err)
		}
		dto.RuleSets = append(dto.RuleSets, cspRuleSet)
	}

	if obj, ok := attrs["advanced_connector_policies"].(basetypes.ObjectValue); ok && !obj.IsNull() && !obj.IsUnknown() {
		connRuleSet, err := convertConnectorManagementModelToDto(obj)
		if err != nil {
			return dto, fmt.Errorf("failed to convert advanced_connector_policies: %w", err)
		}
		dto.RuleSets = append(dto.RuleSets, connRuleSet)
	}

	if obj, ok := attrs["maker_welcome_content"].(basetypes.ObjectValue); ok && !obj.IsNull() && !obj.IsUnknown() {
		var makerWelcomeContent environmentGroupRuleSetMakerWelcomeContentModel
		if diags := obj.As(ctx, &makerWelcomeContent, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
			return dto, fmt.Errorf("failed to convert maker_welcome_content: %v", diags)
		}
		dto.RuleSets = append(dto.RuleSets, policyRuleSetDto{
			Id:      MAKER_ONBOARDING_CONTENT_ID,
			Version: MAKER_ONBOARDING_CONTENT_VERSION,
			Inputs: map[string]any{
				MAKER_ONBOARDING_URL:              makerWelcomeContent.MakerOnboardingUrl.ValueString(),
				MAKER_ONBOARDING_MARKDOWN:         makerWelcomeContent.MakerOnboardingMarkdown.ValueString(),
				MAKER_ONBOARDING_PORTALS:          "",
				MAKER_ONBOARDING_TIMESTAMP:        time.Now().UTC().Format(http.TimeFormat),
				MAKER_ONBOARDING_CONSENT_REQUIRED: false,
			},
		})
	}

	return dto, nil
}

// policyRuleSetsOf returns the policy's rule sets, tolerating a nil policy.
func policyRuleSetsOf(policy *ruleBasedPolicyDto) []policyRuleSetDto {
	if policy == nil {
		return nil
	}
	return policy.RuleSets
}

func findPolicyRuleSet(ruleSets []policyRuleSetDto, id string) *policyRuleSetDto {
	for i := range ruleSets {
		if ruleSets[i].Id == id {
			return &ruleSets[i]
		}
	}
	return nil
}

func convertAdvancedConnectorPoliciesOnlyDtoToModel(ruleSets []policyRuleSetDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{"enabled": types.BoolType}

	ruleSet := findPolicyRuleSet(ruleSets, ADVANCED_CONNECTOR_POLICIES_ONLY_ID)
	if ruleSet == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	enabledValue, ok := ruleSet.Inputs[ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY]
	if !ok {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s input not found in response", ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY)
	}
	enabled, ok := enabledValue.(bool)
	if !ok {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s input is not a boolean", ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY)
	}

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, map[string]attr.Value{"enabled": types.BoolValue(enabled)}), nil
}

// CSP configuration DTO types. The API stores these as JSON strings inside the rule set inputs.
type cspConfigurationDto struct {
	ImgSrc        *cspDirectiveDto `json:"Img-Src,omitempty"`
	StyleSrc      *cspDirectiveDto `json:"Style-Src,omitempty"`
	FormAction    *cspDirectiveDto `json:"Form-Action,omitempty"`
	FrameSrc      *cspDirectiveDto `json:"Frame-Src,omitempty"`
	ConnectSrc    *cspDirectiveDto `json:"Connect-Src,omitempty"`
	FontSrc       *cspDirectiveDto `json:"Font-Src,omitempty"`
	ScriptSrc     *cspDirectiveDto `json:"Script-Src,omitempty"`
	FrameAncestor *cspDirectiveDto `json:"Frame-Ancestor,omitempty"`
}

type cspDirectiveDto struct {
	Sources []cspSourceDto `json:"sources"`
}

type cspSourceDto struct {
	Source string `json:"source"`
}

func cspDirectiveAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"img_src":        types.ListType{ElemType: types.StringType},
		"style_src":      types.ListType{ElemType: types.StringType},
		"form_action":    types.ListType{ElemType: types.StringType},
		"frame_src":      types.ListType{ElemType: types.StringType},
		"connect_src":    types.ListType{ElemType: types.StringType},
		"font_src":       types.ListType{ElemType: types.StringType},
		"script_src":     types.ListType{ElemType: types.StringType},
		"frame_ancestor": types.ListType{ElemType: types.StringType},
	}
}

func cspConfigurationAttrTypes() map[string]attr.Type {
	attrTypes := cspDirectiveAttrTypes()
	attrTypes["strict_csp"] = types.BoolType
	return attrTypes
}

func contentSecurityPolicyAttrTypes() map[string]attr.Type {
	configType := types.ObjectType{AttrTypes: cspConfigurationAttrTypes()}
	return map[string]attr.Type{
		"enabled":                     types.BoolType,
		"enabled_for_canvas":          types.BoolType,
		"enabled_for_code_apps":       types.BoolType,
		"report_uri":                  types.StringType,
		"reporting_endpoint":          types.StringType,
		"configuration":               configType,
		"configuration_for_canvas":    configType,
		"configuration_for_code_apps": types.ObjectType{AttrTypes: cspDirectiveAttrTypes()},
	}
}

func convertContentSecurityPolicyModelToDto(ctx context.Context, objectValue basetypes.ObjectValue) (policyRuleSetDto, error) {
	var csp environmentGroupRuleSetContentSecurityPolicyModel
	if diags := objectValue.As(ctx, &csp, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
		return policyRuleSetDto{}, fmt.Errorf("failed to convert content_security_policy: %v", diags)
	}

	inputs := map[string]any{
		CSP_IS_ENABLED_KEY:               csp.Enabled.ValueBool(),
		CSP_IS_ENABLED_FOR_CANVAS_KEY:    csp.EnabledForCanvas.ValueBool(),
		CSP_IS_ENABLED_FOR_CODE_APPS_KEY: csp.EnabledForCodeApps.ValueBool(),
		CSP_OPTIONS_KEY:                  calculateCspOptions(csp.Configuration, csp.ConfigurationForCanvas),
		CSP_REPORT_URI_KEY:               nil,
		CSP_REPORTING_ENDPOINT_KEY:       nil,
	}

	if !csp.ReportUri.IsNull() && !csp.ReportUri.IsUnknown() {
		inputs[CSP_REPORT_URI_KEY] = csp.ReportUri.ValueString()
	}
	if !csp.ReportingEndpoint.IsNull() && !csp.ReportingEndpoint.IsUnknown() {
		inputs[CSP_REPORTING_ENDPOINT_KEY] = csp.ReportingEndpoint.ValueString()
	}

	configJson, err := convertCspConfigModelToJson(csp.Configuration)
	if err != nil {
		return policyRuleSetDto{}, fmt.Errorf("failed to convert configuration: %w", err)
	}
	inputs[CSP_CONFIGURATION_KEY] = configJson

	canvasJson, err := convertCspConfigModelToJson(csp.ConfigurationForCanvas)
	if err != nil {
		return policyRuleSetDto{}, fmt.Errorf("failed to convert configuration_for_canvas: %w", err)
	}
	inputs[CSP_CONFIGURATION_FOR_CANVAS_KEY] = canvasJson

	codeAppsJson, err := convertCspConfigModelToJson(csp.ConfigurationForCodeApps)
	if err != nil {
		return policyRuleSetDto{}, fmt.Errorf("failed to convert configuration_for_code_apps: %w", err)
	}
	inputs[CSP_CONFIGURATION_FOR_CODE_APPS_KEY] = codeAppsJson

	return policyRuleSetDto{Id: CONTENT_SECURITY_POLICY_ID, Version: CONTENT_SECURITY_POLICY_VERSION, Inputs: inputs}, nil
}

func getStrictCspFromConfig(configObj types.Object) bool {
	if configObj.IsNull() || configObj.IsUnknown() {
		return false
	}
	if b, ok := configObj.Attributes()["strict_csp"].(types.Bool); ok && !b.IsNull() && !b.IsUnknown() {
		return b.ValueBool()
	}
	return false
}

// calculateCspOptions builds the ContentSecurityPolicyOptions bit field: 1 = model-driven strict, 8 = canvas strict.
func calculateCspOptions(modelDrivenConfig, canvasConfig types.Object) int64 {
	var options int64
	if getStrictCspFromConfig(modelDrivenConfig) {
		options |= 1
	}
	if getStrictCspFromConfig(canvasConfig) {
		options |= 8
	}
	return options
}

func convertCspConfigModelToJson(configObj types.Object) (string, error) {
	if configObj.IsNull() || configObj.IsUnknown() {
		return "{}", nil
	}

	attrs := configObj.Attributes()
	directive := func(key string) *cspDirectiveDto {
		list, ok := attrs[key].(types.List)
		if !ok || list.IsNull() || list.IsUnknown() {
			return nil
		}
		sources := make([]cspSourceDto, 0, len(list.Elements()))
		for _, elem := range list.Elements() {
			if strVal, ok := elem.(types.String); ok {
				sources = append(sources, cspSourceDto{Source: strVal.ValueString()})
			}
		}
		return &cspDirectiveDto{Sources: sources}
	}

	jsonBytes, err := json.Marshal(cspConfigurationDto{
		ImgSrc:        directive("img_src"),
		StyleSrc:      directive("style_src"),
		FormAction:    directive("form_action"),
		FrameSrc:      directive("frame_src"),
		ConnectSrc:    directive("connect_src"),
		FontSrc:       directive("font_src"),
		ScriptSrc:     directive("script_src"),
		FrameAncestor: directive("frame_ancestor"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal CSP configuration to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func convertContentSecurityPolicyDtoToModel(ruleSets []policyRuleSetDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := contentSecurityPolicyAttrTypes()
	objType := types.ObjectType{AttrTypes: attrType}

	ruleSet := findPolicyRuleSet(ruleSets, CONTENT_SECURITY_POLICY_ID)
	if ruleSet == nil {
		return objType, types.ObjectNull(attrType), nil
	}
	inputs := ruleSet.Inputs

	getBool := func(key string) types.Bool {
		if b, ok := inputs[key].(bool); ok {
			return types.BoolValue(b)
		}
		return types.BoolValue(false)
	}

	getNullableString := func(key string) types.String {
		if s, ok := inputs[key].(string); ok && s != "" {
			return types.StringValue(s)
		}
		return types.StringNull()
	}

	var options int64
	switch n := inputs[CSP_OPTIONS_KEY].(type) {
	case float64:
		options = int64(n)
	case int64:
		options = n
	}

	getConfig := func(key string, strict bool) (basetypes.ObjectValue, error) {
		jsonStr, ok := inputs[key].(string)
		if !ok {
			return types.ObjectNull(cspConfigurationAttrTypes()), nil
		}
		return convertCspConfigJsonToModel(jsonStr, strict)
	}

	configValue, err := getConfig(CSP_CONFIGURATION_KEY, options&1 != 0)
	if err != nil {
		return objType, types.ObjectNull(attrType), fmt.Errorf("failed to parse %s: %w", CSP_CONFIGURATION_KEY, err)
	}
	canvasValue, err := getConfig(CSP_CONFIGURATION_FOR_CANVAS_KEY, options&8 != 0)
	if err != nil {
		return objType, types.ObjectNull(attrType), fmt.Errorf("failed to parse %s: %w", CSP_CONFIGURATION_FOR_CANVAS_KEY, err)
	}

	codeAppsValue := types.ObjectNull(cspDirectiveAttrTypes())
	if jsonStr, ok := inputs[CSP_CONFIGURATION_FOR_CODE_APPS_KEY].(string); ok && jsonStr != "" && jsonStr != "{}" {
		attrValues, parseErr := parseCspConfigDirectives(jsonStr)
		if parseErr != nil {
			return objType, types.ObjectNull(attrType), fmt.Errorf("failed to parse %s: %w", CSP_CONFIGURATION_FOR_CODE_APPS_KEY, parseErr)
		}
		codeAppsValue = types.ObjectValueMust(cspDirectiveAttrTypes(), attrValues)
	}

	attrValue := map[string]attr.Value{
		"enabled":                     getBool(CSP_IS_ENABLED_KEY),
		"enabled_for_canvas":          getBool(CSP_IS_ENABLED_FOR_CANVAS_KEY),
		"enabled_for_code_apps":       getBool(CSP_IS_ENABLED_FOR_CODE_APPS_KEY),
		"report_uri":                  getNullableString(CSP_REPORT_URI_KEY),
		"reporting_endpoint":          getNullableString(CSP_REPORTING_ENDPOINT_KEY),
		"configuration":               configValue,
		"configuration_for_canvas":    canvasValue,
		"configuration_for_code_apps": codeAppsValue,
	}

	return objType, types.ObjectValueMust(attrType, attrValue), nil
}

func parseCspConfigDirectives(jsonStr string) (map[string]attr.Value, error) {
	var dto cspConfigurationDto
	if err := json.Unmarshal([]byte(jsonStr), &dto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CSP configuration JSON: %w", err)
	}

	directiveToList := func(d *cspDirectiveDto) types.List {
		if d == nil || len(d.Sources) == 0 {
			return types.ListNull(types.StringType)
		}
		elems := make([]attr.Value, len(d.Sources))
		for i, s := range d.Sources {
			elems[i] = types.StringValue(s.Source)
		}
		return types.ListValueMust(types.StringType, elems)
	}

	return map[string]attr.Value{
		"img_src":        directiveToList(dto.ImgSrc),
		"style_src":      directiveToList(dto.StyleSrc),
		"form_action":    directiveToList(dto.FormAction),
		"frame_src":      directiveToList(dto.FrameSrc),
		"connect_src":    directiveToList(dto.ConnectSrc),
		"font_src":       directiveToList(dto.FontSrc),
		"script_src":     directiveToList(dto.ScriptSrc),
		"frame_ancestor": directiveToList(dto.FrameAncestor),
	}, nil
}

func convertCspConfigJsonToModel(jsonStr string, strictCsp bool) (basetypes.ObjectValue, error) {
	attrType := cspConfigurationAttrTypes()
	if jsonStr == "" || jsonStr == "{}" {
		return types.ObjectNull(attrType), nil
	}

	attrValues, err := parseCspConfigDirectives(jsonStr)
	if err != nil {
		return types.ObjectNull(attrType), err
	}
	attrValues["strict_csp"] = types.BoolValue(strictCsp)

	return types.ObjectValueMust(attrType, attrValues), nil
}

const (
	actionsModeAllAllowed     = "all_allowed"
	actionsModeSomeAllowed    = "some_allowed"
	apiActionsModeAllAllowed  = "AllAllowed"
	apiActionsModeSomeAllowed = "SomeAllowed"
)

func allowedConnectorAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"connector_id":    types.StringType,
		"actions_mode":    types.StringType,
		"allowed_actions": types.ListType{ElemType: types.StringType},
	}
}

func advancedConnectorPoliciesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_connectors": types.ListType{ElemType: types.ObjectType{AttrTypes: allowedConnectorAttrTypes()}},
	}
}

func convertConnectorManagementModelToDto(objectValue basetypes.ObjectValue) (policyRuleSetDto, error) {
	connectors := make([]allowedConnectorDto, 0)

	if connList, ok := objectValue.Attributes()["allowed_connectors"].(types.List); ok && !connList.IsNull() && !connList.IsUnknown() {
		for _, elem := range connList.Elements() {
			objVal, ok := elem.(types.Object)
			if !ok {
				continue
			}
			connector, err := convertAllowedConnectorModelToDto(objVal.Attributes())
			if err != nil {
				return policyRuleSetDto{}, err
			}
			connectors = append(connectors, connector)
		}
	}

	return policyRuleSetDto{
		Id:      CONNECTOR_MANAGEMENT_ID,
		Version: CONNECTOR_MANAGEMENT_VERSION,
		Inputs:  map[string]any{CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY: connectors},
	}, nil
}

func convertAllowedConnectorModelToDto(attrs map[string]attr.Value) (allowedConnectorDto, error) {
	connectorId := ""
	if v, ok := attrs["connector_id"].(types.String); ok {
		connectorId = v.ValueString()
	}

	actionsMode := actionsModeAllAllowed
	if v, ok := attrs["actions_mode"].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
		actionsMode = v.ValueString()
	}

	var apiActionsMode string
	switch actionsMode {
	case actionsModeAllAllowed:
		apiActionsMode = apiActionsModeAllAllowed
	case actionsModeSomeAllowed:
		apiActionsMode = apiActionsModeSomeAllowed
	default:
		return allowedConnectorDto{}, fmt.Errorf("invalid actions_mode %q: must be %q or %q", actionsMode, actionsModeAllAllowed, actionsModeSomeAllowed)
	}

	connector := allowedConnectorDto{
		AllowedActionsMode:         apiActionsMode,
		AllowedConnectionTypesMode: apiActionsModeAllAllowed,
		AllowedConnector:           connectorId,
	}
	if !strings.HasPrefix(connectorId, CONNECTOR_API_PREFIX) {
		connector.AllowedConnector = CONNECTOR_API_PREFIX + connectorId
	}

	if actionsMode != actionsModeSomeAllowed {
		return connector, nil
	}

	if actionsAttr, ok := attrs["allowed_actions"].(types.List); ok && !actionsAttr.IsNull() && !actionsAttr.IsUnknown() {
		actions := make([]string, 0, len(actionsAttr.Elements()))
		for _, a := range actionsAttr.Elements() {
			if s, ok := a.(types.String); ok {
				actions = append(actions, s.ValueString())
			}
		}
		connector.AllowedActions = actions
	}

	return connector, nil
}

func convertConnectorManagementDtoToModel(ruleSets []policyRuleSetDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := advancedConnectorPoliciesAttrTypes()
	objType := types.ObjectType{AttrTypes: attrType}
	elemType := types.ObjectType{AttrTypes: allowedConnectorAttrTypes()}

	ruleSet := findPolicyRuleSet(ruleSets, CONNECTOR_MANAGEMENT_ID)
	if ruleSet == nil {
		return objType, types.ObjectNull(attrType), nil
	}

	// Inputs arrive as []any from JSON decoding; round-trip through JSON to get typed values.
	rawList, ok := ruleSet.Inputs[CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY]
	if !ok || rawList == nil {
		return objType, types.ObjectValueMust(attrType, map[string]attr.Value{
			"allowed_connectors": types.ListValueMust(elemType, []attr.Value{}),
		}), nil
	}

	rawJson, err := json.Marshal(rawList)
	if err != nil {
		return objType, types.ObjectNull(attrType), fmt.Errorf("failed to re-marshal %s: %w", CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY, err)
	}
	connectors := []allowedConnectorDto{}
	if err := json.Unmarshal(rawJson, &connectors); err != nil {
		return objType, types.ObjectNull(attrType), fmt.Errorf("failed to unmarshal %s: %w", CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY, err)
	}

	elems := make([]attr.Value, 0, len(connectors))
	for _, connector := range connectors {
		actionsMode := actionsModeAllAllowed
		allowedActions := types.ListNull(types.StringType)
		if connector.AllowedActionsMode == apiActionsModeSomeAllowed {
			actionsMode = actionsModeSomeAllowed
			actionElems := make([]attr.Value, 0, len(connector.AllowedActions))
			for _, action := range connector.AllowedActions {
				actionElems = append(actionElems, types.StringValue(action))
			}
			allowedActions = types.ListValueMust(types.StringType, actionElems)
		}

		elems = append(elems, types.ObjectValueMust(allowedConnectorAttrTypes(), map[string]attr.Value{
			"connector_id":    types.StringValue(strings.TrimPrefix(connector.AllowedConnector, CONNECTOR_API_PREFIX)),
			"actions_mode":    types.StringValue(actionsMode),
			"allowed_actions": allowedActions,
		}))
	}

	return objType, types.ObjectValueMust(attrType, map[string]attr.Value{
		"allowed_connectors": types.ListValueMust(elemType, elems),
	}), nil
}
