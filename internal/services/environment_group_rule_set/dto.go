// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
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
		convertMakerWelcomeContent(ctx, ruleAttrs, &dto)
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
		var aiGenerativeSettings environmentGroupRuleSetAiGenerativeSettingsModel
		if diags := aiGenerativeSettingsObj.(basetypes.ObjectValue).As(ctx, &aiGenerativeSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
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
		var aiGeneratedDesc environmentGroupRuleSetAiGeneratedDescriptionsModel
		if diags := aiGeneratedDescObj.(basetypes.ObjectValue).As(ctx, &aiGeneratedDesc, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
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
		var backupRetention environmentGroupRuleSetBackupRetentionModel
		if diags := backupRetentionObj.(basetypes.ObjectValue).As(ctx, &backupRetention, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
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
		var solutionChecker environmentGroupRuleSetSolutionCheckerEnforcementModel
		solutionCheckerObj.(basetypes.ObjectValue).As(ctx, &solutionChecker, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

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

func convertMakerWelcomeContent(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) {
	makerWelcomeContentObj := attrs["maker_welcome_content"]
	if !makerWelcomeContentObj.IsNull() && !makerWelcomeContentObj.IsUnknown() {
		var makerWelcomeContent environmentGroupRuleSetMakerWelcomeContentModel
		makerWelcomeContentObj.(basetypes.ObjectValue).As(ctx, &makerWelcomeContent, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		hasStatedChanges := true

		rule := environmentGroupRuleSetParameterDto{
			HasStagedChanges: &hasStatedChanges,
			Type:             MAKER_WELCOME_CONTENT,
			ResourceType:     NOT_SPECIFIED,
			Value:            make([]environmentGroupRuleSetValueDto, 0),
		}

		if dto.Parameters == nil {
			dto.Parameters = make([]*environmentGroupRuleSetParameterDto, 0)
		}
		dto.Parameters = append(dto.Parameters, &rule)

		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    MAKER_ONBOARDING_URL,
			Value: makerWelcomeContent.MakerOnboardingUrl.ValueString(),
		})
		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    MAKER_ONBOARDING_MARKDOWN,
			Value: makerWelcomeContent.MakerOnboardingMarkdown.ValueString(),
		})
		rule.Value = append(rule.Value, environmentGroupRuleSetValueDto{
			Id:    MAKER_ONBOARDING_TIMESTAMP,
			Value: time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func convertUsageInsights(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) {
	usageInsightsObj := attrs["usage_insights"]
	if !usageInsightsObj.IsNull() && !usageInsightsObj.IsUnknown() {
		var usageInsights environmentGroupRuleSetUsageInsightsModel
		usageInsightsObj.(basetypes.ObjectValue).As(ctx, &usageInsights, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

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
		var sharingControl environmentGroupRuleSetSharingControlsModel
		sharingControlObj.(basetypes.ObjectValue).As(ctx, &sharingControl, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

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

func convertEnvironmentGroupRuleSetDtoToModel(dto EnvironmentGroupRuleSetValueSetDto) (*environmentGroupRuleSetResourceModel, error) {
	rulesModel, err := convertRulesDtoToModel(dto)
	if err != nil {
		return nil, err
	}

	model := environmentGroupRuleSetResourceModel{
		EnvironmentGroupId: types.StringValue(dto.EnvironmentFilter.Value[0].Id),
		Id:                 types.StringPointerValue(dto.Id),
		Rules:              rulesModel,
	}

	return &model, nil
}

func convertRulesDtoToModel(dto EnvironmentGroupRuleSetValueSetDto) (basetypes.ObjectValue, error) {
	sharingControlType, sharingControlValue, err := convertSharingControlsDtoToModel(getParameterByType(dto, SHARING))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	usageInsightsType, usageInsightsValue, err := convertUsageInsightsDtoToModel(getParameterByType(dto, USAGE_INSIGHTS))
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}
	makerWelcomeContentType, makerWelcomeContentValue, err := convertMakerWelcomeContentDtoToModel(getParameterByType(dto, MAKER_WELCOME_CONTENT))
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

	atr := map[string]attr.Type{
		"sharing_controls":             sharingControlType,
		"usage_insights":               usageInsightsType,
		"maker_welcome_content":        makerWelcomeContentType,
		"solution_checker_enforcement": solutionChekerEnforcementType,
		"backup_retention":             backupRetentionType,
		"ai_generated_descriptions":    aiGeneratedDescType,
		"ai_generative_settings":       aiGenerativeSettingsType,
	}

	if len(dto.Parameters) == 0 {
		return types.ObjectNull(atr), nil
	}

	attrV := map[string]attr.Value{
		"sharing_controls":             sharingControlValue,
		"usage_insights":               usageInsightsValue,
		"maker_welcome_content":        makerWelcomeContentValue,
		"solution_checker_enforcement": solutionChekerEnforcementValue,
		"backup_retention":             backupRetentionValue,
		"ai_generated_descriptions":    aiGeneratedDescValue,
		"ai_generative_settings":       aiGenerativeSettingsValue,
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

	assert.Equal(nil, AI_GENERATIVE_SETTINGS, dto.Type, fmt.Sprintf("Type should be %s", AI_GENERATIVE_SETTINGS))
	assert.Equal(nil, NOT_SPECIFIED, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", NOT_SPECIFIED))

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

	assert.Equal(nil, AI_GENERATED_DESC, dto.Type, fmt.Sprintf("Type should be %s", AI_GENERATED_DESC))
	assert.Equal(nil, APP, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", APP))

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

	assert.Equal(nil, BACKUP_RETENTION, dto.Type, fmt.Sprintf("Type should be %s", BACKUP_RETENTION))
	assert.Equal(nil, NOT_SPECIFIED, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", NOT_SPECIFIED))

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

	assert.Equal(nil, SOLUTION_CHECKER_ENFORCEMENT, dto.Type, fmt.Sprintf("Type should be %s", SOLUTION_CHECKER_ENFORCEMENT))
	assert.Equal(nil, NOT_SPECIFIED, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", NOT_SPECIFIED))

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

func convertMakerWelcomeContentDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"maker_onboarding_url":      types.StringType,
		"maker_onboarding_markdown": types.StringType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	assert.Equal(nil, MAKER_WELCOME_CONTENT, dto.Type, fmt.Sprintf("Type should be %s", MAKER_WELCOME_CONTENT))
	assert.Equal(nil, NOT_SPECIFIED, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", NOT_SPECIFIED))

	makerOnboardingUrl := tryGetRuleValueFromDto(dto.Value, MAKER_ONBOARDING_URL)
	if makerOnboardingUrl == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", MAKER_ONBOARDING_URL)
	}
	makerOnboardingMarkdown := tryGetRuleValueFromDto(dto.Value, MAKER_ONBOARDING_MARKDOWN)
	if makerOnboardingMarkdown == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s value not found in response", MAKER_ONBOARDING_MARKDOWN)
	}
	attrValue := map[string]attr.Value{}
	attrValue["maker_onboarding_url"] = types.StringValue(makerOnboardingUrl.Value)
	attrValue["maker_onboarding_markdown"] = types.StringValue(makerOnboardingMarkdown.Value)

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

func convertUsageInsightsDtoToModel(dto *environmentGroupRuleSetParameterDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"insights_enabled": types.BoolType,
	}

	if dto == nil || len(dto.Value) == 0 {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	assert.Equal(nil, USAGE_INSIGHTS, dto.Type, fmt.Sprintf("Type should be %s", USAGE_INSIGHTS))
	assert.Equal(nil, NOT_SPECIFIED, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", NOT_SPECIFIED))

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

	assert.Equal(nil, SHARING, dto.Type, fmt.Sprintf("Type should be %s", SHARING))
	assert.Equal(nil, APP, dto.ResourceType, fmt.Sprintf("ResourceType should be %s", APP))

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
		assert.Equal(nil, -1.0, maxLimitValue, "Max limit value should be -1, when share mode is 'no limit'")

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

func getParameterByType(params EnvironmentGroupRuleSetValueSetDto, paramType string) *environmentGroupRuleSetParameterDto {
	for paramInx := range params.Parameters {
		if params.Parameters[paramInx].Type == paramType {
			return params.Parameters[paramInx]
		}
	}
	return nil
}
