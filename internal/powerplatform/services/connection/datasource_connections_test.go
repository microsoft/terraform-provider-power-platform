// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/provider"
)

func TestAccConnectionsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "env" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 						  = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "null_resource" "wait_60_seconds" {
					provisioner "local-exec" {
						command = "sleep 60"
					}
					depends_on = [powerplatform_environment.env]
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

					depends_on = [ null_resource.wait_60_seconds ]
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
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.display_name", "OpenAI Connection "+mocks.TestName()),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.connection_parameters", "{\"azureOpenAIResourceName\":\"aaa\",\"azureSearchEndpointUrl\":\"ccc\",\"sku\":\"Enterprise\"}"),
					resource.TestCheckNoResourceAttr("data.powerplatform_connections.all_connections", "connections.0.connection_parameters_set"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.status.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.status.0", "Connected"),
				),
			},
		},
	})
}

func TestUnitConnectionsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connections?api-version=1`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/connections/Validate_Read/get_connections.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: provider.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.TestsUnitProviderConfig + `

				data "powerplatform_connections" "all_connections" {
					environment_id = "00000000-0000-0000-0000-000000000000"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.name", "shared_commondataserviceforapps"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.display_name", "My CDS Connection"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.1.name", "shared_flowmanagement"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.1.display_name", "My Flow Connection"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.1.id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.status.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.status.0", "Connected"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.1.status.#", "1"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.1.status.0", "Connected"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.connection_parameters", "{\"sku\":\"Enterprise\",\"token:grantType\":\"code\"}"),
					resource.TestCheckNoResourceAttr("data.powerplatform_connections.all_connections", "connections.1.connection_parameters"),
					resource.TestCheckNoResourceAttr("data.powerplatform_connections.all_connections", "connections.0.connection_parameters_set"),
					resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.1.connection_parameters_set", "{\"name\":\"firstParty\",\"values\":{}}"),
				),
			},
		},
	})
}
