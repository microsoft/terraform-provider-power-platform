// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccConnectionsResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsAcceptanceProviderConfig + `
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
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection "+mock_helpers.TestName()),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters", "{\"azureOpenAIApiKey\":\"bbb\",\"azureOpenAIResourceName\":\"aaa\",\"azureSearchApiKey\":\"ddd\",\"azureSearchEndpointUrl\":\"ccc\"}"),
					resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connections.0.connection_parameters_set"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
				),
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(`https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("services/connection/tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/connection/tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id = "00000000-0000-0000-0000-000000000000"
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
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters", "{\"azureOpenAIApiKey\":\"bbb\",\"azureOpenAIResourceName\":\"aaa\",\"azureSearchApiKey\":\"ddd\",\"azureSearchEndpointUrl\":\"ccc\"}"),
					resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connections.0.connection_parameters_set"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
				),
			},
		},
	})
}
