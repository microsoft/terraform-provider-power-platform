package powerplatform_mocks

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
)

func NewUnitTestsMockBapiClientInterface(t *testing.T) *MockBapiClientInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockBapiClientInterface(ctrl)

	return clientMock
}

func NewUnitTestMockDataverseClientInterface(t *testing.T) *MockDataverseClientInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockDataverseClientInterface(ctrl)

	return clientMock
}

func NewUnitTestMockPowerPlatformClientInterface(t *testing.T) *MockPowerPlatformClientApiInterface {
	ctrl := gomock.NewController(t)
	clientMock := NewMockPowerPlatformClientApiInterface(ctrl)

	return clientMock
}

const (
	oAuthWellKnownResponse = `{"token_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/token","token_endpoint_auth_methods_supported":["client_secret_post","private_key_jwt","client_secret_basic"],"jwks_uri":"https://login.microsoftonline.com/_/discovery/v2.0/keys","response_modes_supported":["query","fragment","form_post"],"subject_types_supported":["pairwise"],"id_token_signing_alg_values_supported":["RS256"],"response_types_supported":["code","id_token","code id_token","id_token token"],"scopes_supported":["openid","profile","email","offline_access"],"issuer":"https://login.microsoftonline.com/_/v2.0","request_uri_parameter_supported":false,"userinfo_endpoint":"https://graph.microsoft.com/oidc/userinfo","authorization_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/authorize","device_authorization_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/devicecode","http_logout_supported":true,"frontchannel_logout_supported":true,"end_session_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/logout","claims_supported":["sub","iss","cloud_instance_name","cloud_instance_host_name","cloud_graph_host_name","msgraph_host","aud","exp","iat","auth_time","acr","nonce","preferred_username","name","tid","ver","at_hash","c_hash","email"],"kerberos_endpoint":"https://login.microsoftonline.com/_/kerberos","tenant_region_scope":"EU","cloud_instance_name":"microsoftonline.com","cloud_graph_host_name":"graph.windows.net","msgraph_host":"graph.microsoft.com","rbac_url":"https://pas.windows.net"}`
	oAuthTokenResponse     = `{
		"token_type": "Bearer",
		"expires_in": 3599,
		"ext_expires_in": 3599,
		"access_token": "eyJ0eXAiOiJKV1Q"
	}`
)

func ActivateOAuthHttpMocks() {

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		println("No HttpMock responder for: " + req.Method + " " + req.URL.String())
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			println("Body:" + string(body))
		}
		return httpmock.NewStringResponse(http.StatusTeapot, ""), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://login\.microsoftonline\.com/*./v2.0/\.well-known/openid-configuration`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, oAuthWellKnownResponse), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://login\.microsoftonline\.com/*./oauth2/v2.0/token`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, oAuthTokenResponse), nil
		})
}

func ActivateEnvironmentHttpMocks(envId string) {
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%[1]s",
				"type": "Microsoft.BusinessAppPlatform/scopes/environments",
				"location": "europe",
				"name": "%[1]s",
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
						"uniqueName": "%[1]s",
						"domainName": "%[1]s",
						"version": "9.2.23092.00206",
						"instanceUrl": "https://%[1]s.crm4.dynamics.com/",
						"instanceApiUrl": "https://%[1]s.api.crm4.dynamics.com",
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
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{
				"value": [
					{
						"id": "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%[1]s",
						"type": "Microsoft.BusinessAppPlatform/scopes/environments",
						"location": "europe",
						"name": "%[1]s",
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
								"domainName": "%[1]s",
								"version": "9.2.23092.00206",
								"instanceUrl": "https://%[1]s.crm4.dynamics.com/",
								"instanceApiUrl": "https://%[1]s.api.crm4.dynamics.com",
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
					}]}`, envId)), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"isocurrencycode": "PLN"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"_basecurrencyid_value": "xyz"
					}]}`), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

}
