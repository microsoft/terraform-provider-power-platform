// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func covertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	policyModel := dlpPolicyModelDto{
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
		ConnectorConfigurationsDefinition:    dlpConnectorConfigurationsDefinitionDto{},
		CustomConnectorUrlPatternsDefinition: []dlpConnectorUrlPatternsDefinitionDto{},
		ConnectorGroups:                      []dlpConnectorGroupsModelDto{},
	}

	for _, connGroup := range policy.PolicyDefinition.ConnectorGroups {
		connGroupModel := dlpConnectorGroupsModelDto{
			Classification: connGroup.Classification,
			Connectors:     []dlpConnectorModelDto{},
		}
		for _, connector := range connGroup.Connectors {
			nameSplit := strings.Split(connector.Id, "/")
			m := dlpConnectorModelDto{
				Id:   connector.Id,
				Name: nameSplit[len(nameSplit)-1],
				Type: connector.Type,
			}
			if policy.ConnectorConfigurationsDefinition != nil {
				if policy.ConnectorConfigurationsDefinition.ConnectorActionConfigurations != nil {
					for _, connectorActionConfigurations := range policy.ConnectorConfigurationsDefinition.ConnectorActionConfigurations {
						if connectorActionConfigurations.ConnectorId == connector.Id {
							m.DefaultActionRuleBehavior = connectorActionConfigurations.DefaultConnectorActionRuleBehavior
							m.ActionRules = connectorActionConfigurations.ActionRules
						}
					}
				}
				if policy.ConnectorConfigurationsDefinition.EndpointConfigurations != nil {
					for _, endpointConfigurations := range policy.ConnectorConfigurationsDefinition.EndpointConfigurations {
						if endpointConfigurations.ConnectorId == connector.Id {
							m.EndpointRules = endpointConfigurations.EndpointRules
						}
					}
				}
			}
			connGroupModel.Connectors = append(connGroupModel.Connectors, m)
		}
		policyModel.ConnectorGroups = append(policyModel.ConnectorGroups, connGroupModel)
	}

	for _, rule := range policy.CustomConnectorUrlPatternsDefinition.Rules {
		policyModel.CustomConnectorUrlPatternsDefinition = append(policyModel.CustomConnectorUrlPatternsDefinition, dlpConnectorUrlPatternsDefinitionDto{
			Rules: append([]dlpConnectorUrlPatternsRuleDto{}, rule),
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
	}
	return value
}

func convertToAttrValueConnectorsGroup(classification string, connectorsGroup []dlpConnectorGroupsModelDto) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
		}
	}
	return types.SetValueMust(connectorSetObjectType, []attr.Value{})
}

func convertToAttrValueCustomConnectorUrlPatternsDefinition(urlPatterns []dlpConnectorUrlPatternsDefinitionDto) basetypes.SetValue {
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
	}
	return types.SetValueMust(customConnectorPatternSetObjectType, connUrlPattern)
}

func convertToAttrValueEnvironments(environments []dlpEnvironmentDto) []string {
	if len(environments) == 0 {
		return []string{}
	}
	var env []string
	for _, environment := range environments {
		env = append(env, environment.Name)
	}
	return env
}

func convertToAttrValueConnectors(connectorsGroup dlpConnectorGroupsModelDto, connectors []attr.Value) []attr.Value {
	for _, connector := range connectorsGroup.Connectors {
		connectors = append(connectors, types.ObjectValueMust(
			map[string]attr.Type{
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
				// "name":                         types.StringValue(connector.Name),
				"id":                           types.StringValue(connector.Id),
				"default_action_rule_behavior": types.StringValue(connector.DefaultActionRuleBehavior),
				"action_rules":                 types.ListValueMust(actionRuleListObjectType, convertToAtrValueActionRule(connector)),
				"endpoint_rules":               types.ListValueMust(endpointRuleListObjectType, convertToAtrValueEndpointRule(connector)),
			},
		))
	}
	return connectors
}

func getConnectorGroup(ctx context.Context, connectorsAttr basetypes.SetValue) (*dlpConnectorGroupsModelDto, error) {
	var connectors []dataLossPreventionPolicyResourceConnectorModel
	err := connectorsAttr.ElementsAs(ctx, &connectors, true)
	if err != nil {
		return nil, fmt.Errorf("error converting elements: %v", err)
	}

	connectorGroup := dlpConnectorGroupsModelDto{

		Connectors: make([]dlpConnectorModelDto, 0),
	}

	for _, connector := range connectors {
		connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
			Id:                        connector.Id.ValueString(),
			Type:                      "Microsoft.PowerApps/apis",
			DefaultActionRuleBehavior: connector.DefaultActionRuleBehavior.ValueString(),
			ActionRules:               convertToDlpActionRule(connector),
			EndpointRules:             convertToDlpEndpointRule(connector),
		})
	}
	return &connectorGroup, nil
}

func convertToDlpConnectorGroup(ctx context.Context, diags diag.Diagnostics, classification string, connectorsAttr basetypes.SetValue) dlpConnectorGroupsModelDto {
	var connectors []dataLossPreventionPolicyResourceConnectorModel
	err := connectorsAttr.ElementsAs(ctx, &connectors, true)
	if err != nil {
		diags.AddError("Client error when converting DlpConnectorGroups", "")
	}

	connectorGroup := dlpConnectorGroupsModelDto{
		Classification: classification,
		Connectors:     make([]dlpConnectorModelDto, 0),
	}

	for _, connector := range connectors {
		defaultAction := "Allow"

		if connector.DefaultActionRuleBehavior.ValueString() != "" {
			defaultAction = connector.DefaultActionRuleBehavior.ValueString()
		}

		connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
			Id:   connector.Id.ValueString(),
			Type: "Microsoft.PowerApps/apis",

			DefaultActionRuleBehavior: defaultAction,
			ActionRules:               convertToDlpActionRule(connector),
			EndpointRules:             convertToDlpEndpointRule(connector),
		})
	}
	return connectorGroup
}

func convertToDlpEnvironment(environmentsInPolicy []string) []dlpEnvironmentDto {
	environments := make([]dlpEnvironmentDto, 0)
	for _, environment := range environmentsInPolicy {
		environments = append(environments, dlpEnvironmentDto{
			Name: environment,
			Id:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/" + environment,
			Type: "Microsoft.BusinessAppPlatform/scopes/environments",
		})
	}
	return environments
}

func convertToDlpCustomConnectorUrlPatternsDefinition(ctx context.Context, diags diag.Diagnostics, connectorPatternsAttr basetypes.SetValue) []dlpConnectorUrlPatternsDefinitionDto {
	var customConnectorsPatterns []dataLossPreventionPolicyResourceCustomConnectorPattern
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diags.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
	}

	customConnectorUrlPatternsDefinition := make([]dlpConnectorUrlPatternsDefinitionDto, 0)
	for _, customConnectorPattern := range customConnectorsPatterns {
		urlPattern := dlpConnectorUrlPatternsDefinitionDto{
			Rules: []dlpConnectorUrlPatternsRuleDto{},
		}
		urlPattern.Rules = append(urlPattern.Rules, dlpConnectorUrlPatternsRuleDto{
			Order:                       customConnectorPattern.Order.ValueInt64(),
			ConnectorRuleClassification: convertConnectorRuleClassificationValues(customConnectorPattern.DataGroup.ValueString()),
			Pattern:                     customConnectorPattern.HostUrlPattern.ValueString(),
		})
		customConnectorUrlPatternsDefinition = append(customConnectorUrlPatternsDefinition, urlPattern)
	}
	return customConnectorUrlPatternsDefinition
}

func convertToDlpActionRule(connector dataLossPreventionPolicyResourceConnectorModel) []dlpActionRuleDto {
	var actionRules []dlpActionRuleDto
	for _, actionRule := range connector.ActionRules {
		actionRules = append(actionRules, dlpActionRuleDto{
			ActionId: actionRule.ActionId.ValueString(),
			Behavior: actionRule.Behavior.ValueString(),
		})
	}
	return actionRules
}

func convertToDlpEndpointRule(connector dataLossPreventionPolicyResourceConnectorModel) []dlpEndpointRuleDto {
	var endpointRules []dlpEndpointRuleDto
	for _, endpointRule := range connector.EndpointRules {
		endpointRules = append(endpointRules, dlpEndpointRuleDto{
			Order:    endpointRule.Order.ValueInt64(),
			Behavior: endpointRule.Behavior.ValueString(),
			Endpoint: endpointRule.Endpoint.ValueString(),
		})
	}
	return endpointRules
}

func convertToAtrValueActionRule(connector dlpConnectorModelDto) []attr.Value {
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

func convertToAtrValueEndpointRule(connector dlpConnectorModelDto) []attr.Value {
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
