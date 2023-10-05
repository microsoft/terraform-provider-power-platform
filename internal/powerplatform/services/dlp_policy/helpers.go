package powerplatform

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func covertDlpPolicyToPolicyModelDto(policy DlpPolicyDto) (*DlpPolicyModelDto, error) {

	policyModel := DlpPolicyModelDto{
		Name:                                 policy.PolicyDefinition.Name,
		DisplayName:                          policy.PolicyDefinition.DisplayName,
		EnvironmentType:                      policy.PolicyDefinition.EnvironmentType,
		Environments:                         policy.PolicyDefinition.Environments,
		ETag:                                 policy.PolicyDefinition.ETag,
		CreatedBy:                            policy.PolicyDefinition.CreatedBy.DisplayName,
		CreatedTime:                          policy.PolicyDefinition.CreatedTime,
		LastModifiedBy:                       policy.PolicyDefinition.LastModifiedBy.DisplayName,
		LastModifiedTime:                     policy.PolicyDefinition.LastModifiedTime,
		DefaultConnectorsClassification:      policy.PolicyDefinition.DefaultConnectorsClassification,
		ConnectorConfigurationsDefinition:    DlpConnectorConfigurationsDefinitionDto{},
		CustomConnectorUrlPatternsDefinition: []DlpConnectorUrlPatternsDefinitionDto{},
		ConnectorGroups:                      []DlpConnectorGroupsModelDto{},
	}

	for _, connGroup := range policy.PolicyDefinition.ConnectorGroups {
		connGroupModel := DlpConnectorGroupsModelDto{
			Classification: connGroup.Classification,
			Connectors:     []DlpConnectorModelDto{},
		}
		for _, connector := range connGroup.Connectors {
			nameSplit := strings.Split(connector.Id, "/")
			m := DlpConnectorModelDto{
				Id:   connector.Id,
				Name: nameSplit[len(nameSplit)-1],
				Type: connector.Type,
			}
			for _, connectorActionConfigurations := range policy.ConnectorConfigurationsDefinition.ConnectorActionConfigurations {
				if connectorActionConfigurations.ConnectorId == connector.Id {
					m.DefaultActionRuleBehavior = connectorActionConfigurations.DefaultConnectorActionRuleBehavior
					m.ActionRules = connectorActionConfigurations.ActionRules
				}
			}
			for _, endpointConfigurations := range policy.ConnectorConfigurationsDefinition.EndpointConfigurations {
				if endpointConfigurations.ConnectorId == connector.Id {
					m.EndpointRules = endpointConfigurations.EndpointRules
				}
			}
			connGroupModel.Connectors = append(connGroupModel.Connectors, m)

		}
		policyModel.ConnectorGroups = append(policyModel.ConnectorGroups, connGroupModel)
	}

	for _, rule := range policy.CustomConnectorUrlPatternsDefinition.Rules {
		policyModel.CustomConnectorUrlPatternsDefinition = append(policyModel.CustomConnectorUrlPatternsDefinition, DlpConnectorUrlPatternsDefinitionDto{
			Rules: append([]DlpConnectorUrlPatternsRuleDto{}, rule),
		})
	}

	return &policyModel, nil
}

func convertConnectorRuleClassificationValues(value string) string {
	if value == "Business" {
		return "General"
	} else if value == "NonBusiness" {
		return "Confidential"
	} else if value == "General" {
		return "Business"
	} else if value == "Confidential" {
		return "NonBusiness"
	} else {
		return value
	}
}

func convertToAttrValueConnectorsGroup(classification string, connectorsGroup []DlpConnectorGroupsModelDto) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
		}
	}
	return types.SetValueMust(connectorSetObjectType, []attr.Value{})
}

func convertToAttrValueCustomConnectorUrlPatternsDefinition(urlPatterns []DlpConnectorUrlPatternsDefinitionDto) basetypes.SetValue {
	var connUrlPattern []attr.Value
	for _, connectorUrlPattern := range urlPatterns {
		for _, rules := range connectorUrlPattern.Rules {
			connUrlPattern = append(connUrlPattern, types.ObjectValueMust(
				map[string]attr.Type{
					"order":            types.Int64Type,
					"host_url_pattern": types.StringType,
					"data_group":       types.StringType,
				},
				map[string]attr.Value{
					"order":            types.Int64Value(rules.Order),
					"host_url_pattern": types.StringValue(rules.Pattern),
					"data_group":       types.StringValue(convertConnectorRuleClassificationValues(rules.ConnectorRuleClassification)),
				},
			))
		}
	}
	if len(urlPatterns) == 0 {
		return types.SetValueMust(customConnectorPatternSetObjectType, []attr.Value{})
	} else {
		return types.SetValueMust(customConnectorPatternSetObjectType, connUrlPattern)
	}
}

func convertToAttrValueEnvironments(environments []DlpEnvironmentDto) basetypes.SetValue {
	var env []attr.Value
	for _, environment := range environments {
		env = append(env, types.ObjectValueMust(
			map[string]attr.Type{
				"name": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue(environment.Name),
			},
		))
	}

	if len(environments) == 0 {
		return types.SetValueMust(environmentSetObjectType, []attr.Value{})
	} else {
		return types.SetValueMust(environmentSetObjectType, env)
	}
}

func convertToAttrValueConnectors(connectorsGroup DlpConnectorGroupsModelDto, connectors []attr.Value) []attr.Value {
	for _, connector := range connectorsGroup.Connectors {
		connectors = append(connectors, types.ObjectValueMust(
			map[string]attr.Type{
				//"name":                         types.StringType,
				"id":                           types.StringType,
				"default_action_rule_behavior": types.StringType,
				"action_rules": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"action_id": types.StringType,
							"behavior":  types.StringType,
						},
					}},
				"endpoint_rules": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"order":    types.Int64Type,
							"behavior": types.StringType,
							"endpoint": types.StringType,
						},
					}},
			},
			map[string]attr.Value{
				//"name":                         types.StringValue(connector.Name),
				"id":                           types.StringValue(connector.Id),
				"default_action_rule_behavior": types.StringValue(connector.DefaultActionRuleBehavior),
				"action_rules":                 types.ListValueMust(actionRuleListObjectType, convertToAtrValueActionRule(connector)),
				"endpoint_rules":               types.ListValueMust(endpointRuleListObjectType, convertToAtrValueEndpointRule(connector)),
			},
		))
	}
	return connectors
}

func convertToDlpConnectorGroup(ctx context.Context, diag diag.Diagnostics, classification string, connectorsAttr basetypes.SetValue) DlpConnectorGroupsModelDto {
	var connectors []DataLossPreventionPolicyResourceConnectorModel
	err := connectorsAttr.ElementsAs(ctx, &connectors, true)
	if err != nil {
		diag.AddError("Client error when converting DlpConnectorGroups", "")
	}

	connectorGroup := DlpConnectorGroupsModelDto{
		Classification: classification,
		Connectors:     make([]DlpConnectorModelDto, 0),
	}

	for _, connector := range connectors {
		defaultAction := "Allow"

		if connector.DefaultActionRuleBehavior.ValueString() != "" {
			defaultAction = connector.DefaultActionRuleBehavior.ValueString()
		}

		connectorGroup.Connectors = append(connectorGroup.Connectors, DlpConnectorModelDto{
			Id:   connector.Id.ValueString(),
			Type: "Microsoft.PowerApps/apis",

			DefaultActionRuleBehavior: defaultAction,
			ActionRules:               convertToDlpActionRule(connector),
			EndpointRules:             convertToDlpEndpointRule(connector),
		})
	}
	return connectorGroup
}

func convertToDlpEnvironment(ctx context.Context, diag diag.Diagnostics, environmentsAttr basetypes.SetValue) []DlpEnvironmentDto {
	var envs []DataLossPreventionPolicyResourceEnvironmentsModel
	err := environmentsAttr.ElementsAs(ctx, &envs, true)
	if err != nil {
		diag.AddError("Client error when converting DlpEnvironment", "")
	}

	environments := make([]DlpEnvironmentDto, 0)
	for _, environment := range envs {
		environments = append(environments, DlpEnvironmentDto{
			Name: environment.Name.ValueString(),
			Id:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/" + environment.Name.ValueString(),
			Type: "Microsoft.BusinessAppPlatform/scopes/environments",
		})
	}
	return environments
}

func convertToDlpCustomConnectorUrlPatternsDefinition(ctx context.Context, diag diag.Diagnostics, connectorPatternsAttr basetypes.SetValue) []DlpConnectorUrlPatternsDefinitionDto {
	var customConnectorsPatterns []DataLossPreventionPolicyResourceCustomConnectorPattern
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diag.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
	}

	customConnectorUrlPatternsDefinition := make([]DlpConnectorUrlPatternsDefinitionDto, 0)
	for _, customConnectorPattern := range customConnectorsPatterns {
		urlPattern := DlpConnectorUrlPatternsDefinitionDto{
			Rules: []DlpConnectorUrlPatternsRuleDto{},
		}
		urlPattern.Rules = append(urlPattern.Rules, DlpConnectorUrlPatternsRuleDto{
			Order:                       customConnectorPattern.Order.ValueInt64(),
			ConnectorRuleClassification: convertConnectorRuleClassificationValues(customConnectorPattern.DataGroup.ValueString()),
			Pattern:                     customConnectorPattern.HostUrlPattern.ValueString(),
		})
		customConnectorUrlPatternsDefinition = append(customConnectorUrlPatternsDefinition, urlPattern)
	}
	return customConnectorUrlPatternsDefinition
}

func convertToDlpActionRule(connector DataLossPreventionPolicyResourceConnectorModel) []DlpActionRuleDto {
	var actionRules []DlpActionRuleDto
	for _, actionRule := range connector.ActionRules {
		actionRules = append(actionRules, DlpActionRuleDto{
			ActionId: actionRule.ActionId.ValueString(),
			Behavior: actionRule.Behavior.ValueString(),
		})
	}
	return actionRules
}

func convertToDlpEndpointRule(connector DataLossPreventionPolicyResourceConnectorModel) []DlpEndpointRuleDto {
	var endpointRules []DlpEndpointRuleDto
	for _, endpointRule := range connector.EndpointRules {
		endpointRules = append(endpointRules, DlpEndpointRuleDto{
			Order:    endpointRule.Order.ValueInt64(),
			Behavior: endpointRule.Behavior.ValueString(),
			Endpoint: endpointRule.Endpoint.ValueString(),
		})
	}
	return endpointRules
}

func convertToAtrValueActionRule(connector DlpConnectorModelDto) []attr.Value {
	var actionRules []attr.Value
	for _, actionRule := range connector.ActionRules {
		actionRules = append(actionRules, types.ObjectValueMust(
			map[string]attr.Type{
				"action_id": types.StringType,
				"behavior":  types.StringType,
			},
			map[string]attr.Value{
				"action_id": types.StringValue(actionRule.ActionId),
				"behavior":  types.StringValue(actionRule.Behavior),
			},
		))
	}
	return actionRules
}

func convertToAtrValueEndpointRule(connector DlpConnectorModelDto) []attr.Value {
	var endpointRules []attr.Value
	for _, endpointRule := range connector.EndpointRules {
		endpointRules = append(endpointRules, types.ObjectValueMust(
			map[string]attr.Type{
				"order":    types.Int64Type,
				"behavior": types.StringType,
				"endpoint": types.StringType,
			},
			map[string]attr.Value{
				"order":    types.Int64Value(endpointRule.Order),
				"behavior": types.StringValue(endpointRule.Behavior),
				"endpoint": types.StringValue(endpointRule.Endpoint),
			},
		))
	}
	return endpointRules
}
