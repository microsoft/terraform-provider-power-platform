# Issue: Recursive Call Without Max Retry or Backoff

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/api_managed_environment.go

## Problem

The functions `EnableManagedEnvironment` and `DisableManagedEnvironment` both implement recursive retry logic when an operation fails (when `lifecycleResponse.State.Id == "Failed"`). However, there is no limit to the number of retries, nor any backoff, which could lead to stack overflows or unintentional infinite loops.

## Impact

The lack of a maximum retry count, exponential backoff, or circuit breaker pattern could introduce serious stability risks such as stack overflow, resource exhaustion, or API spamming if the remote service experiences prolonged failures. This is a **high severity issue**.

## Location

- `EnableManagedEnvironment`  
- `DisableManagedEnvironment`  
- Lines: Recursive calls within both functions

## Code Issue

```go
if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
    if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
        return err
    }
    tflog.Info(ctx, "Managed Environment Enablement Operation failed. Retrying...")
    return client.EnableManagedEnvironment(ctx, managedEnvSettings, environmentId)
}
```
(similar for `DisableManagedEnvironment`)

## Fix

Implement a bounded retry mechanism with backoff (example for `EnableManagedEnvironment`):

```go
func (client *client) EnableManagedEnvironment(ctx context.Context, managedEnvSettings environment.GovernanceConfigurationDto, environmentId string) error {
    const maxRetries = 5
    var attempt int

    for {
        // ... (make API call as before)

        if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
            attempt++
            if attempt >= maxRetries {
                return fmt.Errorf("max retries reached for enablement operation")
            }
            if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
                return err
            }
            tflog.Info(ctx, "Managed Environment Enablement Operation failed. Retrying...")
            continue
        }
        return nil
    }
}
```
Repeat for `DisableManagedEnvironment` and tune `maxRetries` as needed.

---

This peer-reviewed markdown will be saved to:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_managed_environment_high.md`
