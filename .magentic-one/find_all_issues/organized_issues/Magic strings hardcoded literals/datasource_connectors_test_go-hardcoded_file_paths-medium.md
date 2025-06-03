# Title

Hardcoded File Paths Reduce Test Portability

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors_test.go

## Problem

File paths in the HTTP mock responders are hardcoded. This can break tests when directory structures change, or if the tests are run from a location where the paths are invalid.

## Impact

Reduces test portability and flexibility, leading to fragile tests. Severity: medium.

## Location

Lines 46–62

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
    })
```

## Fix

Store files in locations relative to the test file or use test helpers to resolve absolute paths, and document the requirement.

```go
path := filepath.Join("tests", "Validate_Read", "get_virtual.json")
data := httpmock.File(path).String()
return httpmock.NewStringResponse(http.StatusOK, data), nil
```
