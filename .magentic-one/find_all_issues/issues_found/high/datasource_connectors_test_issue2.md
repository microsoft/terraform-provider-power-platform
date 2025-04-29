# Title

Usage of hardcoded URLs for API endpoints results in fragile test designs.

# Path

`/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors_test.go`

## Problem

Multiple hardcoded URLs are used in tests (`httpmock.RegisterResponder`). If the URL endpoints change or become unavailable, these tests would fail unexpectedly. An approach using environment variables or constants to handle dynamic URLs is recommended.

## Impact

Hardcoded URLs make tests fragile and difficult to maintain. If the base URL changes, multiple tests need updates. Severity: high.

## Location

Function: `TestUnitConnectorsDataSource_Validate_Read`, Line: 41-77

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
    })
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_unblockable.json").String()), nil
    })
httpmock.RegisterResponder("GET", `https://api.powerapps.com/providers/Microsoft.PowerApps/apis?...`),
```

## Fix

Define constants or use configuration variables for these URLs. For example:

```go
const (
    baseAPIURL = "https://api.bap.microsoft.com/providers"
    powerAppsURL = "https://api.powerapps.com/providers/Microsoft.PowerApps/apis"
)

httpmock.RegisterResponder("GET", fmt.Sprintf("%s/PowerPlatform.Governance/v1/connectors/metadata/virtual", baseAPIURL),
    ...)
httpmock.RegisterResponder("GET", fmt.Sprintf("%s/PowerPlatform.Governance/v1/connectors/metadata/unblockable", baseAPIURL),
    ...)
httpmock.RegisterResponder("GET", fmt.Sprintf("%s?api-version=2019-05-01", powerAppsURL),
    ...)
```