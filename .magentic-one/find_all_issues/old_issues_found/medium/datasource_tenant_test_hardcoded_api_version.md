# Title

Improper Use of Hardcoded API Version in URL

## File

`/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant_test.go`

## Problem

The `api-version=2021-04-01` is hardcoded in the URL used in the `httpmock.RegisterResponder` call. Relying on a hardcoded API version can lead to issues if the API changes, as the hardcoded value will not reflect the updated version.

## Impact

Hardcoding API versions can lead to maintainability issues and potential test failures when the API evolves. It reduces flexibility and creates a dependency on a specific version, which may become obsolete. Severity: **medium**

## Location

Code location where the issue exists:

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01`,
```

## Fix

Consider refactoring the code to either:

1. Use a configurable constant for the API version.
2. Fetch the API version dynamically if possible.

Example fix:

```go
const apiVersion = "2021-04-01"
httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=%s", apiVersion),
```
