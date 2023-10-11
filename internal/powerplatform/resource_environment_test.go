package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

func TestAccEnvironmentsResource_Validate_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example2"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					domain									  = "terraformtest2"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example2"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "terraformtest2"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://terraformtest2.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
				),
			},
			{
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example3"
					domain									  = "terraformtest3"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example3"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "terraformtest3"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://terraformtest3.crm4.dynamics.com/"),
				),
			},
		},
	})
}

func TestAccEnvironmentsResource_Validate_Create(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "USD"
					environment_type                          = "Sandbox"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					domain									  = "terraformtest1"
					templates                                 = ["D365_FinOps_Finance"]
					template_metadata						  = "{\"PostProvisioningPackages\": [{ \"applicationUniqueName\": \"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\n \"parameters\": \"DevToolsEnabled=true|DemoDataEnabled=true\"\n }\n ]\n }"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),

					// Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "terraformtest1"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "security_group_id", regexp.MustCompile(powerplatform_helpers.GuidOrEmptyValueRegex)),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://terraformtest1.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "templates", regexp.MustCompile(`D365_FinOps_Finance$`)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "template_metadata", regexp.MustCompile(`{"PostProvisioningPackages": [{ "applicationUniqueName": "msdyn_FinanceAndOperationsProvisioningAppAnchor",\n "parameters": "DevToolsEnabled=true\|DemoDataEnabled=true"\n }\n ]\n }`)),
					resource.TestMatchResourceAttr("powerplatform_environment.development", "linked_app_url", regexp.MustCompile(`\.operations\.dynamics\.com$`)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsResource_Validate_Create_And_Force_Recreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	//mock_helpers.ActivateEnvironmentHttpMocks("00000000-0000-0000-0000-000000000001")

	envIdResponseInx := -1
	envIdResponseArray := []string{"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
		"00000000-0000-0000-0000-000000000003",
		"00000000-0000-0000-0000-000000000004",
		"00000000-0000-0000-0000-000000000005"}

	envPropertiesMap := make(map[string]map[string]string)
	envPropertiesMap["00000000-0000-0000-0000-000000000001"] = make(map[string]string)
	envPropertiesMap["00000000-0000-0000-0000-000000000001"]["currency"] = "PLN"
	envPropertiesMap["00000000-0000-0000-0000-000000000002"] = make(map[string]string)
	envPropertiesMap["00000000-0000-0000-0000-000000000002"]["currency"] = "PLN"
	envPropertiesMap["00000000-0000-0000-0000-000000000003"] = make(map[string]string)
	envPropertiesMap["00000000-0000-0000-0000-000000000003"]["currency"] = "EUR"
	envPropertiesMap["00000000-0000-0000-0000-000000000004"] = make(map[string]string)
	envPropertiesMap["00000000-0000-0000-0000-000000000004"]["currency"] = "EUR"
	envPropertiesMap["00000000-0000-0000-0000-000000000005"] = make(map[string]string)
	envPropertiesMap["00000000-0000-0000-0000-000000000005"]["currency"] = "EUR"

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"value": [
					{
						"isocurrencycode": "%s"
					}]}`, envPropertiesMap[id]["currency"])), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"_basecurrencyid_value": "xyz"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			envIdResponseInx++

			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
			"id": "b03e1e6d-73db-4367-90e1-2e378bf7e2fc",
			"links": {
				"self": {
					"path": "/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc"
				},
				"environment": {
					"path": "/providers/Microsoft.BusinessAppPlatform/environments/%s"
				}
			},
			"type": {
				"id": "Create"
			},
			"typeDisplayName": "Create",
			"state": {
				"id": "Succeeded"
			},
			"createdDateTime": "2023-10-11T07:45:25.3761337Z",
			"lastActionDateTime": "2023-10-11T07:45:43.4915067Z",
			"requestedBy": {
				"id": "8784d9fb-deb0-4811-96ce-fbf21cf3a1fc",
				"displayName": "ServicePrincipal",
				"type": "ServicePrincipal",
				"tenantId": "123"
			},
			"stages": [
				{
					"id": "Validate",
					"name": "Validate",
					"state": {
						"id": "Succeeded"
					},
					"firstActionDateTime": "2023-10-11T07:45:25.9230185Z",
					"lastActionDateTime": "2023-10-11T07:45:25.9230185Z"
				},
				{
					"id": "Prepare",
					"name": "Prepare",
					"state": {
						"id": "Succeeded"
					},
					"firstActionDateTime": "2023-10-11T07:45:25.9230185Z",
					"lastActionDateTime": "2023-10-11T07:45:25.9230185Z"
				},
				{
					"id": "Run",
					"name": "Run",
					"state": {
						"id": "Succeeded"
					},
					"firstActionDateTime": "2023-10-11T07:45:26.0011473Z",
					"lastActionDateTime": "2023-10-11T07:45:33.2570938Z"
				},
				{
					"id": "Finalize",
					"name": "Finalize",
					"state": {
						"id": "Succeeded"
					},
					"firstActionDateTime": "2023-10-11T07:45:33.3352196Z",
					"lastActionDateTime": "2023-10-11T07:45:43.4915067Z"
				}
			]
		}`, envIdResponseArray[envIdResponseInx])), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001",
				"type": "Microsoft.BusinessAppPlatform/scopes/environments",
				"location": "europe",
				"name": "00000000-0000-0000-0000-000000000001",
				"properties": {
					"tenantId": "123",
					"azureRegion": "westeurope",
					"displayName": "displayname",
					"createdTime": "2023-09-27T07:08:27.6057592Z",
					"createdBy": {
						"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
						"displayName": "admin",
						"email": "admin",
						"type": "User",
						"tenantId": "123",
						"userPrincipalName": "admin"
					},
					"lastModifiedTime": "2023-09-27T07:08:34.9205145Z",
					"provisioningState": "Succeeded",
					"creationType": "User",
					"environmentSku": "Sandbox",
					"isDefault": false,
					"capacity": [
						{
							"capacityType": "Database",
							"actualConsumption": 885.0391,
							"ratedConsumption": 1024.0,
							"capacityUnit": "MB",
							"updatedOn": "2023-10-10T03:00:35Z"
						},
						{
							"capacityType": "File",
							"actualConsumption": 1187.142,
							"ratedConsumption": 1187.142,
							"capacityUnit": "MB",
							"updatedOn": "2023-10-10T03:00:35Z"
						},
						{
							"capacityType": "Log",
							"actualConsumption": 0.0,
							"ratedConsumption": 0.0,
							"capacityUnit": "MB",
							"updatedOn": "2023-10-10T03:00:35Z"
						},
						{
							"capacityType": "FinOpsDatabase",
							"actualConsumption": 0.0,
							"ratedConsumption": 0.0,
							"capacityUnit": "MB",
							"updatedOn": "2023-10-10T03:00:35Z"
						},
						{
							"capacityType": "FinOpsFile",
							"actualConsumption": 0.0,
							"ratedConsumption": 0.0,
							"capacityUnit": "MB",
							"updatedOn": "2023-10-10T03:00:35Z"
						}
					],
					"addons": [],
					"clientUris": {
						"admin": "https://admin.powerplatform.microsoft.com/environments/environment/456/hub",
						"maker": "https://make.powerapps.com/environments/456/home"
					},
					"runtimeEndpoints": {
						"microsoft.BusinessAppPlatform": "https://europe.api.bap.microsoft.com",
						"microsoft.CommonDataModel": "https://europe.api.cds.microsoft.com",
						"microsoft.PowerApps": "https://europe.api.powerapps.com",
						"microsoft.PowerAppsAdvisor": "https://europe.api.advisor.powerapps.com",
						"microsoft.PowerVirtualAgents": "https://powervamg.eu-il107.gateway.prod.island.powerapps.com",
						"microsoft.ApiManagement": "https://management.EUROPE.azure-apihub.net",
						"microsoft.Flow": "https://emea.api.flow.microsoft.com"
					},
					"databaseType": "CommonDataService",
					"linkedEnvironmentMetadata": {
						"resourceId": "orgid",
						"friendlyName": "displayname",
						"uniqueName": "00000000-0000-0000-0000-000000000001",
						"domainName": "00000000-0000-0000-0000-000000000001",
						"version": "9.2.23092.00206",
						"instanceUrl": "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/",
						"instanceApiUrl": "https://00000000-0000-0000-0000-000000000001.api.crm4.dynamics.com",
						"baseLanguage": 1033,
						"instanceState": "Ready",
						"createdTime": "2023-09-27T07:08:28.957Z",
						"backgroundOperationsState": "Enabled",
						"scaleGroup": "EURCRMLIVESG705",
						"platformSku": "Standard",
						"schemaType": "Standard"
					},
					"trialScenarioType": "None",
					"notificationMetadata": {
						"state": "NotSpecified",
						"branding": "NotSpecific"
					},
					"retentionPeriod": "P7D",
					"states": {
						"management": {
							"id": "Ready"
						},
						"runtime": {
							"runtimeReasonCode": "NotSpecified",
							"requestedBy": {
								"displayName": "SYSTEM",
								"type": "NotSpecified"
							},
							"id": "Enabled"
						}
					},
					"updateCadence": {
						"id": "Frequent"
					},
					"retentionDetails": {
						"retentionPeriod": "P7D",
						"backupsAvailableFromDateTime": "2023-10-03T09:23:06.1717665Z"
					},
					"protectionStatus": {
						"keyManagedBy": "Microsoft"
					},
					"cluster": {
						"category": "Prod",
						"number": "107",
						"uriSuffix": "eu-il107.gateway.prod.island",
						"geoShortName": "EU",
						"environment": "Prod"
					},
					"connectedGroups": [],
					"lifecycleOperationsEnforcement": {
						"allowedOperations": [
							{
								"type": {
									"id": "DisableGovernanceConfiguration"
								},
								"reason": {
									"message": "DisableGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
									"type": "GovernanceConfig"
								}
							},
							{
								"type": {
									"id": "UpdateGovernanceConfiguration"
								},
								"reason": {
									"message": "UpdateGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
									"type": "GovernanceConfig"
								}
							}
						]
					},
					"governanceConfiguration": {
						"protectionLevel": "Basic"
					}
				}
			}`), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000002`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
			"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000002",
			"type": "Microsoft.BusinessAppPlatform/scopes/environments",
			"location": "unitedstates",
			"name": "00000000-0000-0000-0000-000000000002",
			"properties": {
				"tenantId": "123",
				"azureRegion": "westeurope",
				"displayName": "displayname",
				"createdTime": "2023-09-27T07:08:27.6057592Z",
				"createdBy": {
					"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
					"displayName": "admin",
					"email": "admin",
					"type": "User",
					"tenantId": "123",
					"userPrincipalName": "admin"
				},
				"lastModifiedTime": "2023-09-27T07:08:34.9205145Z",
				"provisioningState": "Succeeded",
				"creationType": "User",
				"environmentSku": "Sandbox",
				"isDefault": false,
				"capacity": [
					{
						"capacityType": "Database",
						"actualConsumption": 885.0391,
						"ratedConsumption": 1024.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "File",
						"actualConsumption": 1187.142,
						"ratedConsumption": 1187.142,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "Log",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsDatabase",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsFile",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					}
				],
				"addons": [],
				"clientUris": {
					"admin": "https://admin.powerplatform.microsoft.com/environments/environment/00000000-0000-0000-0000-000000000002/hub",
					"maker": "https://make.powerapps.com/environments/00000000-0000-0000-0000-000000000002/home"
				},
				"runtimeEndpoints": {
					"microsoft.BusinessAppPlatform": "https://europe.api.bap.microsoft.com",
					"microsoft.CommonDataModel": "https://europe.api.cds.microsoft.com",
					"microsoft.PowerApps": "https://europe.api.powerapps.com",
					"microsoft.PowerAppsAdvisor": "https://europe.api.advisor.powerapps.com",
					"microsoft.PowerVirtualAgents": "https://powervamg.eu-il107.gateway.prod.island.powerapps.com",
					"microsoft.ApiManagement": "https://management.EUROPE.azure-apihub.net",
					"microsoft.Flow": "https://emea.api.flow.microsoft.com"
				},
				"databaseType": "CommonDataService",
				"linkedEnvironmentMetadata": {
					"resourceId": "orgid",
					"friendlyName": "displayname",
					"uniqueName": "00000000-0000-0000-0000-000000000002",
					"domainName": "00000000-0000-0000-0000-000000000002",
					"version": "9.2.23092.00206",
					"instanceUrl": "https://00000000-0000-0000-0000-000000000002.crm4.dynamics.com/",
					"instanceApiUrl": "https://00000000-0000-0000-0000-000000000002.api.crm4.dynamics.com",
					"baseLanguage": 1033,
					"instanceState": "Ready",
					"createdTime": "2023-09-27T07:08:28.957Z",
					"backgroundOperationsState": "Enabled",
					"scaleGroup": "EURCRMLIVESG705",
					"platformSku": "Standard",
					"schemaType": "Standard"
				},
				"trialScenarioType": "None",
				"notificationMetadata": {
					"state": "NotSpecified",
					"branding": "NotSpecific"
				},
				"retentionPeriod": "P7D",
				"states": {
					"management": {
						"id": "Ready"
					},
					"runtime": {
						"runtimeReasonCode": "NotSpecified",
						"requestedBy": {
							"displayName": "SYSTEM",
							"type": "NotSpecified"
						},
						"id": "Enabled"
					}
				},
				"updateCadence": {
					"id": "Frequent"
				},
				"retentionDetails": {
					"retentionPeriod": "P7D",
					"backupsAvailableFromDateTime": "2023-10-03T09:23:06.1717665Z"
				},
				"protectionStatus": {
					"keyManagedBy": "Microsoft"
				},
				"cluster": {
					"category": "Prod",
					"number": "107",
					"uriSuffix": "eu-il107.gateway.prod.island",
					"geoShortName": "EU",
					"environment": "Prod"
				},
				"connectedGroups": [],
				"lifecycleOperationsEnforcement": {
					"allowedOperations": [
						{
							"type": {
								"id": "DisableGovernanceConfiguration"
							},
							"reason": {
								"message": "DisableGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						},
						{
							"type": {
								"id": "UpdateGovernanceConfiguration"
							},
							"reason": {
								"message": "UpdateGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						}
					]
				},
				"governanceConfiguration": {
					"protectionLevel": "Basic"
				}
			}
		}`), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000003`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
			"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000003",
			"type": "Microsoft.BusinessAppPlatform/scopes/environments",
			"location": "unitedstates",
			"name": "00000000-0000-0000-0000-000000000003",
			"properties": {
				"tenantId": "123",
				"azureRegion": "westeurope",
				"displayName": "Example1",
				"createdTime": "2023-09-27T07:08:27.6057592Z",
				"createdBy": {
					"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
					"displayName": "admin",
					"email": "admin",
					"type": "User",
					"tenantId": "123",
					"userPrincipalName": "admin"
				},
				"lastModifiedTime": "2023-09-27T07:08:34.9205145Z",
				"provisioningState": "Succeeded",
				"creationType": "User",
				"environmentSku": "Sandbox",
				"isDefault": false,
				"capacity": [
					{
						"capacityType": "Database",
						"actualConsumption": 885.0391,
						"ratedConsumption": 1024.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "File",
						"actualConsumption": 1187.142,
						"ratedConsumption": 1187.142,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "Log",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsDatabase",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsFile",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					}
				],
				"addons": [],
				"clientUris": {
					"admin": "https://admin.powerplatform.microsoft.com/environments/environment/00000000-0000-0000-0000-000000000003/hub",
					"maker": "https://make.powerapps.com/environments/00000000-0000-0000-0000-000000000003/home"
				},
				"runtimeEndpoints": {
					"microsoft.BusinessAppPlatform": "https://europe.api.bap.microsoft.com",
					"microsoft.CommonDataModel": "https://europe.api.cds.microsoft.com",
					"microsoft.PowerApps": "https://europe.api.powerapps.com",
					"microsoft.PowerAppsAdvisor": "https://europe.api.advisor.powerapps.com",
					"microsoft.PowerVirtualAgents": "https://powervamg.eu-il107.gateway.prod.island.powerapps.com",
					"microsoft.ApiManagement": "https://management.EUROPE.azure-apihub.net",
					"microsoft.Flow": "https://emea.api.flow.microsoft.com"
				},
				"databaseType": "CommonDataService",
				"linkedEnvironmentMetadata": {
					"resourceId": "orgid",
					"friendlyName": "displayname",
					"uniqueName": "00000000-0000-0000-0000-000000000003",
					"domainName": "00000000-0000-0000-0000-000000000003",
					"version": "9.2.23092.00206",
					"instanceUrl": "https://00000000-0000-0000-0000-000000000003.crm4.dynamics.com/",
					"instanceApiUrl": "https://00000000-0000-0000-0000-000000000003.api.crm4.dynamics.com",
					"baseLanguage": 1033,
					"instanceState": "Ready",
					"createdTime": "2023-09-27T07:08:28.957Z",
					"backgroundOperationsState": "Enabled",
					"scaleGroup": "EURCRMLIVESG705",
					"platformSku": "Standard",
					"schemaType": "Standard"
				},
				"trialScenarioType": "None",
				"notificationMetadata": {
					"state": "NotSpecified",
					"branding": "NotSpecific"
				},
				"retentionPeriod": "P7D",
				"states": {
					"management": {
						"id": "Ready"
					},
					"runtime": {
						"runtimeReasonCode": "NotSpecified",
						"requestedBy": {
							"displayName": "SYSTEM",
							"type": "NotSpecified"
						},
						"id": "Enabled"
					}
				},
				"updateCadence": {
					"id": "Frequent"
				},
				"retentionDetails": {
					"retentionPeriod": "P7D",
					"backupsAvailableFromDateTime": "2023-10-03T09:23:06.1717665Z"
				},
				"protectionStatus": {
					"keyManagedBy": "Microsoft"
				},
				"cluster": {
					"category": "Prod",
					"number": "107",
					"uriSuffix": "eu-il107.gateway.prod.island",
					"geoShortName": "EU",
					"environment": "Prod"
				},
				"connectedGroups": [],
				"lifecycleOperationsEnforcement": {
					"allowedOperations": [
						{
							"type": {
								"id": "DisableGovernanceConfiguration"
							},
							"reason": {
								"message": "DisableGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						},
						{
							"type": {
								"id": "UpdateGovernanceConfiguration"
							},
							"reason": {
								"message": "UpdateGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						}
					]
				},
				"governanceConfiguration": {
					"protectionLevel": "Basic"
				}
			}
		}`), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000004`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
			"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000004",
			"type": "Microsoft.BusinessAppPlatform/scopes/environments",
			"location": "unitedstates",
			"name": "00000000-0000-0000-0000-000000000004",
			"properties": {
				"tenantId": "123",
				"azureRegion": "westeurope",
				"displayName": "Example1",
				"createdTime": "2023-09-27T07:08:27.6057592Z",
				"createdBy": {
					"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
					"displayName": "admin",
					"email": "admin",
					"type": "User",
					"tenantId": "123",
					"userPrincipalName": "admin"
				},
				"lastModifiedTime": "2023-09-27T07:08:34.9205145Z",
				"provisioningState": "Succeeded",
				"creationType": "User",
				"environmentSku": "Trial",
				"isDefault": false,
				"capacity": [
					{
						"capacityType": "Database",
						"actualConsumption": 885.0391,
						"ratedConsumption": 1024.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "File",
						"actualConsumption": 1187.142,
						"ratedConsumption": 1187.142,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "Log",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsDatabase",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsFile",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					}
				],
				"addons": [],
				"clientUris": {
					"admin": "https://admin.powerplatform.microsoft.com/environments/environment/00000000-0000-0000-0000-000000000004/hub",
					"maker": "https://make.powerapps.com/environments/00000000-0000-0000-0000-000000000004/home"
				},
				"runtimeEndpoints": {
					"microsoft.BusinessAppPlatform": "https://europe.api.bap.microsoft.com",
					"microsoft.CommonDataModel": "https://europe.api.cds.microsoft.com",
					"microsoft.PowerApps": "https://europe.api.powerapps.com",
					"microsoft.PowerAppsAdvisor": "https://europe.api.advisor.powerapps.com",
					"microsoft.PowerVirtualAgents": "https://powervamg.eu-il107.gateway.prod.island.powerapps.com",
					"microsoft.ApiManagement": "https://management.EUROPE.azure-apihub.net",
					"microsoft.Flow": "https://emea.api.flow.microsoft.com"
				},
				"databaseType": "CommonDataService",
				"linkedEnvironmentMetadata": {
					"resourceId": "orgid",
					"friendlyName": "displayname",
					"uniqueName": "00000000-0000-0000-0000-000000000004",
					"domainName": "00000000-0000-0000-0000-000000000004",
					"version": "9.2.23092.00206",
					"instanceUrl": "https://00000000-0000-0000-0000-000000000004.crm4.dynamics.com/",
					"instanceApiUrl": "https://00000000-0000-0000-0000-000000000004.api.crm4.dynamics.com",
					"baseLanguage": 1033,
					"instanceState": "Ready",
					"createdTime": "2023-09-27T07:08:28.957Z",
					"backgroundOperationsState": "Enabled",
					"scaleGroup": "EURCRMLIVESG705",
					"platformSku": "Standard",
					"schemaType": "Standard"
				},
				"trialScenarioType": "None",
				"notificationMetadata": {
					"state": "NotSpecified",
					"branding": "NotSpecific"
				},
				"retentionPeriod": "P7D",
				"states": {
					"management": {
						"id": "Ready"
					},
					"runtime": {
						"runtimeReasonCode": "NotSpecified",
						"requestedBy": {
							"displayName": "SYSTEM",
							"type": "NotSpecified"
						},
						"id": "Enabled"
					}
				},
				"updateCadence": {
					"id": "Frequent"
				},
				"retentionDetails": {
					"retentionPeriod": "P7D",
					"backupsAvailableFromDateTime": "2023-10-03T09:23:06.1717665Z"
				},
				"protectionStatus": {
					"keyManagedBy": "Microsoft"
				},
				"cluster": {
					"category": "Prod",
					"number": "107",
					"uriSuffix": "eu-il107.gateway.prod.island",
					"geoShortName": "EU",
					"environment": "Prod"
				},
				"connectedGroups": [],
				"lifecycleOperationsEnforcement": {
					"allowedOperations": [
						{
							"type": {
								"id": "DisableGovernanceConfiguration"
							},
							"reason": {
								"message": "DisableGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						},
						{
							"type": {
								"id": "UpdateGovernanceConfiguration"
							},
							"reason": {
								"message": "UpdateGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						}
					]
				},
				"governanceConfiguration": {
					"protectionLevel": "Basic"
				}
			}
		}`), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000005`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
			"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000005",
			"type": "Microsoft.BusinessAppPlatform/scopes/environments",
			"location": "unitedstates",
			"name": "00000000-0000-0000-0000-000000000005",
			"properties": {
				"tenantId": "123",
				"azureRegion": "westeurope",
				"displayName": "Example1",
				"createdTime": "2023-09-27T07:08:27.6057592Z",
				"createdBy": {
					"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
					"displayName": "admin",
					"email": "admin",
					"type": "User",
					"tenantId": "123",
					"userPrincipalName": "admin"
				},
				"lastModifiedTime": "2023-09-27T07:08:34.9205145Z",
				"provisioningState": "Succeeded",
				"creationType": "User",
				"environmentSku": "Trial",
				"isDefault": false,
				"capacity": [
					{
						"capacityType": "Database",
						"actualConsumption": 885.0391,
						"ratedConsumption": 1024.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "File",
						"actualConsumption": 1187.142,
						"ratedConsumption": 1187.142,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "Log",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsDatabase",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					},
					{
						"capacityType": "FinOpsFile",
						"actualConsumption": 0.0,
						"ratedConsumption": 0.0,
						"capacityUnit": "MB",
						"updatedOn": "2023-10-10T03:00:35Z"
					}
				],
				"addons": [],
				"clientUris": {
					"admin": "https://admin.powerplatform.microsoft.com/environments/environment/00000000-0000-0000-0000-000000000005/hub",
					"maker": "https://make.powerapps.com/environments/00000000-0000-0000-0000-000000000005/home"
				},
				"runtimeEndpoints": {
					"microsoft.BusinessAppPlatform": "https://europe.api.bap.microsoft.com",
					"microsoft.CommonDataModel": "https://europe.api.cds.microsoft.com",
					"microsoft.PowerApps": "https://europe.api.powerapps.com",
					"microsoft.PowerAppsAdvisor": "https://europe.api.advisor.powerapps.com",
					"microsoft.PowerVirtualAgents": "https://powervamg.eu-il107.gateway.prod.island.powerapps.com",
					"microsoft.ApiManagement": "https://management.EUROPE.azure-apihub.net",
					"microsoft.Flow": "https://emea.api.flow.microsoft.com"
				},
				"databaseType": "CommonDataService",
				"linkedEnvironmentMetadata": {
					"resourceId": "orgid",
					"friendlyName": "displayname",
					"uniqueName": "00000000-0000-0000-0000-000000000005",
					"domainName": "00000000-0000-0000-0000-000000000005",
					"version": "9.2.23092.00206",
					"instanceUrl": "https://00000000-0000-0000-0000-000000000005.crm4.dynamics.com/",
					"instanceApiUrl": "https://00000000-0000-0000-0000-000000000005.api.crm4.dynamics.com",
					"baseLanguage": 1031,
					"instanceState": "Ready",
					"createdTime": "2023-09-27T07:08:28.957Z",
					"backgroundOperationsState": "Enabled",
					"scaleGroup": "EURCRMLIVESG705",
					"platformSku": "Standard",
					"schemaType": "Standard"
				},
				"trialScenarioType": "None",
				"notificationMetadata": {
					"state": "NotSpecified",
					"branding": "NotSpecific"
				},
				"retentionPeriod": "P7D",
				"states": {
					"management": {
						"id": "Ready"
					},
					"runtime": {
						"runtimeReasonCode": "NotSpecified",
						"requestedBy": {
							"displayName": "SYSTEM",
							"type": "NotSpecified"
						},
						"id": "Enabled"
					}
				},
				"updateCadence": {
					"id": "Frequent"
				},
				"retentionDetails": {
					"retentionPeriod": "P7D",
					"backupsAvailableFromDateTime": "2023-10-03T09:23:06.1717665Z"
				},
				"protectionStatus": {
					"keyManagedBy": "Microsoft"
				},
				"cluster": {
					"category": "Prod",
					"number": "107",
					"uriSuffix": "eu-il107.gateway.prod.island",
					"geoShortName": "EU",
					"environment": "Prod"
				},
				"connectedGroups": [],
				"lifecycleOperationsEnforcement": {
					"allowedOperations": [
						{
							"type": {
								"id": "DisableGovernanceConfiguration"
							},
							"reason": {
								"message": "DisableGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						},
						{
							"type": {
								"id": "UpdateGovernanceConfiguration"
							},
							"reason": {
								"message": "UpdateGovernanceConfiguration cannot be performed on Power Platform environment because of the governance configuration.",
								"type": "GovernanceConfig"
							}
						}
					]
				},
				"governanceConfiguration": {
					"protectionLevel": "Basic"
				}
			}
		}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000001"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
	
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "PLN"),
				),
			},
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "unitedstates"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000002"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "unitedstates"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "PLN"),
				),
			},
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					language_code                             = "1033"
					currency_code                             = "EUR"
					environment_type                          = "Sandbox"
					domain									  = "00000000-0000-0000-0000-000000000003"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
				),
			},
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					language_code                             = "1033"
					currency_code                             = "EUR"
					environment_type                          = "Trial"
					domain									  = "00000000-0000-0000-0000-000000000004"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Trial"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
				),
			},
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "Example1"
					location                                  = "unitedstates"
					language_code                             = "1031"
					currency_code                             = "EUR"
					environment_type                          = "Trial"
					domain									  = "00000000-0000-0000-0000-000000000005"
					security_group_id 						  = "00000000-0000-0000-0000-000000000000"
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", "00000000-0000-0000-0000-000000000005"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1031"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "EUR"),
				),
			},
		},
	})

}

func TestUnitEnvironmentsResource_Validate_Create_And_Update(t *testing.T) {
	clientMock := mocks.NewUnitTestsMockBapiClientInterface(t)
	dataverseClientMock := mocks.NewUnitTestMockDataverseClientInterface(t)

	envId := "00000000-0000-0000-0000-000000000001"
	env := models.EnvironmentDto{
		Name: envId,
		Properties: models.EnvironmentPropertiesDto{
			EnvironmentSku: "Sandbox",
			LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
				ResourceId:      "org1",
				SecurityGroupId: "security1",
				DomainName:      "domain",
				InstanceURL:     "url",
				Version:         "version",
			},
			LinkedAppMetadata: models.LinkedAppMetadataDto{
				Type: "Internal",
				Id:   "00000000-0000-0000-0000-000000000000",
				Url:  "https://url.operations.dynamics.com",
			},
		},
	}

	steps := []resource.TestStep{
		{
			Config: UniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example1"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example1"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security1"),
			),
		},
		{
			Config: UniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example123"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security1"),
			),
		},
		{
			Config: UniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example123"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain123"
				security_group_id 						  = "security1"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security1"),
			),
		},
		{
			Config: UniTestsProviderConfig + `
			resource "powerplatform_environment" "development" {
				display_name                              = "Example123"
				location                                  = "europe"
				language_code                             = "1033"
				currency_code                             = "USD"
				environment_type                          = "Sandbox"
				domain									  = "domain123"
				security_group_id 						  = "security123"
			}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", envId),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "Example123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "domain123"),
				resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "security123"),
			),
		},
	}

	clientMock.EXPECT().GetEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) (*models.EnvironmentDto, error) {
		return &env, nil
	}).AnyTimes()

	dataverseClientMock.EXPECT().GetDefaultCurrencyForEnvironment(gomock.Any(), gomock.Any()).Return(&models.TransactionCurrencyDto{IsoCurrencyCode: "USD"}, nil).AnyTimes()

	clientMock.EXPECT().CreateEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, envToCreate models.EnvironmentCreateDto) (*models.EnvironmentDto, error) {
		env = models.EnvironmentDto{
			Id:       envId,
			Location: envToCreate.Location,
			Name:     envId,
			Properties: models.EnvironmentPropertiesDto{
				DisplayName:    envToCreate.Properties.DisplayName,
				EnvironmentSku: env.Properties.EnvironmentSku,
				LinkedEnvironmentMetadata: models.LinkedEnvironmentMetadataDto{
					DomainName:      "domain",
					InstanceURL:     "url",
					BaseLanguage:    envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage,
					SecurityGroupId: envToCreate.Properties.LinkedEnvironmentMetadata.SecurityGroupId,
					Version:         "version",
					ResourceId:      "org1",
				},
			},
		}
		return &env, nil
	}).Times(1)

	clientMock.EXPECT().UpdateEnvironment(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error) {
		env.Name = environment.Name
		env.Id = environment.Id
		env.Properties.DisplayName = environment.Properties.DisplayName
		env.Properties.LinkedEnvironmentMetadata.DomainName = environment.Properties.LinkedEnvironmentMetadata.DomainName
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = environment.Properties.LinkedEnvironmentMetadata.SecurityGroupId
		return &env, nil
	}).Times(len(steps) - 1)

	clientMock.EXPECT().DeleteEnvironment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id string) error {
		return nil
	}).AnyTimes()

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"powerplatform": powerPlatformProviderServerApiMock(clientMock, dataverseClientMock, nil),
		},
		Steps: steps,
	})

}

func TestUnitEnvironmentsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()
	mock_helpers.ActivateEnvironmentHttpMocks("00000000-0000-0000-0000-000000000001")

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"id": "b03e1e6d-73db-4367-90e1-2e378bf7e2fc",
				"links": {
					"self": {
						"path": "/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc"
					},
					"environment": {
						"path": "/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001"
					}
				},
				"type": {
					"id": "Create"
				},
				"typeDisplayName": "Create",
				"state": {
					"id": "Succeeded"
				},
				"createdDateTime": "2023-10-11T07:45:25.3761337Z",
				"lastActionDateTime": "2023-10-11T07:45:43.4915067Z",
				"requestedBy": {
					"id": "8784d9fb-deb0-4811-96ce-fbf21cf3a1fc",
					"displayName": "ServicePrincipal",
					"type": "ServicePrincipal",
					"tenantId": "123"
				},
				"stages": [
					{
						"id": "Validate",
						"name": "Validate",
						"state": {
							"id": "Succeeded"
						},
						"firstActionDateTime": "2023-10-11T07:45:25.9230185Z",
						"lastActionDateTime": "2023-10-11T07:45:25.9230185Z"
					},
					{
						"id": "Prepare",
						"name": "Prepare",
						"state": {
							"id": "Succeeded"
						},
						"firstActionDateTime": "2023-10-11T07:45:25.9230185Z",
						"lastActionDateTime": "2023-10-11T07:45:25.9230185Z"
					},
					{
						"id": "Run",
						"name": "Run",
						"state": {
							"id": "Succeeded"
						},
						"firstActionDateTime": "2023-10-11T07:45:26.0011473Z",
						"lastActionDateTime": "2023-10-11T07:45:33.2570938Z"
					},
					{
						"id": "Finalize",
						"name": "Finalize",
						"state": {
							"id": "Succeeded"
						},
						"firstActionDateTime": "2023-10-11T07:45:33.3352196Z",
						"lastActionDateTime": "2023-10-11T07:45:43.4915067Z"
					}
				]
			}`), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				resource "powerplatform_environment" "development" {
					display_name                              = "displayname"
					location                                  = "europe"
					language_code                             = "1033"
					currency_code                             = "PLN"
					environment_type                          = "Sandbox"
					domain                                    = "00000000-0000-0000-0000-000000000001"
					security_group_id                         = "00000000-0000-0000-0000-000000000000"
				}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_name", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "display_name", "displayname"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "url", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "domain", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "language_code", "1033"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "currency_code", "PLN"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "organization_id", "orgid"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "security_group_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("powerplatform_environment.development", "version", "9.2.23092.00206"),
				),
			},
		},
	})

}
