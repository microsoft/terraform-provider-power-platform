# Title

Use of Hardcoded API Version in `GetPowerApps`

##

`/workspaces/terraform-provider-power-platform/internal/services/powerapps/api_powerapps.go`

## Problem

The API version "2023-06-01" is hardcoded into the query parameters of the URL. This makes updates to the API version more error-prone and requires changes at multiple points in the codebase.

## Impact

- Severity: **Medium**
- Hardcoding reduces maintainability and flexibility.
- Increased likelihood of bugs during API version updates.

## Location

Line: Inside `GetPowerApps`

## Code Issue

```go
values.Add("api-version", "2023-06-01")
```

## Fix

Define the API version as a constant in the `constants` package and use that variable here instead of hardcoding.

```go
values.Add("api-version", constants.ApiVersionPowerApps)
```

Also, in the `constants` package:

```go
package constants

const ApiVersionPowerApps = "2023-06-01"
```