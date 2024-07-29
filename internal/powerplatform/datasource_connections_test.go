// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccConnectionsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "` + mock_helpers.TestName() + `"
					location                                  = "europe"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
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

				data "powerplatform_connections" "all_connections" {
					environment_id = powerplatform_environment.env.id

					depends_on = [
						powerplatform_connection.azure_openai_connection
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.display_name", "OpenAI Connection"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.connection_parameters", "{\"azureOpenAIResourceName\":\"aaa\",\"azureSearchEndpointUrl\":\"ccc\",\"sku\":\"Enterprise\"}"),
					resource.TestCheckNoResourceAttr("data.powerplatform_connections.all_connections", "connections.0.connection_parameters_set"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.status.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.status.0", "Connected"),
				),
			},
		},
	})
}
