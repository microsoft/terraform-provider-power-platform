# Title

Missing Assertion Messages in Resource Test Checks

##

`/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go`

## Problem

In the `Check` section of `TestAccConnectionsDataSource_Validate_Read`, assertions are made to validate resource attributes (e.g., `resource.TestCheckResourceAttr`). However, no custom messages are provided for these assertions to describe the failure if the test does not pass.

This limits the debugging and interpretability of test failures. Developers will have to manually inspect the code to infer the intention behind each test check.

## Impact

The lack of assertion messages makes debugging test failures harder and reduces the clarity of automated testing. Severity: **Medium**, because tests still run but are less developer-friendly when debugging.

## Location

- Test function: `TestAccConnectionsDataSource_Validate_Read`
- Lines where resource.TestCheckResourceAttr and similar methods are invoked.

## Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.#", "1"),
resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.name", "shared_azureopenai"),
resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.display_name", "OpenAI Connection "+mocks.TestName()),
```

## Fix

Add descriptive messages for assertions using `resource.TestCheckResourceAttrWithMessage` or equivalent custom logging:

```go
resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.#", "1", "Expected number of connections does not match."),
resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.name", "shared_azureopenai", "Expected connection name does not match."),
resource.TestCheckResourceAttr("data.powerplatform_connections.all_connections", "connections.0.display_name", "OpenAI Connection "+mocks.TestName(), "Expected connection display name does not match."),
```

This improves the ability to debug test failures by providing precise context for why assertions failed.