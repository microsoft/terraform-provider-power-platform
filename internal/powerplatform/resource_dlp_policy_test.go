package powerplatform

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

// func TestUnitDataLossPreventionPolicyResource_Validate_Update(t *testing.T) {
// 	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)

// 	policyId := "00000000-0000-0000-0000-000000000000"
// 	policy := models.DlpPolicyModel{
// 		Name:             policyId,
// 		ETag:             "etag",
// 		CreatedBy:        "createdBy",
// 		CreatedTime:      "createdTime",
// 		LastModifiedBy:   "lastModifiedBy",
// 		LastModifiedTime: "lastModifiedTime",
// 	}

// 	steps := []resource.TestStep{
// 		{
// 			Config: UniTestsProviderConfig + `
// 			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
// 				display_name                      = "Block All Policy"
// 				default_connectors_classification = "Blocked"
// 				environment_type                  = "AllEnvironments"
// 				environments = []

// 				business_connectors = []
// 				non_business_connectors = []
// 				blocked_connectors = []
// 				custom_connectors_patterns = []
// 			  }`,

// 			Check: resource.ComposeTestCheckFunc(
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "id", policyId),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "display_name", "Block All Policy"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "default_connectors_classification", "Blocked"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environment_type", "AllEnvironments"),
// 			),
// 		},
// 		{
// 			Config: UniTestsProviderConfig + `
// 			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
// 				display_name                      = "Block All Policy_1"
// 				default_connectors_classification = "Blocked"
// 				environment_type                  = "AllEnvironments"
// 				environments = []

// 				business_connectors = []
// 				non_business_connectors = []
// 				blocked_connectors = []
// 				custom_connectors_patterns = []
// 			  }`,

// 			Check: resource.ComposeTestCheckFunc(
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "display_name", "Block All Policy_1"),
// 			),
// 		},
// 		{
// 			Config: UniTestsProviderConfig + `
// 			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
// 				display_name                      = "Block All Policy"
// 				default_connectors_classification = "General"
// 				environment_type                  = "OnlyEnvironments"
// 				environments = [
// 					{
// 						name = "00000000-0000-0000-0000-000000000000"
// 					}
// 				]

// 				business_connectors = []
// 				non_business_connectors = []
// 				blocked_connectors = []
// 				custom_connectors_patterns = []
// 			  }`,

// 			Check: resource.ComposeTestCheckFunc(
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "default_connectors_classification", "General"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environment_type", "OnlyEnvironments"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.#", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.0.name", "00000000-0000-0000-0000-000000000000"),
// 			),
// 		},
// 		{
// 			Config: UniTestsProviderConfig + `
// 			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
// 				display_name                      = "Block All Policy"
// 				default_connectors_classification = "General"
// 				environment_type                  = "OnlyEnvironments"
// 				environments = [
// 					{
// 						name = "00000000-0000-0000-0000-000000000000"
// 					}
// 				]

// 				business_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 				non_business_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 				blocked_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 				custom_connectors_patterns = []
// 			  }`,

// 			Check: resource.ComposeTestCheckFunc(
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.#", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.default_action_rule_behavior", "Allow"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.#", "0"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.#", "0"),

// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.#", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.default_action_rule_behavior", "Allow"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.action_rules.#", "0"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.endpoint_rules.#", "0"),

// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.#", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_azureblob"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.default_action_rule_behavior", "Allow"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.action_rules.#", "0"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.endpoint_rules.#", "0"),

// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.#", "0"),
// 			),
// 		},
// 		{
// 			Config: UniTestsProviderConfig + `
// 			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
// 				display_name                      = "Block All Policy"
// 				default_connectors_classification = "General"
// 				environment_type                  = "OnlyEnvironments"
// 				environments = [
// 					{
// 						name = "00000000-0000-0000-0000-000000000000"
// 					}
// 				]

// 				business_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [
// 							  {
// 								action_id = "DeleteItem_V2",
// 								behavior  = "Block",
// 							  },
// 							  {
// 								action_id = "ExecutePassThroughNativeQuery_V2",
// 								behavior  = "Block",
// 							  }
// 							],
// 							endpoint_rules = [
// 							  {
// 								order    = 1,
// 								behavior = "Allow",
// 								endpoint = "contoso.com"
// 							  },
// 							  {
// 								order    = 2,
// 								behavior = "Deny",
// 								endpoint = "*"
// 							  }
// 							]
// 						}
// 					])
// 				non_business_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 				blocked_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 				custom_connectors_patterns = []
// 			  }`,

// 			Check: resource.ComposeTestCheckFunc(
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.#", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.default_action_rule_behavior", "Allow"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.#", "2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.action_id", "DeleteItem_V2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.behavior", "Block"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.action_id", "ExecutePassThroughNativeQuery_V2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.behavior", "Block"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.#", "2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.order", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.behavior", "Allow"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.endpoint", "contoso.com"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.order", "2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.behavior", "Deny"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.endpoint", "*"),
// 			),
// 		},
// 		{
// 			Config: UniTestsProviderConfig + `
// 			resource "powerplatform_data_loss_prevention_policy" "my_policy" {
// 				display_name                      = "Block All Policy"
// 				default_connectors_classification = "General"
// 				environment_type                  = "OnlyEnvironments"
// 				environments = [
// 					{
// 						name = "00000000-0000-0000-0000-000000000000"
// 					}
// 				]

// 				business_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [
// 							  {
// 								action_id = "DeleteItem_V2",
// 								behavior  = "Block",
// 							  },
// 							  {
// 								action_id = "ExecutePassThroughNativeQuery_V2",
// 								behavior  = "Block",
// 							  }
// 							],
// 							endpoint_rules = [
// 							  {
// 								order    = 1,
// 								behavior = "Allow",
// 								endpoint = "contoso.com"
// 							  },
// 							  {
// 								order    = 2,
// 								behavior = "Deny",
// 								endpoint = "*"
// 							  }
// 							]
// 						}
// 					])
// 				non_business_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 				blocked_connectors = toset([
// 						{
// 							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
// 							default_action_rule_behavior = "Allow",
// 							action_rules = [],
// 							endpoint_rules = [],
// 						}
// 					])
// 					custom_connectors_patterns = toset([
// 						{
// 						  order            = 1
// 						  host_url_pattern = "https://*.contoso.com"
// 						  data_group       = "Blocked"
// 						},
// 						{
// 						  order            = 2
// 						  host_url_pattern = "*"
// 						  data_group       = "Ignore"
// 						}
// 					  ])
// 			  }`,

// 			Check: resource.ComposeTestCheckFunc(
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.#", "2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.order", "1"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.host_url_pattern", "https://*.contoso.com"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.data_group", "Blocked"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.order", "2"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.host_url_pattern", "*"),
// 				resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.data_group", "Ignore"),
// 			),
// 		},
// 	}

// 	clientMock.EXPECT().UpdatePolicy(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, name string, policyToUpdate models.DlpPolicyModel) (*models.DlpPolicyModel, error) {
// 		policy.DisplayName = policyToUpdate.DisplayName
// 		policy.DefaultConnectorsClassification = policyToUpdate.DefaultConnectorsClassification
// 		policy.EnvironmentType = policyToUpdate.EnvironmentType
// 		policy.Environments = policyToUpdate.Environments
// 		policy.ConnectorGroups = policyToUpdate.ConnectorGroups
// 		policy.CustomConnectorUrlPatternsDefinition = policyToUpdate.CustomConnectorUrlPatternsDefinition
// 		return &policy, nil
// 	}).Times(len(steps) - 1)

// 	clientMock.EXPECT().GetPolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyId string) (*models.DlpPolicyModel, error) {
// 		return &policy, nil
// 	}).AnyTimes()

// 	clientMock.EXPECT().CreatePolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyToCreate models.DlpPolicyModel) (*models.DlpPolicyModel, error) {
// 		policy.DisplayName = policyToCreate.DisplayName
// 		policy.DefaultConnectorsClassification = policyToCreate.DefaultConnectorsClassification
// 		policy.EnvironmentType = policyToCreate.EnvironmentType
// 		policy.Environments = policyToCreate.Environments
// 		policy.ConnectorGroups = policyToCreate.ConnectorGroups
// 		policy.CustomConnectorUrlPatternsDefinition = policyToCreate.CustomConnectorUrlPatternsDefinition

// 		return &policy, nil
// 	}).Times(1)

// 	clientMock.EXPECT().DeletePolicy(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, policyId string) error {
// 		return nil
// 	}).Times(1)

// 	resource.Test(t, resource.TestCase{
// 		IsUnitTest: true,
// 		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
// 			"powerplatform": powerPlatformProviderServerApiMock(clientMock, nil, nil),
// 		},
// 		Steps: steps,
// 	})
// }

func TestUnitDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	policyId := "00000000-0000-0000-0000-000000000000"
	policyResponse := fmt.Sprintf(`{
		"policyDefinition": {
			"name": "%s",
			"displayName": "Block All Policy",
			"defaultConnectorsClassification": "Blocked",
			"connectorGroups": [
				{
					"classification": "Confidential",
					"connectors": [
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
							"name": "shared_sharepointonline",
							"type": "Microsoft.PowerApps/apis"
						}
					]
				},
				{
					"classification": "General",
					"connectors": [
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_sql",
							"name": "shared_sql",
							"type": "Microsoft.PowerApps/apis"
						}
					]
				},
				{
					"classification": "Blocked",
					"connectors": [
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							"name": "shared_azureblob",
							"type": "Microsoft.PowerApps/apis"
						}
					]
				}
			],
			"environmentType": "OnlyEnvironments",
			"environments": [
				{
					"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000000",
					"name": "00000000-0000-0000-0000-000000000000",
					"type": "Microsoft.BusinessAppPlatform/scopes/environments"
				}
			],
			"createdBy": {
				"displayName": "createdBy"
			},
			"createdTime": "createdTime",
			"lastModifiedBy": {
				"displayName": "lastModifiedBy"
			},
			"lastModifiedTime": "lastModifiedTime",
			"etag": "etag",
			"isLegacySchemaVersion": false
		},
		"connectorConfigurationsDefinition": {
			"connectorActionConfigurations": [
				{
					"connectorId": "/providers/Microsoft.PowerApps/apis/shared_sql",
					"actionRules": [
						{
							"actionId": "DeleteItem_V2",
							"behavior": "Block"
						},
						{
							"actionId": "ExecutePassThroughNativeQuery_V2",
							"behavior": "Block"
						}
					],
					"defaultConnectorActionRuleBehavior": "Allow"
				}
			],
			"endpointConfigurations": [
				{
					"connectorId": "/providers/Microsoft.PowerApps/apis/shared_sql",
					"endpointRules": [
						{
							"order": 1,
							"behavior": "Allow",
							"endPoint": "contoso.com"
						},
						{
							"order": 2,
							"behavior": "Deny",
							"endPoint": "*"
						}
					]
				}
			]
		},
		"customConnectorUrlPatternsDefinition": {
			"rules": [
				{
					"order": 1,
					"customConnectorRuleClassification": "Blocked",
					"pattern": "https://*.contoso.com"
				},
				{
					"order": 2,
					"customConnectorRuleClassification": "Ignore",
					"pattern": "*"
				}
			]
		}
	}`, policyId)

	httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, policyResponse), nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf(`https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/%s`, policyId),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, policyResponse), nil
		})

	httpmock.RegisterResponder("DELETE", fmt.Sprintf(`https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/policies/%s`, policyId),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_data_loss_prevention_policy" "my_policy" {
					display_name                      = "Block All Policy"
					default_connectors_classification = "Blocked"
					environment_type                  = "OnlyEnvironments"
					environments = [
						{
							name = "00000000-0000-0000-0000-000000000000"
						}
					]

					business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sql"
							default_action_rule_behavior = "Allow",
							action_rules = [
							  {
								action_id = "DeleteItem_V2",
								behavior  = "Block",
							  },
							  {
								action_id = "ExecutePassThroughNativeQuery_V2",
								behavior  = "Block",
							  }
							],
							endpoint_rules = [
							  {
								order    = 1,
								behavior = "Allow",
								endpoint = "contoso.com"
							  },
							  {
								order    = 2,
								behavior = "Deny",
								endpoint = "*"
							  }
							]
						  }
					])
					non_business_connectors = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
							default_action_rule_behavior = "",
							action_rules                 = [],
							endpoint_rules               = []
						},
					])
					blocked_connectors      = toset([
						{
							id                           = "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							default_action_rule_behavior = "",
							action_rules                 = []
							endpoint_rules               = []
						  },
					])
					custom_connectors_patterns = toset([
					  {
						order            = 1
						host_url_pattern = "https://*.contoso.com"
						data_group       = "Blocked"
					  },
					  {
						order            = 2
						host_url_pattern = "*"
						data_group       = "Ignore"
					  }
					])
				  }`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "id", policyId),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "display_name", "Block All Policy"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "default_connectors_classification", "Blocked"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environment_type", "OnlyEnvironments"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "environments.0.name", "00000000-0000-0000-0000-000000000000"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.default_action_rule_behavior", "Allow"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.action_id", "DeleteItem_V2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.0.behavior", "Block"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.action_id", "ExecutePassThroughNativeQuery_V2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.action_rules.1.behavior", "Block"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.order", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.behavior", "Allow"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.0.endpoint", "contoso.com"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.order", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.behavior", "Deny"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "business_connectors.0.endpoint_rules.1.endpoint", "*"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.default_action_rule_behavior", ""),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.action_rules.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.endpoint_rules.#", "0"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_azureblob"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.default_action_rule_behavior", ""),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.action_rules.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "blocked_connectors.0.endpoint_rules.#", "0"),

					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.order", "1"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.host_url_pattern", "https://*.contoso.com"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.data_group", "Blocked"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.order", "2"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.host_url_pattern", "*"),
					resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.1.data_group", "Ignore"),
				),
			},
		},
	})
}
