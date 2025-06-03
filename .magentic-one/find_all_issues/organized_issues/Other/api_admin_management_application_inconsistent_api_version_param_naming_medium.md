# Title

Inconsistent API Version Parameter Naming

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem

The API version query parameter key is inconsistently named across functions: in `GetAdminApplication`, the key is `constants.API_VERSION_PARAM`, while in `RegisterAdminApplication` and `UnregisterAdminApplication`, a string literal `"api-version"` is used. This inconsistency may lead to maintainability issues and can cause bugs if the API version parameter needs to be updated.

## Impact

Medium. This makes the code harder to maintain, increases the risk of subtle bugs if the version needs to be changed, and potentially could result in runtime errors if the constants diverge or are incorrectly modified in one place but not others.

## Location

`GetAdminApplication`, `RegisterAdminApplication`, and `UnregisterAdminApplication` functions, at the construction of the `apiUrl` variable.

## Code Issue

```go
// GetAdminApplication
RawQuery: url.Values{
	constants.API_VERSION_PARAM: []string{"2020-10-01"},
}.Encode(),

// RegisterAdminApplication / UnregisterAdminApplication
RawQuery: url.Values{
	"api-version": []string{"2020-10-01"},
}.Encode(),
```

## Fix

Use `constants.API_VERSION_PARAM` in all cases for setting the API version query parameter, enhancing maintainability and reducing the chance of typo errors.

```go
// RegisterAdminApplication and UnregisterAdminApplication
RawQuery: url.Values{
	constants.API_VERSION_PARAM: []string{"2020-10-01"},
}.Encode(),
```
