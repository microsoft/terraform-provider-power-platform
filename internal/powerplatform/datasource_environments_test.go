package powerplatform

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccEnvironmentsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.domain", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.environment_name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.environment_type", regexp.MustCompile(`^(Default|Sandbox|Developer)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.language_code", regexp.MustCompile(`^(1033|1031)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.organization_id", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.security_group_id", regexp.MustCompile(powerplatform_helpers.GuidOrEmptyValueRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.url", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.location", regexp.MustCompile(`^(unitedstates|europe)$`)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.version", regexp.MustCompile(powerplatform_helpers.VersionRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "environments.0.currency_code", regexp.MustCompile(powerplatform_helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		println(req.Method + " " + req.URL.String())
		return httpmock.NewStringResponse(http.StatusTeapot, ""), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://(\d+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"isocurrencycode": "PLN"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://(\d+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"_basecurrencyid_value": "xyz"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/(\d+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatchAsUint(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/456",
				"type": "Microsoft.BusinessAppPlatform/scopes/environments",
				"location": "europe",
				"name": "%[1]d",
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
						"uniqueName": "%[1]d",
						"domainName": "%[1]d",
						"version": "9.2.23092.00206",
						"instanceUrl": "https://%[1]d.crm4.dynamics.com/",
						"instanceApiUrl": "https://%[1]d.api.crm4.dynamics.com",
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
			}`, id)), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/456",
						"type": "Microsoft.BusinessAppPlatform/scopes/environments",
						"location": "europe",
						"name": "456",
						"properties": {
							"tenantId": "123",
							"azureRegion": "northeurope",
							"displayName": "Admin AdminOnMicrosoft's Environment",
							"createdTime": "2023-02-15T08:02:36.1799125Z",
							"createdBy": {
								"id": "SYSTEM",
								"displayName": "SYSTEM",
								"type": "NotSpecified"
							},
							"usedBy": {
								"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
								"type": "User",
								"tenantId": "123",
								"userPrincipalName": "admin"
							},
							"provisioningState": "Succeeded",
							"creationType": "Developer",
							"environmentSku": "Developer",
							"isDefault": false,
							"clientUris": {
								"admin": "https://admin.powerplatform.microsoft.com/environments/environment/456/hub",
								"maker": "https://make.powerapps.com/environments/456/home"
							},
							"runtimeEndpoints": {
								"microsoft.BusinessAppPlatform": "https://europe.api.bap.microsoft.com",
								"microsoft.CommonDataModel": "https://europe.api.cds.microsoft.com",
								"microsoft.PowerApps": "https://europe.api.powerapps.com",
								"microsoft.PowerAppsAdvisor": "https://europe.api.advisor.powerapps.com",
								"microsoft.PowerVirtualAgents": "https://powervamg.eu-il106.gateway.prod.island.powerapps.com",
								"microsoft.ApiManagement": "https://management.EUROPE.azure-apihub.net",
								"microsoft.Flow": "https://emea.api.flow.microsoft.com"
							},
							"databaseType": "CommonDataService",
							"linkedEnvironmentMetadata": {
								"resourceId": "6450637c-f9a8-4988-8cf7-b03723d51ab7",
								"friendlyName": "Admin AdminOnMicrosoft's Environment",
								"uniqueName": "unq6450637cf9a849888cf7b03723d51",
								"domainName": "yyy",
								"version": "9.2.23092.00206",
								"instanceUrl": "https://yyy.crm4.dynamics.com/",
								"instanceApiUrl": "https://yyy.api.crm4.dynamics.com",
								"baseLanguage": 1033,
								"instanceState": "Ready",
								"createdTime": "2023-02-15T08:02:46.87Z",
								"backgroundOperationsState": "Enabled",
								"scaleGroup": "EURCRMLIVESG633",
								"platformSku": "Standard",
								"schemaType": "Standard"
							},
							"trialScenarioType": "None",
							"retentionPeriod": "P7D",
							"states": {
								"management": {
									"id": "NotSpecified"
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
								"backupsAvailableFromDateTime": "2023-10-03T08:12:55.5332994Z"
							},
							"protectionStatus": {
								"keyManagedBy": "Microsoft"
							},
							"cluster": {
								"category": "Prod",
								"number": "106",
								"uriSuffix": "eu-il106.gateway.prod.island",
								"geoShortName": "EU",
								"environment": "Prod"
							},
							"connectedGroups": [],
							"lifecycleOperationsEnforcement": {
								"allowedOperations": [
									{
										"type": {
											"id": "Move"
										}
									}
								]
							},
							"governanceConfiguration": {
								"protectionLevel": "Basic"
							}
						}
					},
					{
						"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/123",
						"type": "Microsoft.BusinessAppPlatform/scopes/environments",
						"location": "europe",
						"name": "123",
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
							"clientUris": {
								"admin": "https://admin.powerplatform.microsoft.com/environments/environment/xxx/hub",
								"maker": "https://make.powerapps.com/environments/xxx/home"
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
								"uniqueName": "xxx",
								"domainName": "xxx",
								"version": "9.2.23092.00206",
								"instanceUrl": "https://xxx.crm4.dynamics.com/",
								"instanceApiUrl": "https://xxx.api.crm4.dynamics.com",
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
								"backupsAvailableFromDateTime": "2023-10-03T08:12:55.5332994Z"
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
											"id": "Move"
										}
									}
								],
								"disallowedOperations": [
									{
										"type": {
											"id": "Provision"
										},
										"reason": {
											"message": "Provision cannot be performed because there is no linked CDS instance or the CDS instance version is not supported.",
											"type": "CdsLink"
										}
									}
								]
							},
							"governanceConfiguration": {
								"protectionLevel": "Basic"
							}
						}
					}]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				data "powerplatform_environments" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_environments.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.#", "2"),

					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.display_name", "displayname"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.domain", "xxx"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.environment_name", "123"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.environment_type", "Sandbox"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.language_code", "1033"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.organization_id", "orgid"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.security_group_id", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.url", "https://xxx.crm4.dynamics.com/"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.location", "europe"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.version", "9.2.23092.00206"),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.linked_app_type", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.linked_app_id", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.linked_app_url", ""),
					resource.TestCheckResourceAttr("data.powerplatform_environments.all", "environments.1.currency_code", "PLN"),
				),
			},
		},
	})
}
