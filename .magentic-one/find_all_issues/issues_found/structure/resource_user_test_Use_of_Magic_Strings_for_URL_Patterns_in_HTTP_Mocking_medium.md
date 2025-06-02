# Title

Use of Magic Strings for URL Patterns in HTTP Mocking

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go

## Problem

Throughout the unit test functions, URLs and parts of URLs (including GUIDs) are hardcoded (e.g. `"00000000-0000-0000-0000-000000000001"`). This makes it harder to update environments, switch to new test fixtures, or spot test boilerplate errors. These repeated hardcoded strings create brittle tests that are not easily maintainable.

## Impact

This increases maintenance burden, increases the risk of copy-paste mistakes, and makes the code brittle if endpoints or test data need to change. Severity: medium.

## Location

All HTTP mock responder registration locations with URLs/IDs.

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?...`,
  ...
)
...
resource "powerplatform_user" "new_user" {
    environment_id = "00000000-0000-0000-0000-000000000001"
    aad_id = "00000000-0000-0000-0000-000000000002"
  ...
}
```

## Fix

Extract all GUIDs and base service URLs into constants (or, for extremely common reusable fixtures, utility helpers). This improves clarity, reduces the risk of typos, and makes changes easier.

```go
const (
    envID1 = "00000000-0000-0000-0000-000000000001"
    envID2 = "00000000-0000-0000-0000-000000000002"
    userID = "00000000-0000-0000-0000-000000000002"
    baseAPI = "https://api.bap.microsoft.com/providers/..."
    // etc.
)
...
httpmock.RegisterResponder("GET", fmt.Sprintf("%s/environments/%s?...", baseAPI, envID1), ...)
...
resource "powerplatform_user" ... {
    environment_id = envID1
    aad_id = userID
...
```
