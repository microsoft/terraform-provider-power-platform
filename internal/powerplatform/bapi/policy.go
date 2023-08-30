package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetPolicies(ctx context.Context) ([]models.DlpPolicyModel, error) {

	//https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/policies
	//https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/<tenantId>/policies/<policyId>/policyconnectorconfigurations
	//https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/<tenantId>/policies/<policyId>/urlPatterns

	return nil, nil
}

func (client *ApiClient) GetPolicy(ctx context.Context, name string) (*models.DlpPolicyModel, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", name),
	}
	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
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
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/policies/%s", name),
	}
	request, err := http.NewRequestWithContext(ctx, "DELETE", apiUrl.String(), nil)
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

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.Name),
	}
	request, err := http.NewRequestWithContext(ctx, "PATCH", apiUrl.String(), bytes.NewReader(body))
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

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/PowerPlatform.Governance/v2/policies/",
	}
	request, err := http.NewRequestWithContext(ctx, "POST", apiUrl.String(), bytes.NewReader(body))
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
			nameSplit := strings.Split(connector.Id, "/")
			con := models.DlpConnectorDto{
				Id:   connector.Id,
				Name: nameSplit[len(nameSplit)-1],
				Type: connector.Type,
			}
			conG.Connectors = append(conG.Connectors, con)
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
