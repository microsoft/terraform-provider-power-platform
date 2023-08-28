package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetPolicies(ctx context.Context) ([]models.DlpPolicyModel, error) {

	//https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/policies
	//https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/1dbbeae5-8fa6-462e-a5a1-9932a520a1dc/policies/9faa9dca-9d96-41b3-888e-98b3d8911f88/policyconnectorconfigurations
	//https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/1dbbeae5-8fa6-462e-a5a1-9932a520a1dc/policies/9faa9dca-9d96-41b3-888e-98b3d8911f88/urlPatterns

	return nil, nil
}

func (client *ApiClient) GetPolicy(ctx context.Context, name string) (*models.DlpPolicyModel, error) {

	request, err := http.NewRequestWithContext(ctx, "GET", "https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/"+name, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	policy := models.DlpPolicyDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if err != nil {
		return nil, err
	}

	return covertDlpPolicyToPolicyModel(policy)
}

func (client *ApiClient) DeletePolicy(ctx context.Context, name string) error {
	request, err := http.NewRequestWithContext(ctx, "DELETE", "https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/policies/"+name, nil)
	if err != nil {
		return err
	}

	_, err = client.doRequest(request)
	if err != nil {
		return err
	}
	return nil
}

func (client *ApiClient) UpdatePolicy(ctx context.Context, name string, policy models.DlpPolicyModel) (*models.DlpPolicyModel, error) {
	policyToCreate := convertPolicyModelToDlpPolicy(policy)

	body, err := json.Marshal(policyToCreate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "PATCH", "https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/"+policy.Name, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	createdPolicy := models.DlpPolicyDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&createdPolicy)
	if err != nil {
		return nil, err
	}

	return covertDlpPolicyToPolicyModel(createdPolicy)
}

func (client *ApiClient) CreatePolicy(ctx context.Context, policy models.DlpPolicyModel) (*models.DlpPolicyModel, error) {

	policyToCreate := convertPolicyModelToDlpPolicy(policy)

	body, err := json.Marshal(policyToCreate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", "https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	createdPolicy := models.DlpPolicyDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&createdPolicy)
	if err != nil {
		return nil, err
	}

	return covertDlpPolicyToPolicyModel(createdPolicy)
}

func convertPolicyModelToDlpPolicy(policy models.DlpPolicyModel) models.DlpPolicyDto {
	policyToCreate := models.DlpPolicyDto{
		PolicyDefinition: models.DlpPolicyDefinitionDto{
			Name:                            policy.Name,
			DisplayName:                     policy.DisplayName,
			DefaultConnectorsClassification: policy.DefaultConnectorsClassification,
			EnvironmentType:                 policy.EnvironmentType,
			Environments:                    policy.Environments,
			ConnectorGroups:                 []models.DlpConnectorGroupsDto{},
		},
		ConnectorConfigurationsDefinition:    policy.ConnectorConfigurationsDefinition,
		CustomConnectorUrlPatternsDefinition: models.DlpConnectorUrlPatternsDefinitionDto{},
	}

	for _, policy := range policy.CustomConnectorUrlPatternsDefinition {
		policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, models.DlpConnectorUrlPatternsRuleDto{
			Order:                       policy.Rules[0].Order,
			ConnectorRuleClassification: policy.Rules[0].ConnectorRuleClassification,
			Pattern:                     policy.Rules[0].Pattern,
		})
	}

	for _, connGroups := range policy.ConnectorGroups {
		conG := models.DlpConnectorGroupsDto{
			Classification: connGroups.Classification,
			Connectors:     []models.DlpConnectorDto{},
		}

		for _, connector := range connGroups.Connectors {
			conG.Connectors = append(conG.Connectors, models.DlpConnectorDto{
				Id:   connector.Id,
				Name: connector.Name,
				Type: connector.Type,
			})
		}
		policyToCreate.PolicyDefinition.ConnectorGroups = append(policyToCreate.PolicyDefinition.ConnectorGroups, conG)
	}

	connectorActionConfigurationsDto := []models.DlpConnectorActionConfigurationsDto{}
	endpointConfigurationsDto := []models.DlpEndpointConfigurationsDto{}

	for _, connGroups := range policy.ConnectorGroups {
		for _, connector := range connGroups.Connectors {
			if connector.ActionRules != nil && len(connector.ActionRules) > 0 {
				connectorActionConfigurationsDto = append(connectorActionConfigurationsDto, models.DlpConnectorActionConfigurationsDto{
					ConnectorId:                        connector.Id,
					DefaultConnectorActionRuleBehavior: connector.DefaultActionRuleBehavior,
					ActionRules:                        connector.ActionRules,
				})
			}
			if connector.EndpointRules != nil && len(connector.EndpointRules) > 0 {
				endpointConfigurationsDto = append(endpointConfigurationsDto, models.DlpEndpointConfigurationsDto{
					ConnectorId:   connector.Id,
					EndpointRules: connector.EndpointRules,
				})
			}
		}
	}

	if len(connectorActionConfigurationsDto) > 0 || len(endpointConfigurationsDto) > 0 {
		policyToCreate.ConnectorConfigurationsDefinition = models.DlpConnectorConfigurationsDefinitionDto{}

		if len(connectorActionConfigurationsDto) > 0 {
			policyToCreate.ConnectorConfigurationsDefinition.ConnectorActionConfigurations = connectorActionConfigurationsDto
		}
		if len(endpointConfigurationsDto) > 0 {
			policyToCreate.ConnectorConfigurationsDefinition.EndpointConfigurations = endpointConfigurationsDto
		}
	}
	return policyToCreate
}

func covertDlpPolicyToPolicyModel(policy models.DlpPolicyDto) (*models.DlpPolicyModel, error) {

	policyModel := models.DlpPolicyModel{
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
		ConnectorConfigurationsDefinition:    models.DlpConnectorConfigurationsDefinitionDto{},
		CustomConnectorUrlPatternsDefinition: []models.DlpConnectorUrlPatternsDefinitionDto{},
		ConnectorGroups:                      []models.DlpConnectorGroupsModel{},
	}

	for _, connGroup := range policy.PolicyDefinition.ConnectorGroups {
		connGroupModel := models.DlpConnectorGroupsModel{
			Classification: connGroup.Classification,
			Connectors:     []models.DlpConnectorModel{},
		}
		for _, connector := range connGroup.Connectors {
			nameSplit := strings.Split(connector.Id, "/")
			m := models.DlpConnectorModel{
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
		policyModel.CustomConnectorUrlPatternsDefinition = append(policyModel.CustomConnectorUrlPatternsDefinition, models.DlpConnectorUrlPatternsDefinitionDto{
			Rules: append([]models.DlpConnectorUrlPatternsRuleDto{}, rule),
		})
	}

	return &policyModel, nil
}
