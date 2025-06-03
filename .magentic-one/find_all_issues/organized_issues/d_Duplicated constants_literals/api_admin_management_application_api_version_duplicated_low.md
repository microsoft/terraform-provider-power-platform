# Title

API Version Value Hard-Coded in Each Function

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem

The API version value `"2020-10-01"` is duplicated as a string literal in multiple places. It should be a constant to make version management easier.

## Impact

Low to Medium. Maintainability issue: updating version means touching every occurrence, and it is easy to miss one.

## Location

Occurs in each instance of `url.Values{...}` for constructing requests.

## Code Issue

```go
// Example
RawQuery: url.Values{
	constants.API_VERSION_PARAM: []string{"2020-10-01"},
}.Encode(),
```

## Fix

Define and use a constant, probably in `constants` package:

```go
// In constants package
const ADMIN_MANAGEMENT_APP_API_VERSION = "2020-10-01"

// In usage
RawQuery: url.Values{
	constants.API_VERSION_PARAM: []string{constants.ADMIN_MANAGEMENT_APP_API_VERSION},
}.Encode(),
```
