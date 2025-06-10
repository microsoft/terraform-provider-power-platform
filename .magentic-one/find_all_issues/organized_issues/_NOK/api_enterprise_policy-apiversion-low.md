# Hardcoded API-Version String

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

The API version "2019-10-01" is hardcoded in multiple locations.

## Impact

Severity: **low/medium**. If the API version changes, it's easy to miss an update, risking broken or inconsistent code.

## Location

```go
values.Add("api-version", "2019-10-01")
```

## Code Issue

```go
values.Add("api-version", "2019-10-01")
```

## Fix

Create a constant for the API version at the top of the file:

```go
const apiVersion = "2019-10-01"
// ...
values.Add("api-version", apiVersion)
```
