// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccConnectionsShareDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				//lintignore:AT004
				Config: constants.TestsAcceptanceProviderConfig + `
				provider "azuread" {
				}

				data "azuread_domains" "aad_domains" {
					only_initial = true
				}

				data "azuread_group" "licensing_group" {
					display_name     = "` + mocks.TestsEntraLicesingGroupName() + `"
					security_enabled = true
				}

				resource "azuread_group_member" "example" {
					group_object_id  = data.azuread_group.licensing_group.object_id
					member_object_id = azuread_user.test_user.object_id
				}

				locals {
					domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
				}

				resource "random_password" "passwords" {
				 	min_lower = 1
					min_upper        = 1
					min_numeric      = 1
					min_special      = 1
					length           = 16
					special          = true
					override_special = "_%@"
				}

				resource "azuread_user" "test_user" {
					user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"
					display_name        = "` + mocks.TestName() + `"
					mail_nickname       = "` + mocks.TestName() + `"
					password            = random_password.passwords.result
					usage_location      = "US"
				}

				resource "powerplatform_environment" "env" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
					  language_code     = "1033"
					  currency_code     = "USD"
					  security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "powerplatform_connection" "azure_openai_connection" {
					environment_id = powerplatform_environment.env.id
					name           = "shared_azureopenai"
					display_name   = "OpenAI Connection ` + mocks.TestName() + `"
					connection_parameters = jsonencode({
					  "azureOpenAIResourceName" : "aaa",
					  "azureOpenAIApiKey" : "bbb"
					  "azureSearchEndpointUrl" : "ccc",
					  "azureSearchApiKey" : "ddd"
					})

					lifecycle {
					  ignore_changes = [
						connection_parameters
					  ]
					}
				}

				resource "powerplatform_connection_share" "share_with_user1" {
					environment_id = powerplatform_environment.env.id
					connector_name = powerplatform_connection.azure_openai_connection.name
					connection_id  = powerplatform_connection.azure_openai_connection.id
					role_name      = "CanEdit"
					principal = {
					  entra_object_id = azuread_user.test_user.object_id
					}
				}

				data "powerplatform_connection_shares" "all_shares" {
					environment_id = powerplatform_environment.env.id
					connector_name = "shared_azureopenai"
					connection_id  = powerplatform_connection.azure_openai_connection.id

					depends_on = [
					  powerplatform_connection_share.share_with_user1
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.#", "2"),
				),
			},
		},
	})
}

func TestUnitConnectionsShareDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.01.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000002/permissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000001%27&api-version=1`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/connection_shares/Validate_Read/get_connection_shares.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_connection_shares" "all_shares" {
					environment_id = "00000000-0000-0000-0000-000000000001"
					connector_name = "shared_commondataserviceforapps"
					connection_id  = "00000000-0000-0000-0000-000000000002"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.0.id", "4fac9a58-a88e-417e-bad7-fe0ea64582e3"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.0.role_name", "CanViewWithShare"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.0.principal.entra_object_id", "c52d4a80-3db9-4152-8601-2f8c131c4a90"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.0.principal.display_name", "Power Platform API"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.1.id", "f99f844b-ce3b-49ae-86f3-e374ecae789c"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.1.role_name", "Owner"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.1.principal.entra_object_id", "f99f844b-ce3b-49ae-86f3-e374ecae789c"),
					resource.TestCheckResourceAttr("data.powerplatform_connection_shares.all_shares", "shares.1.principal.display_name", "admin"),
				),
			},
		},
	})
}
