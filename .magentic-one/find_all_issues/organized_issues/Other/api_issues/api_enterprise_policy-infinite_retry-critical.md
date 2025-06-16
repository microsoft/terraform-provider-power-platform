# Potential Infinite Recursion on Policy (Un)Linking Failure

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

If the policy operation fails (i.e., `lifecycleResponse.State.Id == "Failed"`), the function sleeps and then recursively calls itself without any retry limit or backoff. This may cause stack overflow (infinite recursion) and resource exhaustion if the operation perpetually fails.

## Impact

**Severity: high/critical**. This can result in panics, unbounded goroutine and stack growth, OOM errors, or the provider locking up.

## Location

Both `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy`:

```go
if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
	if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
		return err
	}
	tflog.Info(ctx, "Policy Linking Operation failed. Retrying...")
	return client.LinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
}
```

## Code Issue

```go
return client.LinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
// ...
return client.UnLinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
```

## Fix

Introduce a retry counter/policy, either as a parameter or constant, and stop retrying after a certain threshold. Example for `LinkEnterprisePolicy`:

```go
func (client *Client) LinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string, retriesLeft int) error {
	// ... existing code ...
	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if retriesLeft <= 0 {
			return fmt.Errorf("Policy Linking Operation failed after maximum retries")
		}
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Policy Linking Operation failed. Retrying...")
		return client.LinkEnterprisePolicy(ctx, environmentId, environmentType, systemId, retriesLeft-1)
	}
	return nil
}
```

Or, encapsulate retry logic in a loop instead of recursion, which is a more idiomatic Go approach.
