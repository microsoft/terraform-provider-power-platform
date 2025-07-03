// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccConnectionsShareResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"azuread": {
				VersionConstraint: constants.AZURE_AD_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/azuread",
			},
			"random": {
				VersionConstraint: constants.RANDOM_PROVIDER_VERSION_CONSTRAINT,
				Source:            "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			{
				ResourceName: "powerplatform_connection_share.share_with_user1",
				Config: `
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

				resource "powerplatform_connection" "azure_ai_search_connection" {
					environment_id = powerplatform_environment.env.id
					name           = "shared_azureaisearch"
					display_name   = "Azure AI Search Connection ` + mocks.TestName() + `"
					connection_parameters = jsonencode({
						ConnectionEndpoint = "aaa"
						AdminKey           = "bbb"
					})

					lifecycle {
						ignore_changes = [
						connection_parameters
						]
					}
				}

				resource "powerplatform_connection_share" "share_with_user1" {
					environment_id = powerplatform_environment.env.id
					connector_name = powerplatform_connection.azure_ai_search_connection.name
					connection_id  = powerplatform_connection.azure_ai_search_connection.id
					role_name      = "CanEdit"
					principal = {
						entra_object_id = azuread_user.test_user.object_id
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "connector_name", "shared_azureaisearch"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "role_name", "CanEdit"),
					resource.TestCheckResourceAttr("powerplatform_connection_share.share_with_user1", "principal.display_name", mocks.TestName()),
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connection_shares/Validate_Create/get_connection_shares.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
