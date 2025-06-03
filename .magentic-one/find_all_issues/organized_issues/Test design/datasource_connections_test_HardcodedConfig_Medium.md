# Issue 1

Hardcoded and Non-Parametrized Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go

## Problem

The tests embed large Terraform resource and data definitions as inline raw strings within the Go code. This practice leads to maintenance difficulties, as changes to configuration formats or variable values must be made in several scattered locations. It also makes the test code harder to read, and introduces risks of accidentally introducing syntax errors in the configuration.

## Impact

- **Severity:** Medium
- Test maintenance becomes more difficult as configuration complexity or repetition grows.
- Copy-paste errors or inconsistent test configuration become more likely.
- Hard to review or update resource definitions for different test cases.

## Location

```go
Config: `
				resource "powerplatform_environment" "env" {
					display_name                              = "` + mocks.TestName() + `"
					location                                  = "unitedstates"
					environment_type                          = "Sandbox"
					dataverse = {
						language_code                             = "1033"
						currency_code                             = "USD"
						security_group_id 					  = "00000000-0000-0000-0000-000000000000"
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

				data "powerplatform_connections" "all_connections" {
					environment_id = powerplatform_environment.env.id

					depends_on = [
						powerplatform_connection.azure_openai_connection
					]
				}
				`,
```

## Fix

Define reusable test configuration templates (constants or functions) and substitute values dynamically. This also enables sharing across multiple test functions, and reduces risk of error.

```go
const envResource = `
resource "powerplatform_environment" "env" {
	display_name    = "%s"
	location        = "unitedstates"
	environment_type = "Sandbox"
	dataverse = {
		language_code     = "1033"
		currency_code     = "USD"
		security_group_id = "00000000-0000-0000-0000-000000000000"
	}
}
`

const azureOpenAIConnection = `
resource "powerplatform_connection" "azure_openai_connection" {
	environment_id = powerplatform_environment.env.id
	name           = "shared_azureopenai"
	display_name   = "OpenAI Connection %s"
	connection_parameters = jsonencode({
		"azureOpenAIResourceName": "aaa",
		"azureOpenAIApiKey": "bbb",
		"azureSearchEndpointUrl": "ccc",
		"azureSearchApiKey": "ddd"
	})

	lifecycle {
		ignore_changes = [
			connection_parameters
		]
	}
}
`

const dataSource = `
data "powerplatform_connections" "all_connections" {
	environment_id = powerplatform_environment.env.id
	depends_on = [
		powerplatform_connection.azure_openai_connection
	]
}
`

// In your test...
config := fmt.Sprintf(envResource+mocks.TestName(), azureOpenAIConnection+mocks.TestName(), dataSource)
...
Config: config,
```

Or even parameterize only values that need to change. This practice also encourages splitting out test scenario data for easier expansion.
