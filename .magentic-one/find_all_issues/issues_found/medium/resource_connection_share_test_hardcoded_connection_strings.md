# Title

Hardcoded Connection Strings in Unit Tests

##

`/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go`

## Problem

In the unit test `TestUnitConnectionsShareResource_Validate_Create`, there are hardcoded URLs for API endpoints:

```go
https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/modifyPermissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1
https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/permissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1
```

Hardcoding such values in test code can lead to brittle tests, reduced configuration flexibility, and difficulties in maintaining the test framework if environments or API formats change.

## Impact

Hardcoded URLs are problematic because:
1. They are inflexible and cannot be dynamically adjusted based on environment or configuration changes.
2. They increase the risk of errors in tests if the hardcoded values become outdated.
3. They expose implementation details in test code that should ideally be abstracted, harming test portability and readability.

Severity: **Medium**

## Location

Unit Test: `TestUnitConnectionsShareResource_Validate_Create`

## Code Issue

The following snippet contains hardcoded API URLs:

```go
httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/modifyPermissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})

httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/permissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connection_shares/Validate_Create/get_connection_shares.json").String()), nil
	})
```

## Fix

Use configuration variables or constants instead of hardcoding API URLs. Introduce a configuration abstraction to make the test flexible to changes.

```go
const apiBaseUrl = "https://000000000000000000000000000000.00.environment.api.powerplatform.com"

httpmock.RegisterResponder("POST", fmt.Sprintf("%s/connectivity/connectors/shared_commondataserviceforapps/connections/%s/modifyPermissions?%24filter=environment+eq+%s&api-version=1", apiBaseUrl, connectionID, environmentID),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})

httpmock.RegisterResponder("GET", fmt.Sprintf("%s/connectivity/connectors/shared_commondataserviceforapps/connections/%s/permissions?%24filter=environment+eq+%s&api-version=1", apiBaseUrl, connectionID, environmentID),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connection_shares/Validate_Create/get_connection_shares.json").String()), nil
	})
```

Explanation:
- Replace hardcoded strings with dynamic values using a base URL constant and format strings.
- Ensure test cases adapt seamlessly to environment-specific configurations.