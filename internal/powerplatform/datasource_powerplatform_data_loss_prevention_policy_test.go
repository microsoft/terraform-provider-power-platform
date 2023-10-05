package powerplatform

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestUnitDlpPolicyDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateOAuthHttpMocks()

	const policyId1 = "16c21e0d-429e-4e37-b496-f3c1bcd78bfe"
	const policyId2 = "79ce0ded-5539-4bc6-823e-3176d73371fc"

	policiesResponose := fmt.Sprintf(`{
		"value": [
			{
				"policyDefinition": {
					"name": "%s",
					"displayName": "a1",
					"defaultConnectorsClassification": "General",
					"environmentType": "AllEnvironments",
					"environments": [],
					"createdBy": {
						"displayName": "admin"
					},
					"createdTime": "2023-10-02T07:38:50.3269899Z",
					"lastModifiedBy": {
						"displayName": "admin"
					},
					"lastModifiedTime": "2023-10-02T07:38:50.3269899Z",
					"etag": "dcf783da-6eb1-4c5a-a6ee-118a64bafbdb",
					"isLegacySchemaVersion": false
				}
			},
			{
				"policyDefinition": {
					"name": "%s",
					"displayName": "a2",
					"defaultConnectorsClassification": "General",
					"environmentType": "ExceptEnvironments",
					"environments": [
						{
							"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/be0eb809-e58a-ec1b-8fce-ea40b0e53442",
							"name": "be0eb809-e58a-ec1b-8fce-ea40b0e53442",
							"type": "Microsoft.BusinessAppPlatform/scopes/environments"
						}
					],
					"createdBy": {
						"displayName": "admin"
					},
					"createdTime": "2023-10-02T07:38:56.6864176Z",
					"lastModifiedBy": {
						"displayName": "admin"
					},
					"lastModifiedTime": "2023-10-02T07:56:43.9700369Z",
					"etag": "a872cb45-ee20-4f63-a8e6-fcb537bd8aaf",
					"isLegacySchemaVersion": false
				}
			}
		]
	}`, policyId1, policyId2)
	const policy1Response = `{
		"policyDefinition": {
			"name": "16c21e0d-429e-4e37-b496-f3c1bcd78bfe",
			"displayName": "a1",
			"defaultConnectorsClassification": "General",
			"connectorGroups": [
				{
					"classification": "Confidential",
					"connectors": []
				},
				{
					"classification": "General",
					"connectors": []
				},
				{
					"classification": "Blocked",
					"connectors": []
				}
			],
			"environmentType": "AllEnvironments",
			"environments": [],
			"createdBy": {
				"displayName": "admin"
			},
			"createdTime": "2023-10-02T07:38:50.3269899Z",
			"lastModifiedBy": {
				"displayName": "admin"
			},
			"lastModifiedTime": "2023-10-02T07:38:50.3269899Z",
			"etag": "dcf783da-6eb1-4c5a-a6ee-118a64bafbdb",
			"isLegacySchemaVersion": false
		},
		"customConnectorUrlPatternsDefinition": {
			"rules": [
				{
					"order": 1,
					"customConnectorRuleClassification": "Ignore",
					"pattern": "*"
				}
			]
		}
	}`
	const policy2Response = `{
		"policyDefinition": {
			"name": "79ce0ded-5539-4bc6-823e-3176d73371fc",
			"displayName": "a2",
			"defaultConnectorsClassification": "General",
			"connectorGroups": [
				{
					"classification": "Confidential",
					"connectors": [
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_office365users",
							"name": "Office 365 Users",
							"type": "Microsoft.PowerApps/apis"
						},
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_azureblob",
							"name": "Azure Blob Storage",
							"type": "Microsoft.PowerApps/apis"
						}
					]
				},
				{
					"classification": "General",
					"connectors": [
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_powerappsforappmakers",
							"name": "Power Apps for Makers",
							"type": "Microsoft.PowerApps/apis"
						},
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_microsoftspatialservices",
							"name": "Spatial Services",
							"type": "Microsoft.PowerApps/apis"
						},
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_sql",
							"name": "SQL Server",
							"type": "Microsoft.PowerApps/apis"
						},
						{
							"id": "/providers/Microsoft.PowerApps/apis/shared_bttn",
							"name": "bttn",
							"type": "Microsoft.PowerApps/apis"
						}
					]
				},
				{
					"classification": "Blocked",
					"connectors": []
				}
			],
			"environmentType": "ExceptEnvironments",
			"environments": [
				{
					"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/be0eb809-e58a-ec1b-8fce-ea40b0e53442",
					"name": "be0eb809-e58a-ec1b-8fce-ea40b0e53442",
					"type": "Microsoft.BusinessAppPlatform/scopes/environments"
				}
			],
			"createdBy": {
				"displayName": "admin"
			},
			"createdTime": "2023-10-02T07:38:56.6864176Z",
			"lastModifiedBy": {
				"displayName": "admin"
			},
			"lastModifiedTime": "2023-10-02T07:56:43.9700369Z",
			"etag": "a872cb45-ee20-4f63-a8e6-fcb537bd8aaf",
			"isLegacySchemaVersion": false
		},
		"connectorConfigurationsDefinition": {
			"connectorActionConfigurations": [
				{
					"connectorId": "/providers/Microsoft.PowerApps/apis/shared_azureblob",
					"actionRules": [
						{
							"actionId": "CreateFile_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "CreateShareLinkByPath_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "DeleteFile_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "ExtractFolder_V3",
							"behavior": "Allow"
						},
						{
							"actionId": "GetFileMetadata_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "GetFileMetadataByPath_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "GetAccessPolicies_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "GetFileContent_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "GetFileContentByPath_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "ListFolder_V4",
							"behavior": "Allow"
						},
						{
							"actionId": "ListRootFolder_V4",
							"behavior": "Allow"
						},
						{
							"actionId": "SetBlobTierByPath_V2",
							"behavior": "Allow"
						},
						{
							"actionId": "UpdateFile_V2",
							"behavior": "Allow"
						}
					],
					"defaultConnectorActionRuleBehavior": "Block"
				}
			],
			"endpointConfigurations": [
				{
					"connectorId": "/providers/Microsoft.PowerApps/apis/shared_azureblob",
					"endpointRules": [
						{
							"order": 1,
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
					"customConnectorRuleClassification": "Confidential",
					"pattern": "http://aaa.com"
				},
				{
					"order": 2,
					"customConnectorRuleClassification": "Ignore",
					"pattern": "*"
				}
			]
		}
	}`

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, policiesResponose), nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf(`https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/%s`, policyId1),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, policy1Response), nil
		})

	httpmock.RegisterResponder("GET", fmt.Sprintf(`https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/%s`, policyId2),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, policy2Response), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UnitTestsProviderConfig + `
				data "powerplatform_data_loss_prevention_policies" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					//Verify returned count
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.#", "2"),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.blocked_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.non_business_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.created_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.created_time", "2023-10-02T07:38:50.3269899Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.default_connectors_classification", "General"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.id", policyId1),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.display_name", "a1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.last_modified_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.last_modified_time", "2023-10-02T07:38:50.3269899Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.environment_type", "AllEnvironments"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.environments.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.0.data_group", "Ignore"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.0.host_url_pattern", "*"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.custom_connectors_patterns.0.order", "1"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.blocked_connectors.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.business_connectors.#", "4"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.id", "/providers/Microsoft.PowerApps/apis/shared_office365users"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.default_action_rule_behavior", ""),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.action_rules.#", "0"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.1.endpoint_rules.#", "0"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_azureblob"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.default_action_rule_behavior", "Block"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.action_rules.#", "13"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.action_rules.0.action_id", "CreateFile_V2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.action_rules.0.behavior", "Allow"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.0.behavior", "Deny"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.0.endpoint", "*"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.non_business_connectors.0.endpoint_rules.0.order", "1"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.created_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.created_time", "2023-10-02T07:38:56.6864176Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.default_connectors_classification", "General"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.id", policyId2),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.display_name", "a2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.last_modified_by", "admin"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.last_modified_time", "2023-10-02T07:56:43.9700369Z"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.environment_type", "ExceptEnvironments"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.environments.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.environments.0.name", "be0eb809-e58a-ec1b-8fce-ea40b0e53442"),

					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.1.data_group", "NonBusiness"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.1.host_url_pattern", "http://aaa.com"),
					resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.1.custom_connectors_patterns.1.order", "1"),
				),
			},
		},
	})
}
