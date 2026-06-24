// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_based_policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const (
	ADVANCED_CONNECTOR_POLICIES_ONLY_ID    = "AdvancedConnectorPoliciesOnly"
	ADVANCED_CONNECTOR_POLICIES_VERSION    = "1.0"
	ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY = "EnableAdvancedConnectorPoliciesOnly"

	CONTENT_SECURITY_POLICY_ID      = "PowerAppsContentSecurityPolicy"
	CONTENT_SECURITY_POLICY_VERSION = "1.0"

	CSP_IS_ENABLED_KEY                         = "IsContentSecurityPolicyEnabled"
	CSP_IS_ENABLED_FOR_CANVAS_KEY              = "IsContentSecurityPolicyEnabledForCanvas"
	CSP_IS_ENABLED_FOR_CODE_APPS_KEY           = "IsContentSecurityPolicyEnabledForCodeApps"
	CSP_OPTIONS_KEY                            = "ContentSecurityPolicyOptions"
	CSP_REPORT_URI_KEY                         = "ContentSecurityPolicyReportUri"
	CSP_REPORTING_ENDPOINT_KEY                 = "ContentSecurityPolicyReportingEndpoint"
	CSP_CONFIGURATION_KEY                      = "ContentSecurityPolicyConfiguration"
	CSP_CONFIGURATION_FOR_CANVAS_KEY           = "ContentSecurityPolicyConfigurationForCanvas"
	CSP_CONFIGURATION_FOR_CODE_APPS_KEY        = "ContentSecurityPolicyConfigurationForCodeApps"

	CONNECTOR_MANAGEMENT_ID                    = "ConnectorManagement"
	CONNECTOR_MANAGEMENT_VERSION               = "1.0"
	CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY      = "AllowedConnectorList"
	CONNECTOR_API_PREFIX                       = "/providers/Microsoft.PowerApps/apis/"
)

type ruleBasedPolicyDto struct {
	Id           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	TenantId     string       `json:"tenantId,omitempty"`
	LastModified string       `json:"lastModified,omitempty"`
	RuleSetCount int          `json:"ruleSetCount,omitempty"`
	RuleSets     []ruleSetDto `json:"ruleSets,omitempty"`
}

type ruleBasedPolicyRequestDto struct {
	Name     string       `json:"name"`
	RuleSets []ruleSetDto `json:"ruleSets"`
}

type ruleSetDto struct {
	Id      string                 `json:"id"`
	Version string                 `json:"version"`
	Inputs  map[string]interface{} `json:"inputs"`
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
	AssignmentOverrides []interface{} `json:"assignmentOverrides"`
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

func convertModelToDto(ctx context.Context, model environmentGroupRuleBasedPolicyResourceModel) (ruleBasedPolicyRequestDto, error) {
	dto := ruleBasedPolicyRequestDto{
		Name:     fmt.Sprintf("terraform-managed-%s", model.EnvironmentGroupId.ValueString()),
		RuleSets: make([]ruleSetDto, 0),
	}

	if model.RuleSets.IsNull() || model.RuleSets.IsUnknown() {
		return dto, nil
	}

	ruleSetsAttrs := model.RuleSets.Attributes()
	advancedConnectorObj := ruleSetsAttrs["advanced_connector_policies_only"]
	if !advancedConnectorObj.IsNull() && !advancedConnectorObj.IsUnknown() {
		objectValue, ok := advancedConnectorObj.(basetypes.ObjectValue)
		if !ok {
			return dto, fmt.Errorf("expected advanced_connector_policies_only to be of type ObjectValue, got %T", advancedConnectorObj)
		}

		var advancedConnector advancedConnectorPoliciesOnlyModel
		if diags := objectValue.As(ctx, &advancedConnector, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
			return dto, fmt.Errorf("failed to convert advanced_connector_policies_only: %v", diags)
		}

		ruleSet := ruleSetDto{
			Id:      ADVANCED_CONNECTOR_POLICIES_ONLY_ID,
			Version: ADVANCED_CONNECTOR_POLICIES_VERSION,
			Inputs: map[string]interface{}{
				ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY: advancedConnector.Enabled.ValueBool(),
			},
		}
		dto.RuleSets = append(dto.RuleSets, ruleSet)
	}

	cspObj := ruleSetsAttrs["content_security_policy"]
	if !cspObj.IsNull() && !cspObj.IsUnknown() {
		objectValue, ok := cspObj.(basetypes.ObjectValue)
		if !ok {
			return dto, fmt.Errorf("expected content_security_policy to be of type ObjectValue, got %T", cspObj)
		}

		cspRuleSet, err := convertContentSecurityPolicyModelToDto(ctx, objectValue)
		if err != nil {
			return dto, fmt.Errorf("failed to convert content_security_policy: %w", err)
		}
		dto.RuleSets = append(dto.RuleSets, cspRuleSet)
	}

	connMgmtObj := ruleSetsAttrs["advanced_connector_policies"]
	if !connMgmtObj.IsNull() && !connMgmtObj.IsUnknown() {
		objectValue, ok := connMgmtObj.(basetypes.ObjectValue)
		if !ok {
			return dto, fmt.Errorf("expected advanced_connector_policies to be of type ObjectValue, got %T", connMgmtObj)
		}

		connRuleSet, err := convertConnectorManagementModelToDto(ctx, objectValue)
		if err != nil {
			return dto, fmt.Errorf("failed to convert advanced_connector_policies: %w", err)
		}
		dto.RuleSets = append(dto.RuleSets, connRuleSet)
	}

	return dto, nil
}

// managedRuleSetIds maps rule set IDs to their versions for all rules managed by this resource.
var managedRuleSetIds = map[string]string{
	ADVANCED_CONNECTOR_POLICIES_ONLY_ID: ADVANCED_CONNECTOR_POLICIES_VERSION,
	CONTENT_SECURITY_POLICY_ID:          CONTENT_SECURITY_POLICY_VERSION,
	CONNECTOR_MANAGEMENT_ID:             CONNECTOR_MANAGEMENT_VERSION,
}

// mergeRuleSets merges the desired rule sets into the existing policy's rule sets.
// It preserves existing rule sets that are not managed by this resource and replaces/adds managed ones.
func mergeRuleSets(existing []ruleSetDto, desired []ruleSetDto) []ruleSetDto {
	merged := make([]ruleSetDto, 0, len(existing)+len(desired))

	// Keep all existing rule sets that are NOT managed by this resource
	for _, rs := range existing {
		if _, managed := managedRuleSetIds[rs.Id]; !managed {
			merged = append(merged, rs)
		}
	}

	// Add all desired (managed) rule sets
	merged = append(merged, desired...)

	return merged
}

func convertDtoToModel(dto ruleBasedPolicyDto, environmentGroupId string) (*environmentGroupRuleBasedPolicyResourceModel, error) {
	ruleSetsValue, err := convertRuleSetsDtoToModel(dto)
	if err != nil {
		return nil, err
	}

	model := &environmentGroupRuleBasedPolicyResourceModel{
		Id:                 types.StringValue(dto.Id),
		EnvironmentGroupId: types.StringValue(environmentGroupId),
		RuleSets:           ruleSetsValue,
	}

	return model, nil
}

func convertRuleSetsDtoToModel(dto ruleBasedPolicyDto) (basetypes.ObjectValue, error) {
	advancedConnectorType, advancedConnectorValue, err := convertAdvancedConnectorPoliciesDtoToModel(dto.RuleSets)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	cspType, cspValue, err := convertContentSecurityPolicyDtoToModel(dto.RuleSets)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	connMgmtType, connMgmtValue, err := convertConnectorManagementDtoToModel(dto.RuleSets)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	attrTypes := map[string]attr.Type{
		"advanced_connector_policies_only": advancedConnectorType,
		"content_security_policy":          cspType,
		"advanced_connector_policies":      connMgmtType,
	}

	attrValues := map[string]attr.Value{
		"advanced_connector_policies_only": advancedConnectorValue,
		"content_security_policy":          cspValue,
		"advanced_connector_policies":      connMgmtValue,
	}

	return types.ObjectValueMust(attrTypes, attrValues), nil
}

func convertAdvancedConnectorPoliciesDtoToModel(ruleSets []ruleSetDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrType := map[string]attr.Type{
		"enabled": types.BoolType,
	}

	var advancedConnectorRuleSet *ruleSetDto
	for i := range ruleSets {
		if ruleSets[i].Id == ADVANCED_CONNECTOR_POLICIES_ONLY_ID {
			advancedConnectorRuleSet = &ruleSets[i]
			break
		}
	}

	if advancedConnectorRuleSet == nil {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), nil
	}

	enabledValue, ok := advancedConnectorRuleSet.Inputs[ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY]
	if !ok {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s input not found in response", ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY)
	}

	enabled, ok := enabledValue.(bool)
	if !ok {
		return types.ObjectType{AttrTypes: attrType}, types.ObjectNull(attrType), fmt.Errorf("%s input is not a boolean", ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY)
	}

	attrValue := map[string]attr.Value{
		"enabled": types.BoolValue(enabled),
	}

	return types.ObjectType{AttrTypes: attrType}, types.ObjectValueMust(attrType, attrValue), nil
}

// CSP configuration DTO types for JSON serialization
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

// cspConfigurationAttrTypes returns the Terraform attribute types for a CSP configuration block.
func cspConfigurationAttrTypes(includeStrictCsp bool) map[string]attr.Type {
	attrTypes := map[string]attr.Type{
		"img_src":        types.ListType{ElemType: types.StringType},
		"style_src":      types.ListType{ElemType: types.StringType},
		"form_action":    types.ListType{ElemType: types.StringType},
		"frame_src":      types.ListType{ElemType: types.StringType},
		"connect_src":    types.ListType{ElemType: types.StringType},
		"font_src":       types.ListType{ElemType: types.StringType},
		"script_src":     types.ListType{ElemType: types.StringType},
		"frame_ancestor": types.ListType{ElemType: types.StringType},
	}
	if includeStrictCsp {
		attrTypes["strict_csp"] = types.BoolType
	}
	return attrTypes
}

// contentSecurityPolicyAttrTypes returns the Terraform attribute types for the CSP rule set.
func contentSecurityPolicyAttrTypes() map[string]attr.Type {
	configType := types.ObjectType{AttrTypes: cspConfigurationAttrTypes(true)}
	codeAppsConfigType := types.ObjectType{AttrTypes: cspConfigurationAttrTypes(false)}
	return map[string]attr.Type{
		"enabled":                    types.BoolType,
		"enabled_for_canvas":         types.BoolType,
		"enabled_for_code_apps":      types.BoolType,
		"report_uri":                 types.StringType,
		"reporting_endpoint":         types.StringType,
		"configuration":              configType,
		"configuration_for_canvas":   configType,
		"configuration_for_code_apps": codeAppsConfigType,
	}
}

func convertContentSecurityPolicyModelToDto(ctx context.Context, objectValue basetypes.ObjectValue) (ruleSetDto, error) {
	var csp contentSecurityPolicyModel
	if diags := objectValue.As(ctx, &csp, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
		return ruleSetDto{}, fmt.Errorf("failed to convert content_security_policy: %v", diags)
	}

	// Calculate ContentSecurityPolicyOptions from strict_csp booleans
	modelDrivenStrict := getStrictCspFromConfig(ctx, csp.Configuration)
	canvasStrict := getStrictCspFromConfig(ctx, csp.ConfigurationForCanvas)
	options := calculateCspOptions(modelDrivenStrict, canvasStrict)

	inputs := map[string]interface{}{
		CSP_IS_ENABLED_KEY:               csp.Enabled.ValueBool(),
		CSP_IS_ENABLED_FOR_CANVAS_KEY:    csp.EnabledForCanvas.ValueBool(),
		CSP_IS_ENABLED_FOR_CODE_APPS_KEY: csp.EnabledForCodeApps.ValueBool(),
		CSP_OPTIONS_KEY:                  options,
	}

	if !csp.ReportUri.IsNull() && !csp.ReportUri.IsUnknown() {
		inputs[CSP_REPORT_URI_KEY] = csp.ReportUri.ValueString()
	} else {
		inputs[CSP_REPORT_URI_KEY] = nil
	}

	if !csp.ReportingEndpoint.IsNull() && !csp.ReportingEndpoint.IsUnknown() {
		inputs[CSP_REPORTING_ENDPOINT_KEY] = csp.ReportingEndpoint.ValueString()
	} else {
		inputs[CSP_REPORTING_ENDPOINT_KEY] = nil
	}

	configJson, err := convertCspConfigModelToJson(ctx, csp.Configuration)
	if err != nil {
		return ruleSetDto{}, fmt.Errorf("failed to convert configuration: %w", err)
	}
	inputs[CSP_CONFIGURATION_KEY] = configJson

	canvasJson, err := convertCspConfigModelToJson(ctx, csp.ConfigurationForCanvas)
	if err != nil {
		return ruleSetDto{}, fmt.Errorf("failed to convert configuration_for_canvas: %w", err)
	}
	inputs[CSP_CONFIGURATION_FOR_CANVAS_KEY] = canvasJson

	codeAppsJson, err := convertCspConfigModelToJson(ctx, csp.ConfigurationForCodeApps)
	if err != nil {
		return ruleSetDto{}, fmt.Errorf("failed to convert configuration_for_code_apps: %w", err)
	}
	inputs[CSP_CONFIGURATION_FOR_CODE_APPS_KEY] = codeAppsJson

	return ruleSetDto{
		Id:      CONTENT_SECURITY_POLICY_ID,
		Version: CONTENT_SECURITY_POLICY_VERSION,
		Inputs:  inputs,
	}, nil
}

// getStrictCspFromConfig extracts the strict_csp boolean from a configuration object.
func getStrictCspFromConfig(_ context.Context, configObj types.Object) bool {
	if configObj.IsNull() || configObj.IsUnknown() {
		return false
	}
	attrs := configObj.Attributes()
	v, ok := attrs["strict_csp"]
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(types.Bool); ok && !b.IsNull() && !b.IsUnknown() {
		return b.ValueBool()
	}
	return false
}

// calculateCspOptions computes the ContentSecurityPolicyOptions integer from strict_csp flags.
// model-driven strict only = 1, canvas strict only = 8, both = 9, neither = 0
func calculateCspOptions(modelDrivenStrict, canvasStrict bool) int64 {
	var options int64
	if modelDrivenStrict {
		options += 1
	}
	if canvasStrict {
		options += 8
	}
	return options
}

func convertCspConfigModelToJson(_ context.Context, configObj types.Object) (string, error) {
	if configObj.IsNull() || configObj.IsUnknown() {
		return "{}", nil
	}

	attrs := configObj.Attributes()
	dto := cspConfigurationDto{}

	getList := func(key string) types.List {
		v, ok := attrs[key]
		if !ok || v == nil {
			return types.ListNull(types.StringType)
		}
		if l, ok := v.(types.List); ok {
			return l
		}
		return types.ListNull(types.StringType)
	}

	setDirective := func(list types.List) *cspDirectiveDto {
		if list.IsNull() || list.IsUnknown() {
			return nil
		}
		sources := make([]cspSourceDto, 0)
		for _, elem := range list.Elements() {
			if strVal, ok := elem.(types.String); ok {
				sources = append(sources, cspSourceDto{Source: strVal.ValueString()})
			}
		}
		return &cspDirectiveDto{Sources: sources}
	}

	dto.ImgSrc = setDirective(getList("img_src"))
	dto.StyleSrc = setDirective(getList("style_src"))
	dto.FormAction = setDirective(getList("form_action"))
	dto.FrameSrc = setDirective(getList("frame_src"))
	dto.ConnectSrc = setDirective(getList("connect_src"))
	dto.FontSrc = setDirective(getList("font_src"))
	dto.ScriptSrc = setDirective(getList("script_src"))
	dto.FrameAncestor = setDirective(getList("frame_ancestor"))

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CSP configuration to JSON: %w", err)
	}

	return string(jsonBytes), nil
}

func convertContentSecurityPolicyDtoToModel(ruleSets []ruleSetDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	cspAttrTypes := contentSecurityPolicyAttrTypes()
	objType := types.ObjectType{AttrTypes: cspAttrTypes}

	var cspRuleSet *ruleSetDto
	for i := range ruleSets {
		if ruleSets[i].Id == CONTENT_SECURITY_POLICY_ID {
			cspRuleSet = &ruleSets[i]
			break
		}
	}

	if cspRuleSet == nil {
		return objType, types.ObjectNull(cspAttrTypes), nil
	}

	inputs := cspRuleSet.Inputs

	getBool := func(key string) types.Bool {
		v, ok := inputs[key]
		if !ok || v == nil {
			return types.BoolValue(false)
		}
		if b, ok := v.(bool); ok {
			return types.BoolValue(b)
		}
		return types.BoolValue(false)
	}

	getInt64 := func(key string) types.Int64 {
		v, ok := inputs[key]
		if !ok || v == nil {
			return types.Int64Value(0)
		}
		switch n := v.(type) {
		case float64:
			return types.Int64Value(int64(n))
		case int64:
			return types.Int64Value(n)
		}
		return types.Int64Value(0)
	}

	getNullableString := func(key string) types.String {
		v, ok := inputs[key]
		if !ok || v == nil {
			return types.StringNull()
		}
		if s, ok := v.(string); ok {
			return types.StringValue(s)
		}
		return types.StringNull()
	}

	getConfigValue := func(key string, includeStrictCsp bool, strictCspValue bool) (basetypes.ObjectValue, error) {
		v, ok := inputs[key]
		if !ok || v == nil {
			return types.ObjectNull(cspConfigurationAttrTypes(includeStrictCsp)), nil
		}
		jsonStr, ok := v.(string)
		if !ok {
			return types.ObjectNull(cspConfigurationAttrTypes(includeStrictCsp)), nil
		}
		return convertCspConfigJsonToModel(jsonStr, includeStrictCsp, strictCspValue)
	}

	// Derive strict_csp booleans from ContentSecurityPolicyOptions
	optionsVal := getInt64(CSP_OPTIONS_KEY)
	optionsInt := optionsVal.ValueInt64()
	modelDrivenStrict := (optionsInt & 1) != 0
	canvasStrict := (optionsInt & 8) != 0

	configValue, err := getConfigValue(CSP_CONFIGURATION_KEY, true, modelDrivenStrict)
	if err != nil {
		return objType, types.ObjectNull(cspAttrTypes), fmt.Errorf("failed to parse %s: %w", CSP_CONFIGURATION_KEY, err)
	}

	canvasValue, err := getConfigValue(CSP_CONFIGURATION_FOR_CANVAS_KEY, true, canvasStrict)
	if err != nil {
		return objType, types.ObjectNull(cspAttrTypes), fmt.Errorf("failed to parse %s: %w", CSP_CONFIGURATION_FOR_CANVAS_KEY, err)
	}

	codeAppsValue, err := getConfigValue(CSP_CONFIGURATION_FOR_CODE_APPS_KEY, false, false)
	if err != nil {
		return objType, types.ObjectNull(cspAttrTypes), fmt.Errorf("failed to parse %s: %w", CSP_CONFIGURATION_FOR_CODE_APPS_KEY, err)
	}

	attrValues := map[string]attr.Value{
		"enabled":                    getBool(CSP_IS_ENABLED_KEY),
		"enabled_for_canvas":         getBool(CSP_IS_ENABLED_FOR_CANVAS_KEY),
		"enabled_for_code_apps":      getBool(CSP_IS_ENABLED_FOR_CODE_APPS_KEY),
		"report_uri":                 getNullableString(CSP_REPORT_URI_KEY),
		"reporting_endpoint":         getNullableString(CSP_REPORTING_ENDPOINT_KEY),
		"configuration":              configValue,
		"configuration_for_canvas":   canvasValue,
		"configuration_for_code_apps": codeAppsValue,
	}

	return objType, types.ObjectValueMust(cspAttrTypes, attrValues), nil
}

func convertCspConfigJsonToModel(jsonStr string, includeStrictCsp bool, strictCspValue bool) (basetypes.ObjectValue, error) {
	configAttrTypes := cspConfigurationAttrTypes(includeStrictCsp)

	if jsonStr == "" || jsonStr == "{}" {
		return types.ObjectNull(configAttrTypes), nil
	}

	var dto cspConfigurationDto
	if err := json.Unmarshal([]byte(jsonStr), &dto); err != nil {
		return types.ObjectNull(configAttrTypes), fmt.Errorf("failed to unmarshal CSP configuration JSON: %w", err)
	}

	directiveToList := func(d *cspDirectiveDto) types.List {
		if d == nil || len(d.Sources) == 0 {
			return types.ListNull(types.StringType)
		}
		elems := make([]attr.Value, len(d.Sources))
		for i, s := range d.Sources {
			elems[i] = types.StringValue(s.Source)
		}
		list, _ := types.ListValue(types.StringType, elems)
		return list
	}

	attrValues := map[string]attr.Value{
		"img_src":        directiveToList(dto.ImgSrc),
		"style_src":      directiveToList(dto.StyleSrc),
		"form_action":    directiveToList(dto.FormAction),
		"frame_src":      directiveToList(dto.FrameSrc),
		"connect_src":    directiveToList(dto.ConnectSrc),
		"font_src":       directiveToList(dto.FontSrc),
		"script_src":     directiveToList(dto.ScriptSrc),
		"frame_ancestor": directiveToList(dto.FrameAncestor),
	}

	if includeStrictCsp {
		attrValues["strict_csp"] = types.BoolValue(strictCspValue)
	}

	return types.ObjectValueMust(configAttrTypes, attrValues), nil
}

// --- ConnectorManagement (advanced_connector_policies) conversion functions ---

const (
	actionsModeAllAllowed  = "all_allowed"
	actionsModeSomeAllowed = "some_allowed"
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

func convertConnectorManagementModelToDto(ctx context.Context, objectValue basetypes.ObjectValue) (ruleSetDto, error) {
	attrs := objectValue.Attributes()
	connListAttr, ok := attrs["allowed_connectors"]
	if !ok {
		return ruleSetDto{}, fmt.Errorf("missing allowed_connectors attribute")
	}

	connList, ok := connListAttr.(types.List)
	if !ok || connList.IsNull() || connList.IsUnknown() {
		return ruleSetDto{
			Id:      CONNECTOR_MANAGEMENT_ID,
			Version: CONNECTOR_MANAGEMENT_VERSION,
			Inputs:  map[string]interface{}{CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY: []allowedConnectorDto{}},
		}, nil
	}

	dtoList := make([]allowedConnectorDto, 0, len(connList.Elements()))
	for _, elem := range connList.Elements() {
		objVal, ok := elem.(types.Object)
		if !ok {
			continue
		}
		elemAttrs := objVal.Attributes()

		connectorId := ""
		if v, ok := elemAttrs["connector_id"].(types.String); ok {
			connectorId = v.ValueString()
		}

		actionsMode := actionsModeAllAllowed
		if v, ok := elemAttrs["actions_mode"].(types.String); ok && !v.IsNull() && !v.IsUnknown() {
			actionsMode = v.ValueString()
		}

		apiActionsMode := ""
		switch actionsMode {
		case actionsModeAllAllowed:
			apiActionsMode = "AllAllowed"
		case actionsModeSomeAllowed:
			apiActionsMode = "SomeAllowed"
		default:
			return ruleSetDto{}, fmt.Errorf("invalid actions_mode %q: must be %q or %q", actionsMode, actionsModeAllAllowed, actionsModeSomeAllowed)
		}

		allowedConnector := connectorId
		if len(connectorId) < len(CONNECTOR_API_PREFIX) || connectorId[:len(CONNECTOR_API_PREFIX)] != CONNECTOR_API_PREFIX {
			allowedConnector = CONNECTOR_API_PREFIX + connectorId
		}

		entry := allowedConnectorDto{
			AllowedActionsMode:         apiActionsMode,
			AllowedConnectionTypesMode: "AllAllowed",
			AllowedConnector:           allowedConnector,
		}

		if actionsMode == actionsModeSomeAllowed {
			if actionsAttr, ok := elemAttrs["allowed_actions"].(types.List); ok && !actionsAttr.IsNull() && !actionsAttr.IsUnknown() {
				actions := make([]string, 0, len(actionsAttr.Elements()))
				for _, a := range actionsAttr.Elements() {
					if s, ok := a.(types.String); ok {
						actions = append(actions, s.ValueString())
					}
				}
				entry.AllowedActions = actions
			}
		}

		dtoList = append(dtoList, entry)
	}

	return ruleSetDto{
		Id:      CONNECTOR_MANAGEMENT_ID,
		Version: CONNECTOR_MANAGEMENT_VERSION,
		Inputs:  map[string]interface{}{CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY: dtoList},
	}, nil
}

func convertConnectorManagementDtoToModel(ruleSets []ruleSetDto) (basetypes.ObjectType, basetypes.ObjectValue, error) {
	attrTypes := advancedConnectorPoliciesAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}

	var cmRuleSet *ruleSetDto
	for i := range ruleSets {
		if ruleSets[i].Id == CONNECTOR_MANAGEMENT_ID {
			cmRuleSet = &ruleSets[i]
			break
		}
	}

	if cmRuleSet == nil {
		return objType, types.ObjectNull(attrTypes), nil
	}

	rawList, ok := cmRuleSet.Inputs[CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY]
	if !ok || rawList == nil {
		emptyList, _ := types.ListValue(types.ObjectType{AttrTypes: allowedConnectorAttrTypes()}, []attr.Value{})
		return objType, types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"allowed_connectors": emptyList,
		}), nil
	}

	// rawList is []interface{} from JSON unmarshalling
	rawSlice, ok := rawList.([]interface{})
	if !ok {
		emptyList, _ := types.ListValue(types.ObjectType{AttrTypes: allowedConnectorAttrTypes()}, []attr.Value{})
		return objType, types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"allowed_connectors": emptyList,
		}), nil
	}

	connectorElemType := types.ObjectType{AttrTypes: allowedConnectorAttrTypes()}
	connectorElems := make([]attr.Value, 0, len(rawSlice))

	for _, raw := range rawSlice {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		connectorPath := ""
		if v, ok := entry["AllowedConnector"].(string); ok {
			connectorPath = v
		}
		// Strip the API prefix to get the short connector ID
		connectorId := connectorPath
		if len(connectorPath) >= len(CONNECTOR_API_PREFIX) && connectorPath[:len(CONNECTOR_API_PREFIX)] == CONNECTOR_API_PREFIX {
			connectorId = connectorPath[len(CONNECTOR_API_PREFIX):]
		}

		apiActionsMode := "AllAllowed"
		if v, ok := entry["AllowedActionsMode"].(string); ok {
			apiActionsMode = v
		}
		actionsMode := actionsModeAllAllowed
		if apiActionsMode == "SomeAllowed" {
			actionsMode = actionsModeSomeAllowed
		}

		var allowedActionsList types.List
		if actionsMode == actionsModeSomeAllowed {
			if rawActions, ok := entry["AllowedActions"].([]interface{}); ok {
				actionElems := make([]attr.Value, 0, len(rawActions))
				for _, a := range rawActions {
					if s, ok := a.(string); ok {
						actionElems = append(actionElems, types.StringValue(s))
					}
				}
				allowedActionsList, _ = types.ListValue(types.StringType, actionElems)
			} else {
				allowedActionsList, _ = types.ListValue(types.StringType, []attr.Value{})
			}
		} else {
			allowedActionsList = types.ListNull(types.StringType)
		}

		elemValue := types.ObjectValueMust(allowedConnectorAttrTypes(), map[string]attr.Value{
			"connector_id":    types.StringValue(connectorId),
			"actions_mode":    types.StringValue(actionsMode),
			"allowed_actions": allowedActionsList,
		})
		connectorElems = append(connectorElems, elemValue)
	}

	connectorsList, _ := types.ListValue(connectorElemType, connectorElems)

	return objType, types.ObjectValueMust(attrTypes, map[string]attr.Value{
		"allowed_connectors": connectorsList,
	}), nil
}
