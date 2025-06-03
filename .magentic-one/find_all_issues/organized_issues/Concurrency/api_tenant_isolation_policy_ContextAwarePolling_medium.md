# Issue: Direct Polling with Context-Insensitive Sleep Logic (Resource Management)

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

In the `doWaitForLifecycleOperationStatus` function, the code polls for operation completion by calling `client.Api.SleepWithContext(ctx, waitTime)`. While this is generally correct, if the underlying implementation of `SleepWithContext` does not correctly honor context cancellation (e.g., due to timeouts or application shutdown), this could lead to unwanted resource usage or delayed shutdown. If not already ensured, the function should immediately exit and return if `ctx` is done, and avoid orphaning goroutines or consuming compute cycles.

## Impact

Medium. Failing to respect context cancellation can cause inefficient resource usage or delayed termination, especially under heavy load or during shutdown.

## Location

In `doWaitForLifecycleOperationStatus`:

```go
		// Wait before polling again
		err = client.Api.SleepWithContext(ctx, waitTime)
		if err != nil {
			return nil, fmt.Errorf("polling interrupted: %w", err)
		}
```

## Code Issue

```go
		// Wait before polling again
		err = client.Api.SleepWithContext(ctx, waitTime)
		if err != nil {
			return nil, fmt.Errorf("polling interrupted: %w", err)
		}
```

## Fix

Verify and document that `SleepWithContext` immediately returns if `ctx` is done/canceled. Add explicit context check before/after for belt-and-suspenders safety.

```go
		// Wait before polling again (honors context cancellation)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			err = client.Api.SleepWithContext(ctx, waitTime)
			if err != nil {
				return nil, fmt.Errorf("polling interrupted: %w", err)
			}
		}
```
And document context-sensitive wait in method comment.

---
