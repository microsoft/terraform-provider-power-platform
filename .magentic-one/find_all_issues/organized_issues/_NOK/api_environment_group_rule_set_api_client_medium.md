# API Client: Using Hardcoded API-Version in All Method Calls

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/api_environment_group_rule_set.go

## Problem

The `"api-version"` query parameter is hardcoded as `"2021-10-01-preview"` in all requests. This can lead to maintenance issues if the API version changes, as every instance will have to be updated manually.

## Impact

This affects maintainability and upgradability. Changing API versions for different resources could result in inconsistent behavior or bugs, and tight coupling of version text. Severity: Medium.

## Location

```go
values.Add("api-version", "2021-10-01-preview")
```
(Occurs in every method constructing a URL using `values.Add`.)

## Code Issue

```go
values.Add("api-version", "2021-10-01-preview")
```

## Fix

Define the API version as a constant or configuration value, and use it throughout:

```go
const apiVersion = "2021-10-01-preview"
//...

values.Add("api-version", apiVersion)
```

Or, better, make it part of a global config or client struct.

---

This issue will be saved in:
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/api_client/api_environment_group_rule_set_api_client_medium.md
