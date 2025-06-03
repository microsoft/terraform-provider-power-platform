# Issue: Magic constants used for retry/backoff without documentation or configuration

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

Several places use hardcoded retry timings via `api.DefaultRetryAfter()` and constant strings for provisioning state (e.g., `"Succeeded"`, `"Failed"`, `"LinkedDatabaseProvisioning"`, etc) and fixed HTTP status code lists (e.g., `[]int{http.StatusOK, http.StatusAccepted, http.StatusConflict}`). These are magic values and would benefit from documentation, centralization, or configuration for easier adjustments and maintainability.

## Impact

- Severity: Low
- Makes maintenance and debugging harder if a retry interval or provisioning logic changes.
- Hard to find widespread usages for tuning.

## Location

Examples throughout the file, e.g.,

```go
if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
    return err
}
```

And provisioning state checks:

```go
if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
    // ...
} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
    // ...
}
```

## Code Issue

```go
if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
    return err
}
```

Or provisioning states:

```go
if lifecycleEnv.Properties.ProvisioningState == "Succeeded" {
    // ...
} else if lifecycleEnv.Properties.ProvisioningState != "LinkedDatabaseProvisioning" && lifecycleEnv.Properties.ProvisioningState != "Succeeded" {
    // ...
}
```

## Fix

Document the rationale for backoff/retry values and magic strings; consider a centralized definition:

```go
const (
    ProvisioningStateSucceeded = "Succeeded"
    ProvisioningStateFailed = "Failed"
    ProvisioningStateLinkedDatabaseProvisioning = "LinkedDatabaseProvisioning"
    // ...more as needed
)
```
And for retry intervals:

```go
const DefaultRetryAfter = 15 * time.Second  // document why!
```
Or, link documentation directly above with context for the value.

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_environment_magic_constants_low.md`
