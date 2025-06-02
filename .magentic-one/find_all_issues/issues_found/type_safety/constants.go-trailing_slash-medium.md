# Title

Trailing Slash Typos in Constant Values (e.g. USGOV_ANALYTICS_SCOPE)

##

/workspaces/terraform-provider-power-platform/internal/constants/constants.go

## Problem

Some constant values (e.g. `USGOV_ANALYTICS_SCOPE`) have a double slash (`//`) in the URL path, which is likely a typographical error:

```go
USGOV_ANALYTICS_SCOPE = "https://gcc.adminanalytics.powerplatform.microsoft.us//.default"
```
should likely be:
```go
USGOV_ANALYTICS_SCOPE = "https://gcc.adminanalytics.powerplatform.microsoft.us/.default"
```

## Impact

Medium severity. Typos in URLs or resource scopes can lead to authentication problems, incorrect API requests, or subtle bugs that are difficult to trace.

## Location

```go
USGOV_ANALYTICS_SCOPE = "https://gcc.adminanalytics.powerplatform.microsoft.us//.default"
```

## Code Issue

```go
USGOV_ANALYTICS_SCOPE = "https://gcc.adminanalytics.powerplatform.microsoft.us//.default"
```

## Fix

Remove the extraneous `/` from the value:

```go
USGOV_ANALYTICS_SCOPE = "https://gcc.adminanalytics.powerplatform.microsoft.us/.default"
```
Review nearby definitions for similar unintentional redundancy in resource scope strings.

---

**Save location:**  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/constants.go-trailing_slash-medium.md
