package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewDlpPolicyClient(bapi *api.BapiClientApi) DlpPolicyClient {
	return DlpPolicyClient{
		BapiClient: bapi,
	}
}

type DlpPolicyClient struct {
	BapiClient *api.BapiClientApi
}

func (client *DlpPolicyClient) GetPolicies(ctx context.Context) ([]DlpPolicyModelDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BapiClient.GetConfig().Urls.BapiUrl,
		Path:   "providers/PowerPlatform.Governance/v2/policies",
	}
	policiesArray := DlpPolicyDefinitionDtoArray{}
	_, err := client.BapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policiesArray)
	if err != nil {
		return nil, err
	}

	policies := make([]DlpPolicyModelDto, 0)
	for _, policy := range policiesArray.Value {

		apiUrl := &url.URL{
			Scheme: "https",
			Host:   client.BapiClient.GetConfig().Urls.BapiUrl,
			Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.PolicyDefinition.Name),
		}
		policy := DlpPolicyDto{}
		_, err := client.BapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)
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

func (client *DlpPolicyClient) GetPolicy(ctx context.Context, name string) (*DlpPolicyModelDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", name),
	}
	policy := DlpPolicyDto{}
	_, err := client.BapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)
	if err != nil {
		return nil, err
	}
	return covertDlpPolicyToPolicyModel(policy)
}

func (client *DlpPolicyClient) DeletePolicy(ctx context.Context, name string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v1/policies/%s", name),
	}
	_, err := client.BapiClient.Execute(ctx, "DELETE", apiUrl.String(), nil, []int{http.StatusNoContent}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *DlpPolicyClient) UpdatePolicy(ctx context.Context, name string, policy DlpPolicyModelDto) (*DlpPolicyModelDto, error) {
	policyToCreate := convertPolicyModelToDlpPolicy(policy)

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.Name),
	}
	createdPolicy := DlpPolicyDto{}

	_, err := client.BapiClient.Execute(ctx, "PATCH", apiUrl.String(), policyToCreate, []int{http.StatusAccepted}, &createdPolicy)
	if err != nil {
		return nil, err
	}

	return covertDlpPolicyToPolicyModel(createdPolicy)
}

func convertPolicyModelToDlpPolicy(policy DlpPolicyModelDto) DlpPolicyDto {
	policyToCreate := DlpPolicyDto{
		PolicyDefinition: DlpPolicyDefinitionDto{
			Name:                            policy.Name,
			DisplayName:                     policy.DisplayName,
			DefaultConnectorsClassification: policy.DefaultConnectorsClassification,
			EnvironmentType:                 policy.EnvironmentType,
			Environments:                    policy.Environments,
			ConnectorGroups:                 []DlpConnectorGroupsDto{},
		},
		ConnectorConfigurationsDefinition:    &policy.ConnectorConfigurationsDefinition,
		CustomConnectorUrlPatternsDefinition: DlpConnectorUrlPatternsDefinitionDto{},
	}

	for _, policy := range policy.CustomConnectorUrlPatternsDefinition {
		policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, DlpConnectorUrlPatternsRuleDto{
			Order:                       policy.Rules[0].Order,
			ConnectorRuleClassification: policy.Rules[0].ConnectorRuleClassification,
			Pattern:                     policy.Rules[0].Pattern,
		})
	}

	for _, connGroups := range policy.ConnectorGroups {
		conG := DlpConnectorGroupsDto{
			Classification: connGroups.Classification,
			Connectors:     []DlpConnectorDto{},
		}

		for _, connector := range connGroups.Connectors {
			nameSplit := strings.Split(connector.Id, "/")
			con := DlpConnectorDto{
				Id:   connector.Id,
				Name: nameSplit[len(nameSplit)-1],
				Type: connector.Type,
			}
			conG.Connectors = append(conG.Connectors, con)
		}
		policyToCreate.PolicyDefinition.ConnectorGroups = append(policyToCreate.PolicyDefinition.ConnectorGroups, conG)
	}

	connectorActionConfigurationsDto := []DlpConnectorActionConfigurationsDto{}
	endpointConfigurationsDto := []DlpEndpointConfigurationsDto{}

	for _, connGroups := range policy.ConnectorGroups {
		for _, connector := range connGroups.Connectors {
			if connector.ActionRules != nil && len(connector.ActionRules) > 0 {
				connectorActionConfigurationsDto = append(connectorActionConfigurationsDto, DlpConnectorActionConfigurationsDto{
					ConnectorId:                        connector.Id,
					DefaultConnectorActionRuleBehavior: connector.DefaultActionRuleBehavior,
					ActionRules:                        connector.ActionRules,
				})
			}
			if connector.EndpointRules != nil && len(connector.EndpointRules) > 0 {
				endpointConfigurationsDto = append(endpointConfigurationsDto, DlpEndpointConfigurationsDto{
					ConnectorId:   connector.Id,
					EndpointRules: connector.EndpointRules,
				})
			}
		}
	}

	if len(connectorActionConfigurationsDto) > 0 || len(endpointConfigurationsDto) > 0 {
		policyToCreate.ConnectorConfigurationsDefinition = &DlpConnectorConfigurationsDefinitionDto{}

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

func covertDlpPolicyToPolicyModel(policy DlpPolicyDto) (*DlpPolicyModelDto, error) {

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
		policyModel.CustomConnectorUrlPatternsDefinition = append(policyModel.CustomConnectorUrlPatternsDefinition, DlpConnectorUrlPatternsDefinitionDto{
			Rules: append([]DlpConnectorUrlPatternsRuleDto{}, rule),
		})
	}

	return &policyModel, nil
}

func (client *DlpPolicyClient) CreatePolicy(ctx context.Context, policy DlpPolicyModelDto) (*DlpPolicyModelDto, error) {

	policyToCreate := convertPolicyModelToDlpPolicy(policy)

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BapiClient.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v2/policies",
	}

	createdPolicy := DlpPolicyDto{}
	_, err := client.BapiClient.Execute(ctx, "POST", apiUrl.String(), policyToCreate, []int{http.StatusCreated}, &createdPolicy)
	if err != nil {
		return nil, err
	}
	return covertDlpPolicyToPolicyModel(createdPolicy)
}
