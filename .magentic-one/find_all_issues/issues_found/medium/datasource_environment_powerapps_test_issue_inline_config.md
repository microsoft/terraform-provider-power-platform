### Issue Report #1

## Title

Improper Use of Inline Configuration in Acceptance Test

##

File Path:

`/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps_test.go`

## Problem

In the test case `TestAccEnvironmentPowerAppsDataSource_Basic`, the configuration string is written inline within the test function. This violates best practices, as inline configurations make tests hard to read, modify, and debug. Using separate files for test configurations is recommended for clarity and maintainability.

## Impact

Inline configuration reduces the code's readability and makes debugging or fixing potential issues with the configuration more difficult. Furthermore, it is less modular, making reuse across test cases challenging.

Severity: **Medium**

## Location

File location where the issue was found:

```go
						Config: `
				environment" "env" {
						display_name "` + mocks.TestNamelocation"]
		}