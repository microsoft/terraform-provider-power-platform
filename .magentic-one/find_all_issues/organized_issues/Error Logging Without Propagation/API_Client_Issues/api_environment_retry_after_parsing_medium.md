# Issue: Retry-After header parsing logic is inconsistent and misleading

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

The parsing for the `Retry-After` header in `AddDataverseToEnvironment` currently attempts to parse it directly as a duration via `time.ParseDuration(retryHeader)`, then multiplies it by `time.Second` if successful. However, the `Retry-After` header from HTTP responses is usually either a number (seconds as an integer) or a date string, not a Go duration like `"10s"`.

The logic of `else { retryAfter = retryAfter * time.Second }` is misleading and suggests potential unintended long sleep intervals or errors in waiting logic.

## Impact

- Severity: Medium
- May cause incorrect wait intervals (longer/shorter than intended)
- May cause confusion during further maintenance
- Introduces risk of rapid retry or unnecessary delay in Dataverse provisioning loops

## Location

In `AddDataverseToEnvironment`:

```go
retryHeader := apiResponse.GetHeader(constants.HEADER_RETRY_AFTER)
tflog.Debug(ctx, "Retry Header: "+retryHeader)
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    retryAfter = api.DefaultRetryAfter()
} else {
    retryAfter = retryAfter * time.Second
}
```

## Code Issue

```go
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    retryAfter = api.DefaultRetryAfter()
} else {
    retryAfter = retryAfter * time.Second
}
```

## Fix

For an HTTP numeric `Retry-After` header (which is the common case), parse as an integer (seconds), then convert to a `time.Duration`. Only fall back to the default if it cannot be parsed as an int.

```go
retryHeader := apiResponse.GetHeader(constants.HEADER_RETRY_AFTER)
tflog.Debug(ctx, "Retry Header: "+retryHeader)
retryAfter := api.DefaultRetryAfter()

if retryHeader != "" {
    if seconds, err := strconv.Atoi(retryHeader); err == nil {
        retryAfter = time.Duration(seconds) * time.Second
    }
    // Optionally handle date string case if ever expected
}
```

Make sure to import `strconv` at the top:

```go
import (
    // other imports...
    "strconv"
)
```

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/api_client/api_environment_retry_after_parsing_medium.md`
