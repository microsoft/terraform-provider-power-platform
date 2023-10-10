package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccConnectorsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ProviderConfig + `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					// Verify the first power app to ensure all attributes are set
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", regexp.MustCompile(powerplatform_helpers.ApiIdRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", regexp.MustCompile(powerplatform_helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", regexp.MustCompile(powerplatform_helpers.ApiIdRegex)),
				),
			},
		},
	})
}

func TestUnitConnectorsDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock_helpers.ActivateOAuthHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `[
				{
					"id": "Http",
					"metadata": {
						"virtualConnector": true,
						"name": "HTTP",
						"type": "Microsoft.PowerApps/apis",
						"iconUri": "data:image/svg+xml;base64,PCEtLSBQbGVhc2UgYWxzbyBjaGFuZ2UgdGhpcyBpbiBBenVyZS1CUE1VWC4gVGhpcyBpY29uIGlzIGEgZHVwbGljYXRlIG9mIC9zcmMvY29yZS9jb21wb25lbnRzL2ltYWdlcy9idWlsdGlub3BlcmF0aW9uaWNvbnMvaHR0cC5zdmcgLS0+DQo8c3ZnIHdpZHRoPSIzMiIgaGVpZ2h0PSIzMiIgdmVyc2lvbj0iMS4xIiB2aWV3Qm94PSIwIDAgMzIgMzIiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+DQogPHBhdGggZmlsbD0iIzcwOTcyNyIgZD0ibTAgMGgzMnYzMmgtMzJ6Ii8+DQogPGcgZmlsbD0iI2ZmZiI+DQogIDxwYXRoIGQ9Ik0yMS4xMjcgMTAuOTgyYy0xLjA5MS0xLjgxOC0yLjk4Mi0yLjk4Mi01LjE2NC0yLjk4MnMtNC4wNzMgMS4xNjQtNS4wOTEgMi45MDljLS41MDkuODczLS44IDEuODkxLS44IDIuOTgyIDAgMy4wNTUgMi4zMjcgNS41MjcgNS4yMzYgNS44OTF2MS4wMThoMS4zODJ2LTEuMDE4YzIuOTgyLS4zNjQgNS4yMzYtMi44MzYgNS4yMzYtNS44OTEgMC0xLjAxOC0uMjkxLTIuMDM2LS44LTIuOTA5em0tMS4wMTguNTgyYy0uNDM2LjIxOC0xLjA5MS40MzYtMS44OTEuNTgyLS4xNDUtMS4xNjQtLjQzNi0yLjEwOS0uODczLTIuNzY0IDEuMTY0LjM2NCAyLjEwOSAxLjA5MSAyLjc2NCAyLjE4MnptLTIuMjU1IDIuNGMwIC42NTUtLjA3MyAxLjIzNi0uMTQ1IDEuNzQ1LS41MDkuMDczLTEuMDkxLjA3My0xLjc0NS4wNzNzLTEuMjM2IDAtMS43NDUtLjA3M2MtLjA3My0uNTgyLS4xNDUtMS4xNjQtLjE0NS0xLjc0NSAwLS40MzYgMC0uODczLjA3My0xLjMwOS41ODIuMDczIDEuMTY0LjE0NSAxLjgxOC4xNDVzMS4yMzYtLjA3MyAxLjgxOC0uMTQ1bC4wNzMgMS4zMDl6bS0xLjg5MS00LjhjLjIxOCAwIC40MzYgMCAuNjU1LjA3My40MzYuNTA5Ljg3MyAxLjYgMS4wOTEgMi45ODItLjUwOS4wNzMtMS4wOTEuMTQ1LTEuNzQ1LjE0NXMtMS4yMzYtLjA3My0xLjc0NS0uMTQ1Yy4yMTgtMS4zODIuNTgyLTIuNDczIDEuMDkxLTIuOTgyLjIxOC0uMDczLjQzNi0uMDczLjY1NS0uMDczem0tMS4zODIuMjE4Yy0uMzY0LjY1NS0uNzI3IDEuNi0uODczIDIuNzY0LS44LS4xNDUtMS40NTUtLjM2NC0xLjg5MS0uNTgyLjY1NS0xLjA5MSAxLjYtMS44MTggMi43NjQtMi4xODJ6bS0zLjQxOCA0LjU4MmMwLS43MjcuMTQ1LTEuMzgyLjQzNi0yLjAzNi41MDkuMjkxIDEuMjM2LjUwOSAyLjEwOS42NTUtLjA3My40MzYtLjA3My44NzMtLjA3MyAxLjM4MmwuMDczIDEuNzQ1Yy0xLjE2NC0uMTQ1LTEuOTY0LS40MzYtMi40NzMtLjcyN2wtLjA3My0xLjAxOHptLjI5MSAxLjZjLjU4Mi4yOTEgMS40NTUuNDM2IDIuMzI3LjU4Mi4xNDUuOTQ1LjQzNiAxLjgxOC44IDIuNC0xLjQ1NS0uNDM2LTIuNjE4LTEuNTI3LTMuMTI3LTIuOTgyem01LjE2NCAzLjEyN2wtLjY1NS4wNzNzLS40MzYgMC0uNjU1LS4wNzNjLS40MzYtLjUwOS0uOC0xLjMwOS0uOTQ1LTIuNDczLjU4Mi4wNzMgMS4wOTEuMDczIDEuNjczLjA3M3MxLjA5MSAwIDEuNjczLS4wNzNjLS4yOTEgMS4xNjQtLjY1NSAxLjk2NC0xLjA5MSAyLjQ3M3ptLjcyNy0uMTQ1Yy4zNjQtLjU4Mi42NTUtMS40NTUuOC0yLjQuODczLS4xNDUgMS43NDUtLjI5MSAyLjMyNy0uNTgyLS41MDkgMS40NTUtMS42NzMgMi41NDUtMy4xMjcgMi45ODJ6bTMuMjczLTMuNTY0Yy0uNTA5LjI5MS0xLjMwOS41ODItMi40NzMuNzI3LjA3My0uNTA5LjA3My0xLjA5MS4wNzMtMS43NDUgMC0uNDM2IDAtLjk0NS0uMDczLTEuMzgyLjgtLjE0NSAxLjUyNy0uMzY0IDIuMTA5LS42NTUuMjkxLjU4Mi40MzYgMS4zMDkuNDM2IDIuMDM2LjA3My4zNjQgMCAuNjU1LS4wNzMgMS4wMTh6TTEzLjg1NSAyMS4xNjRoNC4yMTh2MS44OTFoLTQuMjE4ek0xOC4zNjQgMjEuNjczaDEuNTI3djEuMzgyaC0xLjUyN3pNMTEuOTY0IDIxLjY3M2gxLjUyN3YxLjM4MmgtMS41Mjd6TTE1LjIzNiAyMy40MThoMS4zODJ2LjU4MmgtMS4zODJ6Ii8+DQogPC9nPg0KPC9zdmc+DQo=",
						"displayName": "HTTP"
					}
				},
				{
					"id": "HttpRequestReceived",
					"metadata": {
						"virtualConnector": true,
						"name": "When a HTTP request is received",
						"type": "Microsoft.PowerApps/apis",
						"iconUri": "data:image/svg+xml;base64,PCEtLSBQbGVhc2UgYWxzbyBjaGFuZ2UgdGhpcyBpbiBBenVyZS1CUE1VWC4gVGhpcyBpY29uIGlzIGEgZHVwbGljYXRlIG9mIC9zcmMvY29yZS9jb21wb25lbnRzL2ltYWdlcy9idWlsdGlub3BlcmF0aW9uaWNvbnMvcmVzcG9uc2Uuc3ZnIC0tPg0KPHN2ZyBlbmFibGUtYmFja2dyb3VuZD0ibmV3IDAgMCAzMiAzMiIgdmVyc2lvbj0iMS4xIiB2aWV3Qm94PSIwIDAgMzIgMzIiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+DQogPHBhdGggZD0ibTAgMGgzMnYzMmgtMzJ6IiBmaWxsPSIjMDA5ZGE1Ii8+DQogPGcgZmlsbD0iI2ZmZiI+DQogIDxwYXRoIGQ9Im0xMy45IDIxLjE2NGg0LjIxOHYxLjg5MWgtNC4yMTh6Ii8+DQogIDxwYXRoIGQ9Im0xOC40MDkgMjEuNjczaDEuNTI3djEuMzgyaC0xLjUyN3oiLz4NCiAgPHBhdGggZD0ibTEyLjAwOSAyMS42NzNoMS41Mjd2MS4zODJoLTEuNTI3eiIvPg0KICA8cGF0aCBkPSJNMTUuMjgyIDIzLjQxOGgxLjM4MnYuNTgyaC0xLjM4MnoiLz4NCiAgPHBhdGggZD0iTTI0Ljg4MiAxOC4xMDlsLTEuNjczLTEuNTI3aC0xLjAxOGwxLjIzNiAxLjE2NGgtMy43MDl2LjhoMy43MDlsLTEuMTY0IDEuMjM2aC45NDV6Ii8+DQogIDxwYXRoIGQ9Ik0xNS4yODIgMjAuODczaDEuMzgydi0xLjAxOGMuOC0uMDczIDEuNTI3LS4zNjQgMi4yNTUtLjcyN3YtMS4zODJjLS40MzYuMzY0LTEuMDE4LjY1NS0xLjYuOC4zNjQtLjU4Mi42NTUtMS40NTUuOC0yLjQuODczLS4xNDUgMS43NDUtLjI5MSAyLjMyNy0uNTgyLS4yMTguNTgyLS41MDkgMS4wMTgtLjg3MyAxLjQ1NWgxLjQ1NWMuNTgyLS44NzMuODczLTEuOTY0Ljg3My0zLjEyNyAwLTEuMDkxLS4yOTEtMi4xMDktLjgtMi45ODItMS4wMTgtMS43NDUtMi45MDktMi45MDktNS4wOTEtMi45MDktMi4xODIgMC00LjA3MyAxLjE2NC01LjA5MSAyLjkwOS0uNTA5Ljg3My0uOCAxLjg5MS0uOCAyLjk4MiAwIDMuMDU1IDIuMzI3IDUuNTI3IDUuMjM2IDUuODkxdjEuMDkxem0yLjQ3My01LjE2NGMtLjUwOS4wNzMtMS4wOTEuMDczLTEuNzQ1LjA3My0uNjU1IDAtMS4yMzYgMC0xLjc0NS0uMDczLS4wNzMtLjU4Mi0uMTQ1LTEuMTY0LS4xNDUtMS43NDUgMC0uNDM2IDAtLjg3My4wNzMtMS4zMDkuNTgyLjA3MyAxLjE2NC4xNDUgMS44MTguMTQ1LjY1NSAwIDEuMjM2LS4wNzMgMS44MTgtLjE0NS4wNzMuNDM2LjA3My44LjA3MyAxLjMwOSAwIC42NTUtLjA3MyAxLjIzNi0uMTQ1IDEuNzQ1em0yLjYxOC0zLjc4MmMuMjkxLjU4Mi40MzYgMS4zMDkuNDM2IDIuMDM2bC0uMDczIDEuMDE4Yy0uNTA5LjI5MS0xLjMwOS41ODItMi40NzMuNzI3LjA3My0uNTA5LjA3My0xLjA5MS4wNzMtMS43NDUgMC0uNDM2IDAtLjk0NS0uMDczLTEuMzgyLjgtLjE0NSAxLjUyNy0uMzY0IDIuMTA5LS42NTV6bS0uMjE4LS4zNjRjLS40MzYuMjE4LTEuMDkxLjQzNi0xLjg5MS41ODItLjE0NS0xLjE2NC0uNDM2LTIuMTA5LS44NzMtMi43NjQgMS4xNjQuMzY0IDIuMTA5IDEuMDkxIDIuNzY0IDIuMTgyem0tNC4xNDUtMi40Yy4yMTggMCAuNDM2IDAgLjY1NS4wNzMuNDM2LjUwOS44NzMgMS42IDEuMDkxIDIuOTgyLS41MDkuMDczLTEuMDkxLjE0NS0xLjc0NS4xNDUtLjY1NSAwLTEuMjM2LS4wNzMtMS43NDUtLjE0NS4yMTgtMS4zODIuNTgyLTIuNDczIDEuMDkxLTIuOTgyLjE0NS0uMDczLjQzNi0uMDczLjY1NS0uMDczem0tMS4zODIuMjE4Yy0uMzY0LjY1NS0uNzI3IDEuNi0uODczIDIuNzY0LS44LS4xNDUtMS40NTUtLjM2NC0xLjg5MS0uNTgyLjU4Mi0xLjA5MSAxLjYtMS44MTggMi43NjQtMi4xODJ6bS0zLjQxOCA0LjU4MmMwLS43MjcuMTQ1LTEuMzgyLjQzNi0yLjAzNi41MDkuMjkxIDEuMjM2LjUwOSAyLjEwOS42NTUtLjA3My40MzYtLjA3My44NzMtLjA3MyAxLjM4MmwuMDczIDEuNzQ1Yy0xLjE2NC0uMTQ1LTEuOTY0LS40MzYtMi40NzMtLjcyN2wtLjA3My0xLjAxOHptLjI5MSAxLjZjLjU4Mi4yOTEgMS40NTUuNDM2IDIuMzI3LjU4Mi4xNDUuOTQ1LjQzNiAxLjgxOC44IDIuNC0xLjQ1NS0uNDM2LTIuNjE4LTEuNTI3LTMuMTI3LTIuOTgyem0yLjgzNi42NTVjLjU4Mi4wNzMgMS4wOTEuMDczIDEuNjczLjA3My41ODIgMCAxLjA5MSAwIDEuNjczLS4wNzMtLjIxOCAxLjE2NC0uNTgyIDIuMDM2LS45NDUgMi40NzNsLS42NTUuMDczYy0uMjE4IDAtLjQzNiAwLS42NTUtLjA3My0uNTA5LS41MDktLjg3My0xLjMwOS0xLjA5MS0yLjQ3M3oiLz4NCiA8L2c+DQo8L3N2Zz4NCg==",
						"displayName": "When a HTTP request is received"
					}
				}]`), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `[
				{
					"id": "/providers/Microsoft.PowerApps/apis/shared_approvals",
					"metadata": {
						"unblockable": true
					}
				},
				{
					"id": "/providers/Microsoft.PowerApps/apis/shared_cloudappsecurity",
					"metadata": {
						"unblockable": true
					}
				},
				{
					"id": "/providers/Microsoft.PowerApps/apis/shared_commondataservice",
					"metadata": {
						"unblockable": true
					}
				}],`), nil
		})

	httpmock.RegisterResponder("GET", `https://api.powerapps.com/providers/Microsoft.PowerApps/apis?%24filter=environment+eq+%27~Default%27&api-version=2023-06-01&hideDlpExemptApis=true&showAllDlpEnforceableApis=true&showApisWithToS=true`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"name": "shared_sharepointonline",
						"id": "/providers/Microsoft.PowerApps/apis/shared_sharepointonline",
						"type": "Microsoft.PowerApps/apis",
						"properties": {
							"displayName": "SharePoint",
							"iconUri": "https://connectoricons-prod.azureedge.net/u/henryorsborn/partial-builds/asev3migrations-with-resourceTemplate/1.0.1653.3414/sharepointonline/icon.png",
							"iconBrandColor": "#036C70",
							"apiEnvironment": "Shared",
							"isCustomApi": false,
							"connectionParameters": {
								"token": {
									"type": "oauthSetting",
									"oAuthSettings": {
										"identityProvider": "sharepointonlinecertificateV2",
										"clientId": "7ab7862c-4c57-491e-8a45-d52a7e023983",
										"scopes": [],
										"redirectMode": "GlobalPerConnector",
										"redirectUrl": "https://global.consent.azure-apim.net/redirect/sharepointonline",
										"properties": {
											"IsFirstParty": "True",
											"IsOnbehalfofLoginSupported": true
										},
										"customParameters": {
											"resourceUriAAD": {
												"value": "https://graph.microsoft.com/"
											},
											"loginUri": {
												"value": "https://login.windows.net"
											},
											"loginUriAAD": {
												"value": "https://login.windows.net"
											},
											"resourceUri": {
												"value": "https://graph.microsoft.com"
											}
										}
									},
									"uiDefinition": {
										"displayName": "Log in with SharePoint Credentials",
										"description": "Log in with SharePoint Credentials",
										"tooltip": "Provide SharePoint Credentials",
										"constraints": {
											"required": "true",
											"capability": [
												"cloud"
											]
										}
									}
								},
								"token:TenantId": {
									"type": "string",
									"metadata": {
										"sourceType": "AzureActiveDirectoryTenant"
									},
									"uiDefinition": {
										"displayName": "Tenant",
										"description": "The tenant ID of for the Azure Active Directory application",
										"constraints": {
											"required": "false",
											"hidden": "true"
										}
									}
								},
								"gateway": {
									"type": "gatewaySetting",
									"gatewaySettings": {
										"dataSourceType": "SharePoint",
										"connectionDetails": []
									},
									"uiDefinition": {
										"tabIndex": 1,
										"constraints": {
											"hidden": "false",
											"capability": [
												"gateway"
											]
										}
									}
								},
								"authType": {
									"type": "string",
									"allowedValues": [
										{
											"value": "windows"
										}
									],
									"uiDefinition": {
										"displayName": "Authentication Type",
										"description": "Authentication type to connect to your database",
										"tooltip": "Authentication type to connect to your database",
										"constraints": {
											"tabIndex": 2,
											"required": "false",
											"allowedValues": [
												{
													"text": "Windows",
													"value": "windows"
												}
											],
											"capability": [
												"gateway"
											]
										}
									}
								},
								"username": {
									"type": "securestring",
									"uiDefinition": {
										"displayName": "Username",
										"description": "Username credential",
										"tooltip": "Username credential",
										"constraints": {
											"tabIndex": 3,
											"clearText": true,
											"required": "true",
											"capability": [
												"gateway"
											]
										}
									}
								},
								"password": {
									"type": "securestring",
									"uiDefinition": {
										"displayName": "Password",
										"description": "Password credential",
										"tooltip": "Password credential",
										"constraints": {
											"tabIndex": 4,
											"required": "true",
											"capability": [
												"gateway"
											]
										}
									}
								}
							},
							"runtimeUrls": [
								"https://europe-001.azure-apim.net/apim/sharepointonline"
							],
							"primaryRuntimeUrl": "https://europe-001.azure-apim.net/apim/sharepointonline",
							"metadata": {
								"source": "marketplace",
								"brandColor": "#036C70",
								"useNewApimVersion": "true",
								"version": {
									"previous": "releases/v1.0.1656\\1.0.1656.3432",
									"current": "u/henryorsborn/partial-builds/asev3migrations-with-resourceTemplate\\1.0.1653.3414"
								}
							},
							"capabilities": [
								"tabular",
								"gateway",
								"cloud"
							],
							"interfaces": {
								"CDPTabular1": {
									"revisions": {
										"1": {
											"baseUrl": "/",
											"status": "Production"
										}
									}
								}
							},
							"description": "SharePoint helps organizations share and collaborate with colleagues, partners, and customers. You can connect to SharePoint Online or to an on-premises SharePoint 2013 or 2016 farm using the On-Premises Data Gateway to manage documents and list items.",
							"createdTime": "2016-10-07T18:40:04.372652Z",
							"changedTime": "2023-09-26T00:11:38.2143635Z",
							"releaseTag": "Production",
							"tier": "Standard",
							"publisher": "Microsoft",
							"scopes": {
								"will": [
									"Read list and library names, as well as the names of the columns",
									"Create, read, update, copy and delete files and metadata",
									"Create, read, update, and delete list items"
								],
								"wont": []
							}
						}
					},
					{
						"name": "shared_onedriveforbusiness",
						"id": "/providers/Microsoft.PowerApps/apis/shared_onedriveforbusiness",
						"type": "Microsoft.PowerApps/apis",
						"properties": {
							"displayName": "OneDrive for Business",
							"iconUri": "https://connectoricons-prod.azureedge.net/releases/v1.0.1656/1.0.1656.3432/onedriveforbusiness/icon.png",
							"iconBrandColor": "#0078D4",
							"apiEnvironment": "Shared",
							"isCustomApi": false,
							"connectionParameters": {
								"token": {
									"type": "oauthSetting",
									"oAuthSettings": {
										"identityProvider": "OneDriveForBusinessCertificate",
										"clientId": "7ab7862c-4c57-491e-8a45-d52a7e023983",
										"scopes": [],
										"redirectMode": "GlobalPerConnector",
										"redirectUrl": "https://global.consent.azure-apim.net/redirect/onedriveforbusiness",
										"properties": {
											"IsFirstParty": "True",
											"IsOnbehalfofLoginSupported": true
										},
										"customParameters": {
											"capability": {
												"value": "MyFiles"
											},
											"grantType": {
												"value": "code"
											},
											"resourceUri": {
												"value": "https://graph.microsoft.com"
											},
											"resourceUriAAD": {
												"value": "https://graph.microsoft.com"
											},
											"loginUriAAD": {
												"value": "https://login.windows.net"
											}
										}
									},
									"uiDefinition": {
										"displayName": "Log in with OneDrive for Business Credentials",
										"description": "Log in with OneDrive for Business Credentials",
										"tooltip": "Provide OneDrive for Business Credentials",
										"constraints": {
											"required": "true"
										}
									}
								}
							},
							"runtimeUrls": [
								"https://europe-001.azure-apim.net/apim/onedriveforbusiness"
							],
							"primaryRuntimeUrl": "https://europe-001.azure-apim.net/apim/onedriveforbusiness",
							"metadata": {
								"source": "marketplace",
								"brandColor": "#0078D4",
								"useNewApimVersion": "true",
								"version": {
									"previous": "releases/v1.0.1647\\1.0.1647.3361",
									"current": "releases/v1.0.1656\\1.0.1656.3432"
								}
							},
							"capabilities": [
								"blob"
							],
							"interfaces": {
								"CDPBlob0": {
									"revisions": {
										"1": {
											"baseUrl": "/",
											"status": "Production",
											"deprecated": true
										}
									}
								},
								"CDPBlob1": {
									"revisions": {
										"1": {
											"baseUrl": "/",
											"status": "Production"
										}
									}
								}
							},
							"description": "OneDrive for Business is a cloud storage, file hosting service that allows users to sync files and later access them from a web browser or mobile device. Connect to OneDrive for Business to manage your files. You can perform various actions such as upload, update, get, and delete files.",
							"createdTime": "2016-09-30T04:12:48.5709476Z",
							"changedTime": "2023-09-21T23:21:20.0282274Z",
							"releaseTag": "Production",
							"tier": "Standard",
							"publisher": "Microsoft",
							"scopes": {
								"will": [
									"Read your user profile",
									"Create, read, update, and delete files"
								],
								"wont": []
							}
						}
					}
				]}`), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: UniTestsProviderConfig + `
				data "powerplatform_connectors" "all" {}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify placeholder id attribute
					resource.TestMatchResourceAttr("data.powerplatform_connectors.all", "id", regexp.MustCompile(`^[1-9]\d*$`)),

					//Verify returned count
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.#", "4"),

					// // Verify the first power app to ensure all attributes are set
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.description", "SharePoint helps organizations share and collaborate with colleagues, partners, and customers. You can connect to SharePoint Online or to an on-premises SharePoint 2013 or 2016 farm using the On-Premises Data Gateway to manage documents and list items."),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.display_name", "SharePoint"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.name", "shared_sharepointonline"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.publisher", "Microsoft"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.tier", "Standard"),
					resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.0.type", "Microsoft.PowerApps/apis"),
				),
			},
		},
	})
}
