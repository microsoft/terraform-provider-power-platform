# Title

Hard-coded URL pattern usage in HTTP mock registration

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave_test.go`

## Problem

In multiple locations in the file, hard-coded URL patterns are used for HTTP mock registration (e.g., `https://api.admin.powerplatform.microsoft.com/api/environments/...`). Using hard-coded values in test cases can make the code less flexible and harder to maintain. If the endpoint URLs were to change, updates would be required in multiple places in the code.

## Impact

This issue impacts code maintainability and extensibility. Any changes to the endpoint structure would require manual updates across several hard-coded places. Furthermore, it introduces brittleness in tests, as the tests are tightly coupled with specific URL patterns.

**Severity:** Medium

## Location

This issue is found in the following functions:
- `registerOrganizationsMock`
- `registerEnvironmentMock`
- `TestUnitEnvironmentWaveResource_Create`
- `TestUnitEnvironmentWaveResource_Error`
- `TestUnitEnvironmentWaveResource_NotFound`
- `TestUnitEnvironmentWaveResource_FailedDuringUpgrade`
- `TestUnitEnvironmentWaveResource_UnsupportedState`

## Code Issue

An example snippet is below:

```go
httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/tenants/mytenant/organizations$`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, testFolder, "get_organizations.json")), nil
    })
```

## Fix

Introduce constants or configuration for endpoint URL patterns. This approach makes the code easier to maintain and less prone to errors during updates.

Example fix:

```go
const (
    baseAdminAPIURL = "https://api.admin.powerplatform.microsoft.com/api"
    baseEnvironmentAPIURL = "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments"
)

httpmock.RegisterResponder("GET", fmt.Sprintf(`=~^%s/tenants/mytenant/organizations$`, baseAdminAPIURL),
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, testFolder, "get_organizations.json")), nil
    })
```

In this fix, base URLs are defined as constants, and the test functions dynamically construct the full URL using these constants. If the base URL changes, it only needs adjustment at one place.