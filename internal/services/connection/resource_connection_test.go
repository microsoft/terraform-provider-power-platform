// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccConnectionsResource_Validate_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
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

					resource "time_sleep" "wait_for_dataverse" {
						create_duration = "120s"

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

						depends_on = [time_sleep.wait_for_dataverse]
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection "+mocks.TestName()),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters", "{\"azureOpenAIApiKey\":\"bbb\",\"azureOpenAIResourceName\":\"aaa\",\"azureSearchApiKey\":\"ddd\",\"azureSearchEndpointUrl\":\"ccc\"}"),
					resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connections.0.connection_parameters_set"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
				),
			},
		},
	})
}

func TestAccConnectionsResource_Validate_Create_Without_Config_Parameters(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
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

					resource "time_sleep" "wait_for_dataverse" {
						create_duration = "120s"

						depends_on = [powerplatform_environment.env]
					}

					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id = powerplatform_environment.env.id
						name           = "shared_azureopenai"
						display_name   = "OpenAI Connection ` + mocks.TestName() + `"

						depends_on = [time_sleep.wait_for_dataverse]
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection "+mocks.TestName()),
					// Without real credentials the connector genuinely can't connect, so status is "Error".
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Error"),
				),
			},
			{
				// Real API omits connectionParameters/connectionParametersSet when they were never
				// configured; a refresh must not leave the attributes unknown or produce drift.
				RefreshState: true,
			},
		},
	})
}

func TestAccConnectionsResource_Validate_Create_Invalid_Connection_Parameters(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id         = "00000000-0000-0000-0000-000000000000"
						name                   = "shared_azureopenai"
						display_name           = "OpenAI Connection ` + mocks.TestName() + `"
						connection_parameters  = "{"
					}
					`,
				ExpectError: regexp.MustCompile("Failed to convert connection parameters"),
			},
		},
	})
}

func TestAccConnectionsResource_Validate_Create_Invalid_Connection_Parameters_Set(t *testing.T) {
	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id             = "00000000-0000-0000-0000-000000000000"
						name                       = "shared_azureopenai"
						display_name               = "OpenAI Connection ` + mocks.TestName() + `"
						connection_parameters_set  = "{"
					}
					`,
				ExpectError: regexp.MustCompile("Failed to convert connection parameters set"),
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Create_Retries_When_Environment_Not_Yet_Visible(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	putAttempts := 0
	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			putAttempts++
			if putAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/connections/Validate_Create_Retries_When_Environment_Not_Yet_Visible/put_connection_environment_not_found.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/connections/Validate_Create_Retries_When_Environment_Not_Yet_Visible/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connections/Validate_Create_Retries_When_Environment_Not_Yet_Visible/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id = "00000000-0000-0000-0000-000000000000"
						name           = "shared_azureopenai"
						display_name   = "OpenAI Connection"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
					func(*terraform.State) error {
						if putAttempts < 2 {
							return fmt.Errorf("expected the connection create to be retried, got %d attempt(s)", putAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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

func TestUnitConnectionsResource_Validate_Create_Without_Config_Parameters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id = "00000000-0000-0000-0000-000000000000"
						name           = "shared_azureopenai"
						display_name   = "OpenAI Connection"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters", "{\"azureOpenAIResourceName\":\"aaa\",\"azureSearchEndpointUrl\":\"ccc\",\"sku\":\"Enterprise\"}"),
					resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters_set"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
				),
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Create_Invalid_Connection_Parameters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id         = "00000000-0000-0000-0000-000000000000"
						name                   = "shared_azureopenai"
						display_name           = "OpenAI Connection"
						connection_parameters  = "{"
					}
					`,
				ExpectError: regexp.MustCompile("Failed to convert connection parameters"),
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Create_Invalid_Connection_Parameters_Set(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id             = "00000000-0000-0000-0000-000000000000"
						name                       = "shared_azureopenai"
						display_name               = "OpenAI Connection"
						connection_parameters_set  = "{"
					}
					`,
				ExpectError: regexp.MustCompile("Failed to convert connection parameters set"),
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Create_Missing_Parameters_In_Response(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/connections/Validate_Create_Missing_Parameters/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connections/Validate_Create_Missing_Parameters/put_connection.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest: true,

		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// The mocked API response omits connectionParameters and connectionParametersSet
				// entirely; both attributes must normalize to null (never unknown) so the apply
				// and the follow-up refresh/plan succeed without producing a diff.
				Config: `
					resource "powerplatform_connection" "azure_openai_connection" {
						environment_id = "00000000-0000-0000-0000-000000000000"
						name           = "shared_azureopenai"
						display_name   = "OpenAI Connection"
					}
					`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection"),
					resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters"),
					resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters_set"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
				),
			},
			{
				// Explicit refresh exercises the Read path against the same
				// parameters-omitted response; state must remain stable.
				RefreshState: true,
			},
		},
	})
}

func TestUnitConnectionsResource_Validate_Read_ParentDeleted(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
		})

	var getConnectionCallCount = 0
	// First two calls (Create + post-Create plan check) succeed; step 2 refresh simulates a 403 because the environment is gone.
	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			getConnectionCallCount++
			if getConnectionCallCount <= 2 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
			}
			// Simulate 403 ConnectionAuthorizationFailed when the parent environment has been deleted.
			return httpmock.NewStringResponse(http.StatusForbidden, `{"error":{"code":"ConnectionAuthorizationFailed","message":"The caller does not have permission."}}`), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(`^https://000000000000000000000000000000\.00\.environment\.api\.powerplatform\.com/connectivity/connectors/shared_azureopenai/connections/(.*)?%24filter=environment\+eq\+%2700000000-0000-0000-0000-000000000000%27&api-version=1$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	// When the connection GET fails with a non-ErrObjectNotFound error, the provider checks
	// whether the parent environment still exists. Return 404 to indicate it is gone.
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000000\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, httpmock.File("tests/resource/connections/Validate_Read_ParentDeleted/get_environment_00000000-0000-0000-0000-000000000000.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: create the connection successfully.
				Config: `
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
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
					resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection"),
				),
			},
			{
				// Step 2: refresh when the parent environment has been deleted out-of-band.
				// The connection GET returns 403; the provider detects the parent is gone and removes the resource from state.
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
