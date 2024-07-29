// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectionsShareDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				//lintignore:AT004
				Config: TestsProviderConfig + `
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
					user_principal_name = "test@${local.domain_name}"
					display_name        = "test"
					mail_nickname       = "test"
					password            = random_password.passwords.result
					usage_location      = "US"
				  }
				  
				  resource "powerplatform_environment" "env" {
					display_name     = "wercker"
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
					display_name   = "OpenAI Connection"
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
