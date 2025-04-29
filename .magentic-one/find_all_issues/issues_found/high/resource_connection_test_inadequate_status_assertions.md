# Title

Inadequate Test Assertions on Status Values

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

The test cases verify the `status` attribute but do not check for multiple scenarios (e.g., status changes in error conditions). Only a generic "Connected" value is validated.

#### Violating Code:

```go
resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
```

## Impact

- **Severity**: High
- The test does not fully validate the behavior of the `status` attribute during various lifecycle states (e.g., when failed or disconnected). This reduces confidence in robustness and could lead to undetected issues in production.

## Location

Occurs in both test cases:

1. `TestAccConnectionsResource_Validate_Create`: at line 48.
2. `TestUnitConnectionsResource_Validate_Create`: at line 110.

## Fix

Add test assertions to validate other potential values for `status`:

```go
resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Connected"),
resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Disconnected"),
resource.TestCheckResourceAttr("powerplatform_connection.azure_openai_connection", "status.0", "Failed"),
```