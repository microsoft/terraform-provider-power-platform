# Hardcoded API Version in Multiple Places

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go

## Problem

The API version `2022-03-01-preview` is hardcoded in multiple places, which can make upgrades difficult and error-prone.

## Impact

Low. This does not cause immediate bugs but affects maintainability and upgradability.

## Location

All occurrences of:

```go
values.Add("api-version", "2022-03-01-preview")
```

## Fix

Define the API version as a constant at the top of the file (if not already in use):

```go
const apiVersion = "2022-03-01-preview"
```
Then use:

```go
values.Add("api-version", apiVersion)
```
