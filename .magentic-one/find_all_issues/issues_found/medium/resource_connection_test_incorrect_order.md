# Title

Incorrect `Check` List Order in Test Definitions

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

In both test cases, the aggregate checking functions specified in the `Check` field don't follow consistent logical ordering (e.g., verifying `name` before `display_name`). Additionally, checks for nested attributes such as `connection_parameters` are inconsistently grouped.

#### Violating Code:

```go
Check: resource.ComposeAggregateTestCheckFunc(
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection "+mocks.TestName()),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters", "{\"azureOpenAIApiKey\":\"bbb\",\"azureOpenAIResourceName\":\"aaa\",\"azureSearchApiKey\":\"ddd\",\"azureSearchEndpointUrl\":\"ccc\"}"),
	resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connections.0.connection_parameters_set"),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
),
```

## Impact

- **Severity**: Medium
- Reduces readability, understanding, and maintainability of the tests, especially when debugging test failures.

## Location

Two occurrences found:

1. Part of `TestAccConnectionsResource_Validate_Create` function at line 42.
2. Part of `TestUnitConnectionsResource_Validate_Create` function at line 104.

## Fix

Improve logical grouping and order of the checks for better clarity:

```go
Check: resource.ComposeAggregateTestCheckFunc(
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "name", "shared_azureopenai"),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "display_name", "OpenAI Connection "+mocks.TestName()),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "connection_parameters", "{\"azureOpenAIApiKey\":\"bbb\",\"azureOpenAIResourceName\":\"aaa\",\"azureSearchApiKey\":\"ddd\",\"azureSearchEndpointUrl\":\"ccc\"}"),
	resource.TestCheckNoResourceAttr("powerplatform_connection.azure_openai_connection", "connections.0.connection_parameters_set"),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.#", "1"),
	resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
),
```