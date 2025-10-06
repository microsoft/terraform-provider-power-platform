// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newDlpPolicyClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetPolicies(ctx context.Context) ([]dlpPolicyModelDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "providers/PowerPlatform.Governance/v2/policies",
	}
	policiesArray := dlpPolicyDefinitionArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policiesArray)
	if err != nil {
		return nil, err
	}

	policies := make([]dlpPolicyModelDto, 0)
	for _, policy := range policiesArray.Value {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   client.Api.GetConfig().Urls.BapiUrl,
			Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.PolicyDefinition.Name),
		}
		policy := dlpPolicyDto{}
		_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)
		if err != nil {
			return nil, err
		}
		v, err := covertDlpPolicyToPolicyModelDto(policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, *v)
	}
	return policies, nil
}

func (client *client) GetPolicy(ctx context.Context, name string) (*dlpPolicyModelDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", name),
	}
	policy := dlpPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)
	if err != nil {
		var httpError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("DLP Policy '%s' not found", name))
		}
		return nil, err
	}
	return covertDlpPolicyToPolicyModel(policy)
}

func (client *client) DeletePolicy(ctx context.Context, name string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", name),
	}
	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) UpdatePolicy(ctx context.Context, policy dlpPolicyModelDto) (*dlpPolicyModelDto, error) {
	policyToCreate := convertPolicyModelToDlpPolicy(policy)

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.Name),
	}
	createdPolicy := dlpPolicyDto{}

	_, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, policyToCreate, []int{http.StatusOK}, &createdPolicy)
	if err != nil {
		return nil, err
	}

	return covertDlpPolicyToPolicyModel(createdPolicy)
}

func convertPolicyModelToDlpPolicy(policy dlpPolicyModelDto) dlpPolicyDto {
	policyToCreate := dlpPolicyDto{
		PolicyDefinition: dlpPolicyDefinitionDto{
			Name:                            policy.Name,
			DisplayName:                     policy.DisplayName,
			DefaultConnectorsClassification: policy.DefaultConnectorsClassification,
			EnvironmentType:                 policy.EnvironmentType,
			Environments:                    policy.Environments,
			ConnectorGroups:                 []dlpConnectorGroupsDto{},
		},
		ConnectorConfigurationsDefinition:    &policy.ConnectorConfigurationsDefinition,
		CustomConnectorUrlPatternsDefinition: dlpConnectorUrlPatternsDefinitionDto{},
	}

	for _, policy := range policy.CustomConnectorUrlPatternsDefinition {
		policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, dlpConnectorUrlPatternsRuleDto{
			Order:                       policy.Rules[0].Order,
			ConnectorRuleClassification: policy.Rules[0].ConnectorRuleClassification,
			Pattern:                     policy.Rules[0].Pattern,
		})
	}

	for _, connGroups := range policy.ConnectorGroups {
		conG := dlpConnectorGroupsDto{
			Classification: connGroups.Classification,
			Connectors:     []dlpConnectorDto{},
		}

		for _, connector := range connGroups.Connectors {
			nameSplit := strings.Split(connector.Id, "/")
			con := dlpConnectorDto{
				Id:   connector.Id,
				Name: nameSplit[len(nameSplit)-1],
				Type: connector.Type,
			}
			conG.Connectors = append(conG.Connectors, con)
		}
		policyToCreate.PolicyDefinition.ConnectorGroups = append(policyToCreate.PolicyDefinition.ConnectorGroups, conG)
	}

	connectorActionConfigurationsDto := []dlpConnectorActionConfigurationsDto{}
	endpointConfigurationsDto := []dlpEndpointConfigurationsDto{}

	for _, connGroups := range policy.ConnectorGroups {
		for _, connector := range connGroups.Connectors {
			if len(connector.ActionRules) > 0 {
				connectorActionConfigurationsDto = append(connectorActionConfigurationsDto, dlpConnectorActionConfigurationsDto{
					ConnectorId:                        connector.Id,
					DefaultConnectorActionRuleBehavior: connector.DefaultActionRuleBehavior,
					ActionRules:                        connector.ActionRules,
				})
			}
			if len(connector.EndpointRules) > 0 {
				endpointConfigurationsDto = append(endpointConfigurationsDto, dlpEndpointConfigurationsDto{
					ConnectorId:   connector.Id,
					EndpointRules: connector.EndpointRules,
				})
			}
		}
	}

	if len(connectorActionConfigurationsDto) > 0 || len(endpointConfigurationsDto) > 0 {
		policyToCreate.ConnectorConfigurationsDefinition = &dlpConnectorConfigurationsDefinitionDto{}

		if len(connectorActionConfigurationsDto) > 0 {
			policyToCreate.ConnectorConfigurationsDefinition.ConnectorActionConfigurations = connectorActionConfigurationsDto
		}
		if len(endpointConfigurationsDto) > 0 {
			policyToCreate.ConnectorConfigurationsDefinition.EndpointConfigurations = endpointConfigurationsDto
		}
	} else {
		policyToCreate.ConnectorConfigurationsDefinition = nil
	}
	return policyToCreate
}

func covertDlpPolicyToPolicyModel(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
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

	sort.Slice(policyModel.CustomConnectorUrlPatternsDefinition, func(i, j int) bool {
		// Compare the Order of the first rule in each definition (assuming at least one rule per definition)
		return policyModel.CustomConnectorUrlPatternsDefinition[i].Rules[0].Order < policyModel.CustomConnectorUrlPatternsDefinition[j].Rules[0].Order
	})

	return &policyModel, nil
}

func (client *client) CreatePolicy(ctx context.Context, policy dlpPolicyModelDto) (*dlpPolicyModelDto, error) {
	policyToCreate := convertPolicyModelToDlpPolicy(policy)

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v2/policies",
	}

	createdPolicy := dlpPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, &createdPolicy)
	if err != nil {
		return nil, err
	}
	return covertDlpPolicyToPolicyModel(createdPolicy)
}
