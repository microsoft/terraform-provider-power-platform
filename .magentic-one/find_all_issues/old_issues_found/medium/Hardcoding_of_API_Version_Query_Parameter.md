# Title: Hardcoding of API Version Query Parameter

##
`/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/client.go`

## Problem

The query parameter `api-version` is hardcoded as "2.0". This is not flexible and makes the code less maintainable, especially if the version needs to be updated.

## Impact

Reduced flexibility and potential structural errors when updating the version. Severity: **medium**.

## Location

When query parameters are added to the rulesUrl.

## Code Issue

```go
queryParams.Add("api-version", "2.0")
```

## Fix

Add a constant or configuration setting for API versioning.

```go
queryParams.Add("api-version", constants.API_VERSION)
```
Where `API_VERSION` is a new constant defined in the `constants` package.