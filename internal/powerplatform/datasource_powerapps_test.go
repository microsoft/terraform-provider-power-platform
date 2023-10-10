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

func TestAccPowerAppsDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				data "powerplatform_powerapps" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.name", regexp.MustCompile(powerplatform_helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.environment_name", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.display_name", regexp.MustCompile(powerplatform_helpers.UrlValidStringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`)),
				),
			},
		},
	})
}

func TestUnitPowerAppsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		println(req.Method + " " + req.URL.String())
		return httpmock.NewStringResponse(http.StatusTeapot, ""), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://api\.powerapps\.com/providers/Microsoft\.PowerApps/scopes/admin/environments/(\d+)/apps`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatchAsUint(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"value": [
					{
						"name": "%[1]d",
						"id": "/providers/Microsoft.PowerApps/scopes/admin/environments/%[1]d/apps/3fec9f57-83bc-4fb8-981e-4b6b45aaa2da",
						"type": "Microsoft.PowerApps/scopes/admin/apps",
						"tags": {
							"primaryDeviceWidth": "1366",
							"primaryDeviceHeight": "768",
							"supportsPortrait": "true",
							"supportsLandscape": "true",
							"primaryFormFactor": "Tablet",
							"publisherVersion": "3.23081.15",
							"minimumRequiredApiVersion": "2.2.0",
							"hasComponent": "false",
							"hasUnlockedComponent": "false",
							"isUnifiedRootApp": "false",
							"sienaVersion": "20230927T203137Z-3.23081.15.0",
							"showStatusBar": "false"
						},
						"properties": {
							"appVersion": "2023-09-27T20:31:37Z",
							"lastDraftVersion": "2023-09-27T20:31:37Z",
							"lifeCycleId": "Published",
							"status": "Ready",
							"createdByClientVersion": {
								"major": 3,
								"minor": 23081,
								"build": 15,
								"revision": 0,
								"majorRevision": 0,
								"minorRevision": 0
							},
							"minClientVersion": {
								"major": 3,
								"minor": 23081,
								"build": 15,
								"revision": 0,
								"majorRevision": 0,
								"minorRevision": 0
							},
							"owner": {
								"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
								"displayName": "admin",
								"email": "admin",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "admin"
							},
							"createdBy": {
								"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
								"displayName": "admin",
								"email": "admin",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "admin"
							},
							"lastModifiedBy": {
								"id": "00000000-0000-0000-0000-5157eaa02fcd",
								"displayName": "SYSTEM",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "00000000-0000-0000-0000-5157eaa02fcd"
							},
							"lastPublishedBy": {
								"id": "00000000-0000-0000-0000-5157eaa02fcd",
								"displayName": "SYSTEM",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "00000000-0000-0000-0000-5157eaa02fcd"
							},
							"backgroundColor": "RGBA(0,176,240,1)",
							"displayName": "Overview",
							"description": "",
							"commitMessage": "",
							"publisher": "",
							"createdTime": "2023-09-27T07:08:47.1964785Z",
							"lastModifiedTime": "2023-09-27T20:31:37.2197567Z",
							"lastPublishTime": "2023-09-27T20:31:37Z",
							"sharedGroupsCount": 0,
							"sharedUsersCount": 0,
							"appOpenProtocolUri": "ms-apps:///providers/Microsoft.PowerApps/apps/3fec9f57-83bc-4fb8-981e-4b6b45aaa2da",
							"appOpenUri": "https://apps.powerapps.com/play/e/%[1]d/a/3fec9f57-83bc-4fb8-981e-4b6b45aaa2da?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&hint=3c7206de-f9cd-4179-9604-c7bf733c7b8c&sourcetime=1695846697184",
							"appPlayUri": "https://apps.powerapps.com/play/e/%[1]d/a/3fec9f57-83bc-4fb8-981e-4b6b45aaa2da?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&hint=3c7206de-f9cd-4179-9604-c7bf733c7b8c&sourcetime=1696937557640",
							"appPlayEmbeddedUri": "https://apps.powerapps.com/play/e/%[1]d/a/3fec9f57-83bc-4fb8-981e-4b6b45aaa2da?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&hint=3c7206de-f9cd-4179-9604-c7bf733c7b8c&telemetryLocation=eu&sourcetime=1696937557640",
							"appPlayTeamsUri": "https://apps.powerapps.com/play/e/%[1]d/a/3fec9f57-83bc-4fb8-981e-4b6b45aaa2da?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&source=teamstab&hint=3c7206de-f9cd-4179-9604-c7bf733c7b8c&telemetryLocation=eu&locale={locale}&channelId={channelId}&channelType={channelType}&chatId={chatId}&groupId={groupId}&hostClientType={hostClientType}&isFullScreen={isFullScreen}&entityId={entityId}&subEntityId={subEntityId}&teamId={teamId}&teamType={teamType}&theme={theme}&userTeamRole={userTeamRole}&sourcetime=1696937557640",
							"connectionReferences": {},
							"authorizationReferences": [],
							"databaseReferences": {
								"default.cds": {
									"databaseDetails": {
										"referenceType": "Environmental",
										"environmentName": "default.cds",
										"overrideValues": {
											"status": "NotSpecified"
										},
										"linkedEnvironmentMetadata": {
											"resourceId": "xxx",
											"friendlyName": "displayName",
											"uniqueName": "unq11",
											"domainName": "xxx",
											"version": "9.2.23092.00206",
											"instanceUrl": "https://xxx.crm4.dynamics.com/",
											"instanceApiUrl": "https://xxx.api.crm4.dynamics.com",
											"baseLanguage": 1033,
											"instanceState": "Ready",
											"createdTime": "2023-09-27T07:08:28.957Z",
											"platformSku": "Standard"
										}
									},
									"dataSources": {
										"Entities": {
											"entitySetName": "entities",
											"logicalName": "entity"
										}
									}
								}
							},
							"userAppMetadata": {
								"favorite": "NotSpecified",
								"includeInAppsList": false
							},
							"isFeaturedApp": false,
							"bypassConsent": false,
							"isHeroApp": false,
							"environment": {
								"id": "/providers/Microsoft.PowerApps/environments/%[1]d",
								"name": "%[1]d",
								"location": "europe"
							},
							"almMode": "Solution",
							"performanceOptimizationEnabled": true,
							"unauthenticatedWebPackageHint": "3c7206de-f9cd-4179-9604-c7bf733c7b8c",
							"canConsumeAppPass": true,
							"enableModernRuntimeMode": false,
							"executionRestrictions": {
								"isTeamsOnly": false,
								"dataLossPreventionEvaluationResult": {
									"status": "Compliant",
									"lastEvaluationDate": "2023-09-27T07:09:02.8310948Z",
									"violations": [],
									"violationsByPolicy": [],
									"violationErrorMessage": "The app uses the following connectors: shared_commondataservice."
								}
							},
							"appPlanClassification": "Premium",
							"usesPremiumApi": true,
							"usesOnlyGrandfatheredPremiumApis": false,
							"usesCustomApi": false,
							"usesOnPremiseGateway": false,
							"usesPcfExternalServiceUsage": false,
							"isCustomizable": true
						},
						"logicalName": "cat_overview_3dbf5",
						"appLocation": "europe",
						"isAppComponentLibrary": false,
						"appType": "CustomCanvasPage"
					},
					{
						"name": "123",
						"id": "/providers/Microsoft.PowerApps/scopes/admin/environments/%[1]d/apps/123",
						"type": "Microsoft.PowerApps/scopes/admin/apps",
						"tags": {
							"primaryDeviceWidth": "1366",
							"primaryDeviceHeight": "768",
							"supportsPortrait": "true",
							"supportsLandscape": "true",
							"primaryFormFactor": "Tablet",
							"publisherVersion": "3.23093.9",
							"minimumRequiredApiVersion": "2.2.0",
							"hasComponent": "false",
							"hasUnlockedComponent": "false",
							"isUnifiedRootApp": "false",
							"sienaVersion": "20230927T203135Z-3.23091.14.0",
							"showStatusBar": "false",
							"optimizedForTeamsMeeting": "true"
						},
						"properties": {
							"appVersion": "2023-09-27T20:31:35Z",
							"lastDraftVersion": "2023-09-27T20:31:35Z",
							"lifeCycleId": "Published",
							"status": "Ready",
							"createdByClientVersion": {
								"major": 3,
								"minor": 23091,
								"build": 14,
								"revision": 0,
								"majorRevision": 0,
								"minorRevision": 0
							},
							"minClientVersion": {
								"major": 3,
								"minor": 23091,
								"build": 14,
								"revision": 0,
								"majorRevision": 0,
								"minorRevision": 0
							},
							"owner": {
								"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
								"displayName": "admin",
								"email": "admin",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "admin"
							},
							"createdBy": {
								"id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
								"displayName": "admin",
								"email": "admin",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "admin"
							},
							"lastModifiedBy": {
								"id": "00000000-0000-0000-0000-5157eaa02fcd",
								"displayName": "SYSTEM",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "00000000-0000-0000-0000-5157eaa02fcd"
							},
							"lastPublishedBy": {
								"id": "00000000-0000-0000-0000-5157eaa02fcd",
								"displayName": "SYSTEM",
								"type": "User",
								"tenantId": "1dbbeae5-8fa6-462e-a5a1-9932a520a1dc",
								"userPrincipalName": "00000000-0000-0000-0000-5157eaa02fcd"
							},
							"backgroundColor": "RGBA(0,176,240,1)",
							"displayName": "Dataverse Actions Page",
							"description": "",
							"commitMessage": "",
							"publisher": "",
							"createdTime": "2023-09-27T07:08:47.2791282Z",
							"lastModifiedTime": "2023-09-27T20:31:35.7685201Z",
							"lastPublishTime": "2023-09-27T20:31:35Z",
							"sharedGroupsCount": 0,
							"sharedUsersCount": 0,
							"appOpenProtocolUri": "ms-apps:///providers/Microsoft.PowerApps/apps/123",
							"appOpenUri": "https://apps.powerapps.com/play/e/%[1]d/a/123?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&hint=dd5bc598-756c-4ffd-ab8a-aa8bd2b50aa3&sourcetime=1695846695701",
							"appPlayUri": "https://apps.powerapps.com/play/e/%[1]d/a/123?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&hint=dd5bc598-756c-4ffd-ab8a-aa8bd2b50aa3&sourcetime=1696937557640",
							"appPlayEmbeddedUri": "https://apps.powerapps.com/play/e/%[1]d/a/123?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&hint=dd5bc598-756c-4ffd-ab8a-aa8bd2b50aa3&telemetryLocation=eu&sourcetime=1696937557640",
							"appPlayTeamsUri": "https://apps.powerapps.com/play/e/%[1]d/a/123?tenantId=1dbbeae5-8fa6-462e-a5a1-9932a520a1dc&source=teamstab&hint=dd5bc598-756c-4ffd-ab8a-aa8bd2b50aa3&telemetryLocation=eu&locale={locale}&channelId={channelId}&channelType={channelType}&chatId={chatId}&groupId={groupId}&hostClientType={hostClientType}&isFullScreen={isFullScreen}&entityId={entityId}&subEntityId={subEntityId}&teamId={teamId}&teamType={teamType}&theme={theme}&userTeamRole={userTeamRole}&sourcetime=1696937557640",
							"connectionReferences": {},
							"authorizationReferences": [],
							"databaseReferences": {
								"default.cds": {
									"databaseDetails": {
										"referenceType": "Environmental",
										"environmentName": "default.cds",
										"overrideValues": {
											"status": "NotSpecified"
										},
										"linkedEnvironmentMetadata": {
											"resourceId": "xxx",
											"friendlyName": "displayName",
											"uniqueName": "unq11",
											"domainName": "xxx",
											"version": "9.2.23092.00206",
											"instanceUrl": "https://xxx.crm4.dynamics.com/",
											"instanceApiUrl": "https://xxx.api.crm4.dynamics.com",
											"baseLanguage": 1033,
											"instanceState": "Ready",
											"createdTime": "2023-09-27T07:08:28.957Z",
											"platformSku": "Standard"
										}
									},
									"dataSources": {
										"Solutions": {
											"entitySetName": "solutions",
											"logicalName": "solution"
										}
									}
								}
							},
							"userAppMetadata": {
								"favorite": "NotSpecified",
								"includeInAppsList": false
							},
							"isFeaturedApp": false,
							"bypassConsent": false,
							"isHeroApp": false,
							"environment": {
								"id": "/providers/Microsoft.PowerApps/environments/%[1]d",
								"location": "europe"
							},
							"almMode": "Solution",
							"performanceOptimizationEnabled": true,
							"unauthenticatedWebPackageHint": "dd5bc598-756c-4ffd-ab8a-aa8bd2b50aa3",
							"canConsumeAppPass": true,
							"enableModernRuntimeMode": false,
							"executionRestrictions": {
								"isTeamsOnly": false,
								"dataLossPreventionEvaluationResult": {
									"status": "Compliant",
									"lastEvaluationDate": "2023-09-27T07:09:03.313275Z",
									"violations": [],
									"violationsByPolicy": [],
									"violationErrorMessage": "The app uses the following connectors: shared_commondataservice."
								}
							},
							"appPlanClassification": "Premium",
							"usesPremiumApi": true,
							"usesOnlyGrandfatheredPremiumApis": false,
							"usesCustomApi": false,
							"usesOnPremiseGateway": false,
							"usesPcfExternalServiceUsage": false,
							"isCustomizable": true
						},
						"logicalName": "cat_dataverseactiondetailspage_eec36",
						"appLocation": "europe",
						"isAppComponentLibrary": false,
						"appType": "CustomCanvasPage"
					}
				]
			}`, id)), nil
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
				data "powerplatform_powerapps" "all" {}`,

				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_powerapps.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.#", "4"),

					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.name", "456"),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.environment_name", "456"),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.display_name", "Overview"),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.0.created_time", "2023-09-27T07:08:47.1964785Z"),

					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.2.name", "123"),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.2.environment_name", "123"),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.2.display_name", "Overview"),
					resource.TestCheckResourceAttr("data.powerplatform_powerapps.all", "powerapps.2.created_time", "2023-09-27T07:08:47.1964785Z"),
				),
			},
		},
	})
}
