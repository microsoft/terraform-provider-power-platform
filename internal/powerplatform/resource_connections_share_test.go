// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccConnectionsShareResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				//lintignore:AT004
				Config: TestsAcceptanceProviderConfig + `
				terraform {
					required_providers {
					  azuread = {
						source = "hashicorp/azuread"
					  }
					  random = {
						source = "hashicorp/random"
					  }
					}
				}

				provider "azuread" {
					use_cli = true
				}
				  
				data "azuread_domains" "aad_domains" {
					only_initial = true
				}
				
				locals {
					domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
				}
				
				resource "random_password" "passwords" {
					length           = 16
					special          = true
					override_special = "_%@"
				}
				
				resource "azuread_user" "test_user" {
					user_principal_name = "` + mock_helpers.TestName() + `@${local.domain_name}"
					display_name        = "` + mock_helpers.TestName() + `"
					mail_nickname       = "` + mock_helpers.TestName() + `"
					password            = random_password.passwords.result
					usage_location      = "US"
				}
				
				resource "powerplatform_environment" "env" {
					display_name     = "` + mock_helpers.TestName() + `"
					location         = "europe"
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
					display_name   = "OpenAI Connection ` + mock_helpers.TestName() + `"
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
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "connector_name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "role_name", "CanEdit"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "principal.display_name", mock_helpers.TestName()),
				),
			},
		},
	})
}

func TestUnitConnectionsShareResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/modifyPermissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/permissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/connection/tests/resource/connection_shares/Validate_Create/get_connection_shares.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
				resource "powerplatform_connection_share" "share_with_user1" {
					environment_id = "00000000-0000-0000-0000-000000000000"
					connector_name = "shared_commondataserviceforapps"
					connection_id  = "00000000-0000-0000-0000-000000000001"
					role_name      = "CanViewWithShare"
					principal = {
						entra_object_id = "00000000-0000-0000-0000-000000000002"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "connector_name", "shared_commondataserviceforapps"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "role_name", "CanViewWithShare"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "principal.display_name", "Power Platform API"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "principal.entra_object_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "connection_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}
